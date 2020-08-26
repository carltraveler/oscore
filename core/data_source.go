package core

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/core/http"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	http2 "net/http"
	"net/url"
	"sort"
	"strings"
	"sync/atomic"
)

// the crypto info maybe fill both in body or headers.
func GetWithCrptoBodyHeaders(bodymapsVars map[string]interface{}, apiCryptoInfos []*tables.ApiCryptoInfo, baseUrl string, referParams []*tables.RequestParam, params []*tables.RequestParam) (map[string]interface{}, []*tables.ApiHeadValues, string, error) {
	cryptoHeaderKV := make([]*tables.ApiHeadValues, 0)
	firstQueryArg := !strings.Contains(baseUrl, "?")

	for _, apiCryptoInfo := range apiCryptoInfos {
		var cryptoResult string
		switch apiCryptoInfo.CryptoType {
		case tables.API_CRYTPO_TYPE_BANDCARD_VERIFY:
			bodySortsKey := make([]string, 0)
			for k := range bodymapsVars {
				bodySortsKey = append(bodySortsKey, k)
			}

			sort.Strings(bodySortsKey)

			var bodySign string
			for _, k := range bodySortsKey {
				tmps := bodymapsVars[k].(string)
				bodySign = bodySign + k + "=" + tmps
			}
			log.Debugf("HandleDataSourceReqCore Y.y bodySign %s", bodySign)

			bodySign = bodySign + apiCryptoInfo.CryptoValue
			cryptoResult = fmt.Sprintf("%x", md5.Sum([]byte(bodySign)))
		case tables.API_CRYTPO_TYPE_YI_MEI:
			if len(apiCryptoInfos) != 1 {
				return nil, nil, "", fmt.Errorf("yi mei cryptoinfo must only one key called sign")
			}

			if apiCryptoInfo.CryptoKey != "sign" {
				return nil, nil, "", fmt.Errorf("CryptoKey must be sign")
			}

			var appId string
			var timestamp string

			for _, p := range referParams {
				if p.ParamName == "appId" {
					appId = p.Note
				}

				if p.ParamName == "timestamp" {
					timestamp = p.Note
				}
			}

			if appId == "" || timestamp == "" {
				return nil, nil, "", fmt.Errorf("yi mei crypto info need appId and timestamp info")
			}

			bodySign := appId + apiCryptoInfo.CryptoValue + timestamp
			cryptoResult = fmt.Sprintf("%x", md5.Sum([]byte(bodySign)))
		default:
			return nil, nil, "", fmt.Errorf("GetWithCrptoBodyHeaders can not support %d kind of crypto method.", apiCryptoInfo.CryptoType)
		}

		switch apiCryptoInfo.CryptoWhere {
		case tables.API_CRYTPO_FILL_IN_BODY:
			bodymapsVars[apiCryptoInfo.CryptoKey] = cryptoResult
		case tables.API_CRYTPO_FILL_IN_HEADER:
			cryptoHeaderKV = append(cryptoHeaderKV, &tables.ApiHeadValues{
				HeaderKey:   apiCryptoInfo.CryptoKey,
				HeaderValue: cryptoResult,
			})
		case tables.API_CRYTPO_FILL_IN_QUERY:
			if firstQueryArg {
				baseUrl = baseUrl + "?" + apiCryptoInfo.CryptoKey + "=" + cryptoResult
				firstQueryArg = true
			} else {
				baseUrl = baseUrl + "&" + apiCryptoInfo.CryptoKey + "=" + cryptoResult
			}
		default:
			return nil, nil, "", fmt.Errorf("can not fill crypto info only in body or headers.")
		}
	}

	return bodymapsVars, cryptoHeaderKV, baseUrl, nil
}

func HandleDataSourceReqCore(tx *sqlx.Tx, oscoreUrlKey string, params []*tables.RequestParam, apiKey string, publishTestOnly bool) ([]byte, error) {
	log.Debugf("HandleDataSourceReqCore : %v", params)
	var apiState []int32

	if publishTestOnly {
		apiState = []int32{tables.API_STATE_PUBLISH}
	} else {
		apiState = []int32{tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER}
	}

	info, err := dao.DefOscoreApiDB.QueryApiBasicInfoByOscoreUrlKeyStates(tx, oscoreUrlKey, apiState)
	if err != nil {
		return nil, err
	}

	cryptoInfos, err := dao.DefOscoreApiDB.QueryApiCryptoInfoByApiId(tx, info.ApiId)
	if err != nil {
		log.Debugf("HandleDataSourceReqCore N.0.0.0 %s", err)
		return nil, err
	}

	referParams, err := dao.DefOscoreApiDB.QueryReferRequestParamByApiId(tx, info.ApiId)
	if err != nil {
		return nil, err
	}
	var varParmCount int
	for _, p := range referParams {
		if p.ParamTag != tables.REQUEST_PARAM_TAG_FIX {
			varParmCount += 1
		}
	}

	if varParmCount != len(params) {
		return nil, fmt.Errorf("params len error. should be %d", varParmCount)
	}

	var firstQueryArg bool
	var bodyParam []byte
	// use param
	var j uint32
	baseUrl := info.ApiProvider
	firstQueryArg = !strings.Contains(baseUrl, "?")
	bodymapsVars := make(map[string]interface{})
	bodyParamNum := uint32(0)
	headers := make([]*tables.ApiHeadValues, 0)
	for i, p := range referParams {
		log.Debugf("publish params[%d]: %v", i, p)
		if p.ParamTag != tables.REQUEST_PARAM_TAG_FIX && p.ParamTag != tables.REQUEST_PARAM_TAG_PARAM {
			return nil, fmt.Errorf("param tag input error.")
		}

		if p.ParamTag != tables.REQUEST_PARAM_TAG_FIX {
			if p.ParamName != params[j].ParamName || p.ParamWhere != params[j].ParamWhere {
				return nil, fmt.Errorf("params error. should be %v", referParams)
			}

			if publishTestOnly {
				if params[j].Required && params[j].Note == "" {
					return nil, fmt.Errorf("params[%d] %s is Required. but empty", j, params[j].ParamName)
				}
			} else {
				if params[j].Required && params[j].ValueDesc == "" {
					return nil, fmt.Errorf("params[%d] %s is Required. but empty", j, params[j].ParamName)
				}
			}
		}

		switch p.ParamWhere {
		case tables.URL_PARAM_HEADER:
			if p.ParamTag == tables.REQUEST_PARAM_TAG_FIX {
				h := &tables.ApiHeadValues{
					HeaderKey:   p.ParamName,
					HeaderValue: p.Note,
				}

				headers = append(headers, h)
			} else {
				var valueStr string
				if publishTestOnly {
					valueStr = params[j].Note
				} else {
					valueStr = params[j].ValueDesc
				}
				if valueStr != "" {
					h := &tables.ApiHeadValues{
						HeaderKey:   params[j].ParamName,
						HeaderValue: valueStr,
					}

					headers = append(headers, h)
				}
			}
		case tables.URL_PARAM_RESTFUL:
			if !firstQueryArg {
				return nil, fmt.Errorf("params error. restful url after query.")
			}

			if p.ParamTag == tables.REQUEST_PARAM_TAG_FIX {
				baseUrl = baseUrl + "/" + p.Note
			} else {
				if publishTestOnly {
					if params[j].Note != "" {
						baseUrl = baseUrl + "/" + params[j].Note
					}
				} else {
					if params[j].ValueDesc != "" {
						baseUrl = baseUrl + "/" + params[j].ValueDesc
					}
				}
			}
		case tables.URL_PARAM_QUERY:
			if firstQueryArg {
				if p.ParamTag == tables.REQUEST_PARAM_TAG_FIX {
					p.Note = url.QueryEscape(p.Note)
					baseUrl = baseUrl + "?" + p.ParamName + "=" + p.Note
				} else {
					if publishTestOnly {
						if params[j].Note != "" {
							params[j].Note = url.QueryEscape(params[j].Note)
							baseUrl = baseUrl + "?" + params[j].ParamName + "=" + params[j].Note
						}
					} else {
						if params[j].ValueDesc != "" {
							params[j].ValueDesc = url.QueryEscape(params[j].ValueDesc)
							baseUrl = baseUrl + "?" + params[j].ParamName + "=" + params[j].ValueDesc
						}
					}
				}
				firstQueryArg = false
			} else {
				if p.ParamTag == tables.REQUEST_PARAM_TAG_FIX {
					p.Note = url.QueryEscape(p.Note)
					baseUrl = baseUrl + "&" + p.ParamName + "=" + p.Note
				} else {
					if publishTestOnly {
						if params[j].Note != "" {
							params[j].Note = url.QueryEscape(params[j].Note)
							baseUrl = baseUrl + "&" + params[j].ParamName + "=" + params[j].Note
						}
					} else {
						if params[j].ValueDesc != "" {
							params[j].ValueDesc = url.QueryEscape(params[j].ValueDesc)
							baseUrl = baseUrl + "&" + params[j].ParamName + "=" + params[j].ValueDesc
						}
					}
				}
			}
		case tables.URL_PARAM_BODY:
			bodymapsTmps := make(map[string]interface{})
			if info.RequestType == tables.API_REQUEST_GET {
				return nil, fmt.Errorf("params error. can not set body param in get request.")
			}

			if bodyParamNum != 0 {
				return nil, fmt.Errorf("params error. can not pass multi body param.")
			}

			var valueStr string
			if p.ParamTag != tables.REQUEST_PARAM_TAG_FIX {
				bodyParamNum += 1
				if publishTestOnly {
					valueStr = params[j].Note
				} else {
					valueStr = params[j].ValueDesc
				}

				if valueStr != "" {
					err = json.Unmarshal([]byte(valueStr), &bodymapsTmps)
					if err != nil {
						log.Errorf("HandleDataSourceReqCore N.restart.0 %s", err)
						return nil, fmt.Errorf("HandleDataSourceReqCore. body json format error. %s", err)
					}

					for k, v := range bodymapsTmps {
						if _, ok := bodymapsVars[k]; ok {
							log.Errorf("HandleDataSourceReqCore N.restart.1 %s", err)
							return nil, fmt.Errorf("HandleDataSourceReqCore key %s already exists", p.ParamName)
						}
						bodymapsVars[k] = v
					}
				}
			} else {
				if _, ok := bodymapsVars[p.ParamName]; ok {
					log.Errorf("HandleDataSourceReqCore N.restart.2 %s", err)
					return nil, fmt.Errorf("HandleDataSourceReqCore key %s already exists", p.ParamName)
				}
				bodymapsVars[p.ParamName] = p.Note
			}
		}

		if p.ParamTag != tables.REQUEST_PARAM_TAG_FIX {
			j += 1
		}
	}

	var cryptoHeaderKV []*tables.ApiHeadValues
	if len(cryptoInfos) != 0 {
		bodymapsVars, cryptoHeaderKV, baseUrl, err = GetWithCrptoBodyHeaders(bodymapsVars, cryptoInfos, baseUrl, referParams, params)
		if err != nil {
			log.Errorf("HandleDataSourceReqCore N.0 GetBodyWithCrpto err %s", err)
			return nil, err
		}
	}

	bodyParam, err = json.Marshal(bodymapsVars)
	if err != nil {
		log.Errorf("HandleDataSourceReqCore N.1 %s", err)
		return nil, err
	}

	log.Debugf("HandleDataSourceReqCore Y.0 bodyParam %s", bodyParam)

	log.Debugf("baseUrl: %s", baseUrl)

	var key *tables.APIKey
	var apiCounterP *uint64
	if !publishTestOnly {
		key, apiCounterP, err = DefOscoreApi.Cache.BeforeCheckApiKey(apiKey, info.ApiId)
		if err != nil {
			return nil, err
		}
	}

	log.Debugf("HandleDataSourceReqCore Y.1 headers %v", headers)

	// fill cryptoHeaderKV to header
	if cryptoHeaderKV != nil {
		headers = append(headers, cryptoHeaderKV...)
	}

	switch info.RequestType {
	case tables.API_REQUEST_GET:
		log.Debugf("headers: len %d,%v", len(headers), headers)
		res, statuscode, err := http.DefClient.GetWithHeader(baseUrl, headers)
		if err != nil {
			if !publishTestOnly {
				atomic.AddUint64(&key.UsedNum, ^uint64(0))
				atomic.AddUint64(apiCounterP, ^uint64(0))
			}
			return nil, err
		}

		if statuscode != http2.StatusOK {
			return nil, fmt.Errorf("request get error msg: %d, %s", statuscode, string(res))
		}

		if !publishTestOnly && key.ApiKeyType == tables.API_KEY_TYPE_COUNT {
			DefOscoreApi.Cache.UpdateFreq <- apiKey
		}
		log.Debugf("%s", string(res))
		return res, nil
	case tables.API_REQUEST_POST:
		res, statuscode, err := http.DefClient.PostWithHeader(baseUrl, headers, bodyParam)
		if err != nil {
			if !publishTestOnly {
				atomic.AddUint64(&key.UsedNum, ^uint64(0))
				atomic.AddUint64(apiCounterP, ^uint64(0))
			}
			return nil, err
		}
		if !publishTestOnly && key.ApiKeyType == tables.API_KEY_TYPE_COUNT {
			DefOscoreApi.Cache.UpdateFreq <- apiKey
		}

		if statuscode != http2.StatusOK {
			return nil, fmt.Errorf("request get error msg: %s", string(res))
		}

		log.Debugf("%s", string(res))
		return res, nil
	}

	return nil, nil
}
