package tables

import (
	"encoding/json"
	"fmt"
	"time"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	x := time.Time(t).Unix()
	return json.Marshal(x)
}

// for ApiBasicInfo.ApiState
const (
	API_STATE_IGNOR         int32 = -1
	API_STATE_INVALID       int32 = 0
	API_STATE_BUILTIN       int32 = 1
	API_STATE_PUBLISH       int32 = 2
	API_STATE_DISABLE_ORDER int32 = 3
	API_STATE_DELETE        int32 = 4
	API_STATE_LAST          int32 = 5
)

const (
	API_KIND_INVALID      int32 = 0
	API_KIND_DATA_NORMAL  int32 = 1
	API_KIND_DATA_PROCESS int32 = 2
)

// for RequestParam.ParamWhere
const (
	URL_PARAM_RESTFUL int32 = 1
	URL_PARAM_QUERY   int32 = 2
	URL_PARAM_BODY    int32 = 3
	URL_PARAM_HEADER  int32 = 4
)

const (
	API_REQUEST_POST string = "POST"
	API_REQUEST_GET  string = "GET"
)

type ApiBasicInfo struct {
	ApiId               uint32 `json:"apiId" db:"ApiId"`
	ApiType             string `json:"type" db:"ApiType"`
	Icon                string `json:"icon" db:"Icon"`
	Title               string `json:"title" db:"Title"`
	ApiProvider         string `json:"provider" db:"ApiProvider"`            //source url. join args can access.
	ApiOscoreUrlKey     string `json:"apiOscoreUrlKey" db:"ApiOscoreUrlKey"` //oscoreurlkey
	ApiUrl              string `json:"apiUrl" db:"ApiUrl"`                   // oscoreurl
	Price               string `json:"price" db:"Price"`
	ApiDesc             string `json:"description" db:"ApiDesc"`
	ErrorDesc           string `json:"errorDescription" db:"ErrorDesc"`
	Specifications      uint32 `json:"specifications" db:"Specifications"`
	Popularity          uint32 `json:"popularity" db:"Popularity"`
	Delay               uint32 `json:"delay" db:"Delay"`
	SuccessRate         uint32 `json:"successRate" db:"SuccessRate"`
	InvokeFrequency     uint64 `json:"invokeFrequency" db:"InvokeFrequency"`
	ApiState            int32  `json:"apiState" db:"ApiState"`
	RequestType         string `json:"requestType" db:"RequestType"`
	Mark                string `json:"mark" db:"Mark"`
	ResponseParam       string `json:"responseParam" db:"ResponseParam"`
	ResponseExample     string `json:"responseExample" db:"ResponseExample"`
	DataDesc            string `json:"dataDesc" db:"DataDesc"`
	DataSource          string `json:"dataSource" db:"DataSource"`
	ApplicationScenario string `json:"applicationScenario" db:"ApplicationScenario"`
	ApiKind             int32  `json:"apiKind" db:"ApiKind"`
	OntId               string `json:"ontId" db:"OntId"`
	UserId              string `json:"userId" db:"UserId"`
	Author              string `json:"author" db:"Author"`
	ResourceId          string `json:"resourceId" db:"ResourceId"`
	TokenHash           string `json:"tokenHash" db:"TokenHash"`
	OwnerAddress        string `json:"ownerAddress" db:"OwnerAddress"`
	Abstract            string `json:"abstract" db:"Abstract"`
	NotifyDelete        int32  `json:"notifyDelete" db:"NotifyDelete"`
	CreateTime          int64  `json:"createTime" db:"CreateTime"`
	UpdateTime          int64  `json:"updateTime" db:"UpdateTime"`
}

type ApiTag struct {
	Id         uint32 `json:"id" db:"Id"`
	ApiId      uint32 `json:"apiId" db:"ApiId"`
	TagId      uint32 `json:"tagId" db:"TagId"`
	State      byte   `json:"state" db:"State"`
	CreateTime Time   `json:"createTime" db:"CreateTime"`
}

type Tag struct {
	Id         uint32 `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	CategoryId uint32 `json:"categoryId" db:"category_id"`
	State      byte   `json:"state" db:"state"`
	CreateTime Time   `json:"createTime" db:"create_time"`
}

type Category struct {
	Id     uint32 `json:"id" db:"id"`
	NameZh string `json:"nameZh" db:"name_zh"`
	NameEn string `json:"nameEn" db:"name_en"`
	Icon   string `json:"icon" db:"icon"`
	State  byte   `json:"state" db:"state"`
	Sort   int32  `json:"sort" db:"sort"`
}

const (
	SPEC_TYPE_COUNT    int32 = 1
	SPEC_TYPE_DURATION int32 = 2
)

type Specifications struct {
	Id     uint32 `json:"id" db:"Id"`
	ApiId  uint32 `json:"apiId" db:"ApiId"`
	Price  string `json:"price" db:"Price"`
	Amount uint64 `json:"amount" db:"Amount"`
	// month num. to ensure the db do not overflow. do not use int to insert.
	EffectiveDuration int32 `json:"effectiveDuration" db:"EffectiveDuration"`
	SpecType          int32 `json:"specType" db:"SpecType"`
}

const (
	REQUEST_PARAM_TAG_PARAM = 0
	REQUEST_PARAM_TAG_FIX   = 1
)

type RequestParam struct {
	Id         uint32 `json:"id" db:"Id"`
	ApiId      uint32 `json:"apiId" db:"ApiId"`
	ParamName  string `json:"paramName" db:"ParamName"`
	Required   bool   `json:"required" db:"Required"`
	ParamWhere int32  `json:"paramWhere" db:"ParamWhere"`
	ParamType  string `json:"paramType" db:"ParamType"`
	Note       string `json:"note" db:"Note"`
	ValueDesc  string `json:"valueDesc" db:"ValueDesc"`
	ParamTag   int32  `json:"paramTag" db:"ParamTag"`
}

func (self *RequestParam) UnmarshalJSON(buf []byte) error {
	//fmt.Println("RequestParam UnmarshalJSON begin")
	param := struct {
		Id         uint32      `json:"id"`
		ApiId      uint32      `json:"apiId"`
		ParamName  string      `json:"paramName"`
		Required   bool        `json:"required"`
		ParamWhere int32       `json:"paramWhere"`
		ParamType  string      `json:"paramType"`
		Note       interface{} `json:"note"`
		ValueDesc  interface{} `json:"valueDesc"`
		ParamTag   int32       `json:"paramTag"`
	}{}

	err := json.Unmarshal(buf, &param)
	if err != nil {
		return err
	}

	//fmt.Printf("RequestParam UnmarshalJSON %v\n", param)
	if param.ParamWhere == URL_PARAM_BODY && param.ParamTag == REQUEST_PARAM_TAG_PARAM {
		if param.ValueDesc != nil {
			vstr, err := json.Marshal(param.ValueDesc)
			if err != nil {
				return err
			}
			self.ValueDesc = string(vstr)
		}
		if param.Note != nil {
			vstr, err := json.Marshal(param.Note)
			if err != nil {
				return err
			}
			self.Note = string(vstr)
		}
	} else {
		if param.ValueDesc != nil {
			var ok bool
			self.ValueDesc, ok = param.ValueDesc.(string)
			if !ok {
				return fmt.Errorf("can not unmarshal object to ValueDesc as string")
			}
		}

		if param.Note != nil {
			var ok bool
			self.Note, ok = param.Note.(string)
			if !ok {
				return fmt.Errorf("can not unmarshal object to Note as string")
			}
		}
	}

	self.Id = param.Id
	self.ApiId = param.ApiId
	self.ParamName = param.ParamName
	self.Required = param.Required
	self.ParamWhere = param.ParamWhere
	self.ParamType = param.ParamType
	self.ParamTag = param.ParamTag
	return nil
}

type ErrorCode struct {
	Id        uint32 `json:"id" db:"Id"`
	ErrorCode int32  `json:"code" db:"ErrorCode"`
	ErrorDesc string `json:"description" db:"ErrorDesc"`
}

const (
	API_CRYTPO_TYPE_NONE            = 0
	API_CRYTPO_TYPE_BANDCARD_VERIFY = 1
	API_CRYTPO_TYPE_YI_MEI          = 2
)

const (
	API_CRYTPO_FILL_IN_BODY   = 1
	API_CRYTPO_FILL_IN_HEADER = 2
	API_CRYTPO_FILL_IN_QUERY  = 3
)

// this info will not be explorue to users.
type ApiCryptoInfo struct {
	ApiId       uint32 `json:"apiId" db:"ApiId"`
	CryptoType  int32  `json:"cryptoType" db:"CryptoType"`
	CryptoWhere int32  `json:"cryptoWhere" db:"CryptoWhere"`
	CryptoKey   string `json:"cryptoKey" db:"CryptoKey"`
	CryptoValue string `json:"cryptoValue" db:"CryptoValue"`
}
