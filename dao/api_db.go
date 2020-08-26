package dao

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/common"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/ontio/oscore/utils"
	"math/big"
	"strings"
)

func IsErrNoRows(err error) bool {
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func (this *OscoreApiDB) ClearApiBasicDB() error {
	strSql := "delete from tbl_api_basic_info"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *OscoreApiDB) ClearRequestParamDB() error {
	strSql := "delete from tbl_request_param"
	_, err := this.DB.Exec(strSql)
	return err
}

func (this *OscoreApiDB) ClearSpecificationsDB() error {
	strSql := "delete from tbl_specifications"
	_, err := this.DB.Exec(strSql)
	return err
}
func (this *OscoreApiDB) ClearApiKeyDB() error {
	strSql := "delete from tbl_api_key"
	_, err := this.DB.Exec(strSql)
	if err != nil {
		return err
	}
	strSql2 := "delete from tbl_api_test_key"
	_, err = this.DB.Exec(strSql2)
	return err
}

func (this *OscoreApiDB) ClearAll() error {
	strSql := "delete from tbl_api_test_key"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_qr_code"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_key"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_order"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_error_code"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_request_param"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_specifications"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_tag"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_tag"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_category"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_header_values"
	this.DB.Exec(strSql)
	strSql = "delete from tbl_api_basic_info"
	this.DB.Exec(strSql)

	strSql = "delete from tbl_country_city"
	this.DB.Exec(strSql)

	strSql = "delete from tbl_algorithm_env"
	this.DB.Exec(strSql)

	strSql = "delete from tbl_api_algorithm"
	this.DB.Exec(strSql)

	strSql = "delete from tbl_env"
	this.DB.Exec(strSql)

	strSql = "delete from tbl_algorithm"
	this.DB.Exec(strSql)
	return nil
}

func (this *OscoreApiDB) InsertApiBasicInfo(tx *sqlx.Tx, infos []*tables.ApiBasicInfo) error {
	var err error
	if len(infos) == 0 {
		return nil
	}

	sqlStrArr := make([]string, len(infos))
	for i, info := range infos {
		sqlStrArr[i] = fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%d','%d','%d','%d','%d','%d','%s','%s','%s','%s','%s','%s','%s','%d','%s','%s','%s','%s','%d','%d')",
			info.ApiType, info.Icon, info.Title, info.ApiProvider, info.ApiOscoreUrlKey, info.ApiUrl, info.Price,
			info.ApiDesc, info.ErrorDesc, info.Specifications, info.Popularity, info.Delay, info.SuccessRate, info.InvokeFrequency, info.ApiState, info.RequestType, info.Mark, info.ResponseParam, info.ResponseExample, info.DataDesc, info.DataSource, info.ApplicationScenario, info.ApiKind, info.OntId, info.UserId, info.Author, info.Abstract, info.CreateTime, info.CreateTime)
	}
	strSql := `insert into tbl_api_basic_info (ApiType,Icon,Title,ApiProvider,ApiOscoreUrlKey,ApiUrl,Price,ApiDesc,ErrorDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType,Mark,ResponseParam,ResponseExample,DataDesc,DataSource,ApplicationScenario,ApiKind,OntId,UserId,Author,Abstract,CreateTime,UpdateTime) values`
	strSql += strings.Join(sqlStrArr, ",")
	err = this.Exec(tx, strSql)
	return err
}

func (this *OscoreApiDB) UpdateApiBasicPrice(tx *sqlx.Tx, apiId uint32, price string) error {
	strSql := `update tbl_api_basic_info set Price=? where ApiId=?`
	return this.Exec(tx, strSql, price, apiId)
}

func (this *OscoreApiDB) UpdateApiBasicResponseExample(tx *sqlx.Tx, apiId uint32, responseExample string) error {
	strSql := `update tbl_api_basic_info set ResponseExample=? where ApiId=?`
	return this.Exec(tx, strSql, responseExample, apiId)
}

func (this *OscoreApiDB) QueryApiBasicInfoByApiId(tx *sqlx.Tx, apiId uint32, apiState int32) (*tables.ApiBasicInfo, error) {
	var err error
	info := &tables.ApiBasicInfo{}
	if apiState != tables.API_STATE_IGNOR {
		strSql := `select * from tbl_api_basic_info where ApiId=? and ApiState=?`

		err = this.Get(tx, info, strSql, apiId, apiState)
		if err != nil {
			return nil, err
		}
	} else {
		strSql := `select * from tbl_api_basic_info where ApiId=?`

		err = this.Get(tx, info, strSql, apiId)
		if err != nil {
			return nil, err
		}
	}

	return info, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByApiIds(tx *sqlx.Tx, apiIds []uint32, apiState int32) ([]*tables.ApiBasicInfo, error) {
	var err error
	infos := make([]*tables.ApiBasicInfo, 0)
	strInArgs := getUint32InSql(apiIds)
	if apiState != tables.API_STATE_IGNOR {
		strSql := fmt.Sprintf("select * from tbl_api_basic_info where ApiId in (%s) and ApiState=? order by CreateTime desc", strInArgs)
		err = this.Select(tx, &infos, strSql, apiState)
		if err != nil {
			return nil, err
		}
	} else {
		strSql := fmt.Sprintf("select * from tbl_api_basic_info where ApiId in (%s) order by CreateTime desc", strInArgs)
		err = this.Select(tx, &infos, strSql)
		if err != nil {
			return nil, err
		}
	}

	return infos, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByApiIdInState(tx *sqlx.Tx, apiId uint32, apiState []int32) (*tables.ApiBasicInfo, error) {
	var err error
	strState := getStrSate(apiState)
	strSql := fmt.Sprintf("select * from tbl_api_basic_info where ApiId=? and ApiState in (%s)", strState)

	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, apiId)
	if err != nil {
		log.Errorf("strSql: %s. apiId %d", strSql, apiId)
		log.Errorf("QueryApiBasicInfoByApiIdInState: N.0 %s", err)
		return nil, err
	}

	var enterPrise *tables.EnterpriseInfo
	if info.OntId != "" {
		enterPrise, err = this.QueryEnterpriseInfoByUserId(tx, info.UserId)
		if err != nil {
			log.Debugf("QueryApiBasicInfoByApiIdInState: %s. OntId: %s", err, info.OntId)
			info.Author = ""
			return info, nil
		}

		if enterPrise.State == 1 {
			info.Author = enterPrise.EnterpriseLegalName
		} else {
			info.Author = ""
		}

		log.Debugf("QueryApiBasicInfoByApiIdInState %d. author: %s", enterPrise.State, info.Author)
	}

	return info, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByProvider(tx *sqlx.Tx, apiProvider string) (*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where ApiProvider=?`
	info := &tables.ApiBasicInfo{}
	err := this.Get(tx, info, strSql, apiProvider)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByTitle(tx *sqlx.Tx, title string) (*tables.ApiBasicInfo, error) {
	strSql := `select * from tbl_api_basic_info where Title=?`
	info := &tables.ApiBasicInfo{}
	err := this.Get(tx, info, strSql, title)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *OscoreApiDB) DeleteApiBasicInfoByApiId(tx *sqlx.Tx, apiId uint32) error {
	strSqlt := `delete from tbl_api_tag where ApiId=?`
	err := this.Exec(tx, strSqlt, apiId)
	if err != nil {
		return err
	}
	strSqlp := `delete from tbl_request_param where ApiId=?`
	err = this.Exec(tx, strSqlp, apiId)
	if err != nil {
		return err
	}
	strSqls := `delete from tbl_specifications where ApiId=?`
	err = this.Exec(tx, strSqls, apiId)
	if err != nil {
		return err
	}
	strSqlh := `delete from tbl_api_header_values where ApiId=?`
	err = this.Exec(tx, strSqlh, apiId)
	if err != nil {
		return err
	}
	strSqlo := `delete from tbl_api_other_info where ApiId=?`
	err = this.Exec(tx, strSqlo, apiId)
	if err != nil {
		return err
	}
	strSql := `delete from tbl_api_basic_info where ApiId=?`
	err = this.Exec(tx, strSql, apiId)
	if err != nil {
		return err
	}
	return err
}

func (this *OscoreApiDB) QueryApiBasicInfoByOntId(tx *sqlx.Tx, ontId string, apiState []int32, pageNum, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	var err error
	sqlStrArr := make([]string, len(apiState))
	for i, state := range apiState {
		sqlStrArr[i] = fmt.Sprintf("%d", state)
	}
	strState := strings.Join(sqlStrArr, ",")

	strSql := fmt.Sprintf("select * from tbl_api_basic_info where OntId=? and ApiState in (%s) order by CreateTime desc limit ?,?", strState)
	fmt.Printf("%s", strSql)

	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize

	res := make([]*tables.ApiBasicInfo, 0)
	err = this.Select(tx, &res, strSql, ontId, start, pageSize)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getUint32InSql(args []uint32) string {
	sqlStrArr := make([]string, len(args))
	for i, state := range args {
		sqlStrArr[i] = fmt.Sprintf("%d", state)
	}
	strState := strings.Join(sqlStrArr, ",")
	return strState
}

func getStrSate(apiState []int32) string {
	sqlStrArr := make([]string, len(apiState))
	for i, state := range apiState {
		sqlStrArr[i] = fmt.Sprintf("%d", state)
	}
	strState := strings.Join(sqlStrArr, ",")
	return strState
}

func (this *OscoreApiDB) QueryApiBasicInfoByOscoreUrlKeyStates(tx *sqlx.Tx, urlkey string, apiState []int32) (*tables.ApiBasicInfo, error) {
	var err error
	strState := getStrSate(apiState)
	strSql := fmt.Sprintf("select * from tbl_api_basic_info where ApiOscoreUrlKey=? and ApiState in (%s)", strState)

	// here must only one. so use get not select.
	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, urlkey)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByOscoreUrlKey(tx *sqlx.Tx, urlkey string, apiState int32) (*tables.ApiBasicInfo, error) {
	var err error
	strSql := `select * from tbl_api_basic_info where ApiOscoreUrlKey=? and ApiState=?`
	info := &tables.ApiBasicInfo{}
	err = this.Get(tx, info, strSql, urlkey, apiState)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (this *OscoreApiDB) SearchApi(tx *sqlx.Tx) (map[string][]*tables.ApiBasicInfo, error) {
	res := make(map[string][]*tables.ApiBasicInfo)
	newestApi := make([]*tables.ApiBasicInfo, 0)
	hottestApi := make([]*tables.ApiBasicInfo, 0)
	freeApi := make([]*tables.ApiBasicInfo, 0)
	strNew := "select * from tbl_api_basic_info where ApiState=? order by CreateTime desc limit ?"
	strHot := "select * from tbl_api_basic_info where ApiState=? order by InvokeFrequency desc limit ?"
	strFree := "select * from tbl_api_basic_info where Price='0' and ApiState=? limit ?"
	err := this.Select(tx, &newestApi, strNew, tables.API_STATE_BUILTIN, 10)
	if err != nil {
		return nil, err
	}
	res["newest"] = newestApi

	err = this.Select(tx, &hottestApi, strHot, tables.API_STATE_BUILTIN, 10)
	if err != nil {
		return nil, err
	}
	res["hottest"] = hottestApi

	err = this.Select(tx, &freeApi, strFree, tables.API_STATE_BUILTIN, 10)
	if err != nil {
		return nil, err
	}

	res["free"] = freeApi
	return res, nil
}

func (this *OscoreApiDB) SearchFreeApi(tx *sqlx.Tx, start, pageSize int) (map[string]interface{}, error) {
	strSql := "select * from tbl_api_basic_info where Price='0' and ApiState=? order by CreateTime desc limit ?,?"
	infos := make([]*tables.ApiBasicInfo, 0)
	err := this.Select(tx, &infos, strSql, tables.API_STATE_BUILTIN, start, pageSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": len(infos),
		"list":  infos,
	}, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByCategoryId(tx *sqlx.Tx, categoryId, start, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	var strSql string
	var err error
	res := make([]*tables.ApiBasicInfo, 0)

	if categoryId != oscoreconfig.CategoryAllId {
		strSql = `select * from tbl_api_basic_info where ApiState=? and ApiId in (select ApiId from tbl_api_tag where TagId=(select id from tbl_tag where category_id=?)) order by CreateTime desc limit ?, ?`
		err = this.Select(tx, &res, strSql, tables.API_STATE_BUILTIN, categoryId, start, pageSize)
	} else {
		strSql = `select * from tbl_api_basic_info where ApiState=? and ApiId in (select ApiId from tbl_api_tag where TagId in (select id from tbl_tag)) order by CreateTime desc limit ?, ?`
		err = this.Select(tx, &res, strSql, tables.API_STATE_BUILTIN, start, pageSize)
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByPage(pageNum, pageSize uint32, apiState int32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize

	strSql := `select * from tbl_api_basic_info where ApiState=? limit ?, ?`
	res := make([]*tables.ApiBasicInfo, 0)
	err := this.DB.Select(&res, strSql, apiState, start, pageSize)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoCount(tx *sqlx.Tx, apiState int32) (uint64, error) {
	strSql := `select count(*) from tbl_api_basic_info where ApiState=?`
	var count uint64
	err := this.Get(tx, &count, strSql, apiState)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoOntIdCount(tx *sqlx.Tx, OntId string, apiState []int32) (uint64, error) {
	strState := getStrSate(apiState)
	strSql := fmt.Sprintf("select count(*) from tbl_api_basic_info where OntId=? and ApiState in (%s)", strState)

	var count uint64
	err := this.Get(tx, &count, strSql, OntId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (this *OscoreApiDB) SearchApiByKey(key string, apiState int32) ([]*tables.ApiBasicInfo, error) {
	k := "%" + key + "%"
	strSql := `select * from tbl_api_basic_info where ApiState=? and (ApiDesc like ? or Title like ? or ApiId in (select ApiId from tbl_api_tag where TagId=(select id from tbl_tag where name=?))) limit 30`

	infos := make([]*tables.ApiBasicInfo, 0)
	err := this.DB.Select(&infos, strSql, apiState, k, k, key)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (this *OscoreApiDB) InsertRequestParam(tx *sqlx.Tx, params []*tables.RequestParam) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		var require int32
		if param.Required {
			require = 1
			if param.Note == "" {
				return fmt.Errorf("%s param is requre. can not empty exzample value", param.ParamName)
			}
		} else {
			require = 0
		}
		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%d','%d','%s','%s','%s','%d')", param.ApiId, param.ParamName, require, param.ParamWhere, param.ParamType, param.Note, param.ValueDesc, param.ParamTag)
	}
	strSql := `insert into tbl_request_param (ApiId,ParamName,Required,ParamWhere,ParamType,Note,ValueDesc,ParamTag) values`
	strSql += strings.Join(sqlStrArr, ",")
	err := this.Exec(tx, strSql)
	return err
}

func (this *OscoreApiDB) QueryRequestParamByApiId(tx *sqlx.Tx, apiId uint32, paramTag int32) ([]*tables.RequestParam, error) {
	strSql := `select * from tbl_request_param where ApiId=? and ParamTag=?`
	params := make([]*tables.RequestParam, 0)
	err := this.Select(tx, &params, strSql, apiId, paramTag)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func (this *OscoreApiDB) QueryReferRequestParamByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.RequestParam, error) {
	strSql := `select * from tbl_request_param where ApiId=?`
	params := make([]*tables.RequestParam, 0)
	err := this.Select(tx, &params, strSql, apiId)
	if err != nil {
		return nil, err
	}
	return params, nil
}

// unit test none.
func (this *OscoreApiDB) InsertErrorCode(tx *sqlx.Tx, params []*tables.ErrorCode) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		sqlStrArr[i] = fmt.Sprintf("('%d','%s')",
			param.ErrorCode, param.ErrorDesc)
	}
	strSql := `insert into tbl_error_code (ErrorCode,ErrorDesc) values`
	strSql += strings.Join(sqlStrArr, ",")
	err := this.Exec(tx, strSql)
	return err
}

// unit test none.
func (this *OscoreApiDB) QueryErrorCode(tx *sqlx.Tx) ([]*tables.ErrorCode, error) {
	strSql := `select * from tbl_error_code`
	params := make([]*tables.ErrorCode, 0)
	err := this.Select(tx, &params, strSql)
	return params, err
}

func (this *OscoreApiDB) InsertSpecifications(tx *sqlx.Tx, params []*tables.Specifications) error {
	if len(params) == 0 {
		return nil
	}
	sqlStrArr := make([]string, len(params))
	for i, param := range params {
		// if not specified. it will consider as counter.
		if param.SpecType == 0 {
			param.SpecType = tables.SPEC_TYPE_COUNT
		}

		if param.SpecType != tables.SPEC_TYPE_COUNT && param.SpecType != tables.SPEC_TYPE_DURATION {
			return errors.New("error spec type.")
		}

		if param.SpecType == tables.SPEC_TYPE_DURATION {
			if param.Amount != 0 {
				return errors.New("Duration specifications type can not specifiy amount.")
			}

			if param.EffectiveDuration == 0 {
				return errors.New("Duration specifications must specifiy effectiveDuration.")
			}
		} else if param.SpecType == tables.SPEC_TYPE_COUNT {
			if param.Amount == 0 {
				return errors.New("Count specifications must specifiy count.")
			}

			if param.EffectiveDuration != 0 {
				return errors.New("Count specifications can not specifiy effectiveDuration.")
			}
		}

		sqlStrArr[i] = fmt.Sprintf("('%d','%s','%d','%d','%d')",
			param.ApiId, param.Price, param.Amount, param.EffectiveDuration, param.SpecType)
		if param.EffectiveDuration > 120 || param.EffectiveDuration < 0 {
			return fmt.Errorf("can not over ten years or less than zero.")
		}
	}

	strSql := `insert into tbl_specifications (ApiId,Price,Amount,EffectiveDuration,SpecType) values`
	strSql += strings.Join(sqlStrArr, ",")
	err := this.Exec(tx, strSql)
	return err
}

func (this *OscoreApiDB) QuerySpecificationsById(tx *sqlx.Tx, id uint32) (*tables.Specifications, error) {
	strSql := `select * from tbl_specifications where Id=?`
	ss := &tables.Specifications{}
	err := this.Get(tx, ss, strSql, id)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (this *OscoreApiDB) QuerySpecificationsByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.Specifications, error) {
	strSql := `select * from tbl_specifications where ApiId=? order by Amount`
	ss := make([]*tables.Specifications, 0)
	err := this.Select(tx, &ss, strSql, apiId)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (this *OscoreApiDB) QuerySpecificationsByApiIdSpecType(tx *sqlx.Tx, apiId uint32, specType int32) ([]*tables.Specifications, error) {
	var strSql string
	switch specType {
	case tables.SPEC_TYPE_COUNT:
		strSql = `select * from tbl_specifications where ApiId=? and SpecType=? order by Amount`
	case tables.SPEC_TYPE_DURATION:
		strSql = `select * from tbl_specifications where ApiId=? and SpecType=? order by EffectiveDuration`
	default:
		return nil, fmt.Errorf("QuerySpecificationsByApiIdSpecType error spec type %d", specType)
	}

	ss := make([]*tables.Specifications, 0)
	err := this.Select(tx, &ss, strSql, apiId, specType)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (this *OscoreApiDB) GetMinPriceSpecOfApiId(tx *sqlx.Tx, apiId uint32) (*tables.Specifications, error) {
	specsCount, _ := this.QuerySpecificationsByApiIdSpecType(tx, apiId, tables.SPEC_TYPE_COUNT)
	specsDuration, _ := this.QuerySpecificationsByApiIdSpecType(tx, apiId, tables.SPEC_TYPE_DURATION)

	if len(specsCount) == 0 && len(specsDuration) == 0 {
		return nil, errors.New("api spec both noting")
	}

	if len(specsCount) == 0 {
		return specsDuration[0], nil
	} else if len(specsDuration) == 0 {
		return specsCount[0], nil
	} else {
		specC := specsCount[0]
		price := utils.ToIntByPrecise(specC.Price, oscoreconfig.ONG_DECIMALS)
		specifications := new(big.Int).SetUint64(uint64(specC.Amount))
		amountC := new(big.Int).Mul(price, specifications)

		specD := specsDuration[0]
		amountD := utils.ToIntByPrecise(specD.Price, oscoreconfig.ONG_DECIMALS)

		res := amountC.Cmp(amountD)
		if res == 1 {
			return specD, nil
		} else {
			return specC, nil
		}
	}
}

//dependent on orderId.
func (this *OscoreApiDB) InsertApiKey(tx *sqlx.Tx, key *tables.APIKey) error {
	strSql := `insert into tbl_api_key (ApiKey,OrderId, ApiId, RequestLimit, UsedNum, OntId, UserId, OutDate, ApiKeyType,CreateTime) values (?,?,?,?,?,?,?,?,?,?)`
	err := this.Exec(tx, strSql, key.ApiKey, key.OrderId, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId, key.UserId, key.OutDate, key.ApiKeyType, key.CreateTime)
	return err
}

//dependent on orderId. use default.
func (this *OscoreApiDB) InsertApiTestKey(tx *sqlx.Tx, key *tables.APIKey) error {
	strSql := `insert into tbl_api_test_key (ApiKey, ApiId, RequestLimit, UsedNum, OntId, UserId, ApiKeyType,CreateTime) values (?,?,?,?,?,?,?,?)`
	err := this.Exec(tx, strSql, key.ApiKey, key.ApiId, key.RequestLimit, key.UsedNum, key.OntId, key.UserId, key.ApiKeyType, key.CreateTime)
	return err
}

func (this *OscoreApiDB) QueryInvokeFreByApiId(tx *sqlx.Tx, apiId uint32) (uint64, error) {
	var freq uint64
	strSql := `select InvokeFrequency from tbl_api_basic_info where ApiId =?`
	err := this.Get(tx, &freq, strSql, apiId)
	if err != nil {
		return 0, err
	}

	return freq, nil
}

func (this *OscoreApiDB) QueryApiKeyByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.APIKey, error) {
	strSql := `select * from tbl_api_key where ApiId=?`

	res := make([]*tables.APIKey, 0)
	err := this.Select(tx, &res, strSql, apiId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryApiKeyByApiKey(tx *sqlx.Tx, apiKey string) (*tables.APIKey, error) {
	return this.queryApiKey(tx, apiKey, "")
}
func (this *OscoreApiDB) QueryApiKeyByOrderId(tx *sqlx.Tx, orderId string) (*tables.APIKey, error) {
	return this.queryApiKey(tx, "", orderId)
}

func (this *OscoreApiDB) QueryApiTestKeyByOntidAndApiId(tx *sqlx.Tx, ontid string, apiId uint32) (*tables.APIKey, error) {
	strSql := "select * from tbl_api_test_key where OntId=? and ApiId=?"
	key := &tables.APIKey{}
	err := this.Get(tx, key, strSql, ontid, apiId)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (this *OscoreApiDB) queryApiKey(tx *sqlx.Tx, key, orderId string) (*tables.APIKey, error) {
	var strSql string
	var where string
	if key != "" {
		if common.IsTestKey(key) {
			strSql = "select * from tbl_api_test_key where ApiKey=?"
		} else {
			strSql = "select * from tbl_api_key where ApiKey=?"
		}

		where = key
	} else if orderId != "" {
		strSql = "select * from tbl_api_key where OrderId=?"
		where = orderId
	} else {
		return nil, errors.New("both queryApiKeykey and order can not null.")
	}

	k := &tables.APIKey{}
	err := this.Get(tx, k, strSql, where)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (this *OscoreApiDB) QueryApiKeyByOntId(tx *sqlx.Tx, ontId string, start, pageSize int) ([]*tables.APIKey, error) {
	res := make([]*tables.APIKey, 0)
	strSql := "select * from tbl_api_key where OntId=? order by CreateTime desc limit ?,?"
	err := this.Select(tx, &res, strSql, ontId, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *OscoreApiDB) QueryApiKeyCountByOntId(tx *sqlx.Tx, ontId string) (uint64, error) {
	strSql := `select count(*) from tbl_api_key where OntId=?`
	var count uint64
	err := this.Get(tx, &count, strSql, ontId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (this *OscoreApiDB) UpdateApiKeyReqLimit(tx *sqlx.Tx, apiKey string, requestLimit uint64, outDate int64) error {
	strSql := "update tbl_api_key set RequestLimit=?,OutDate=? where ApiKey=?"
	return this.Exec(tx, strSql, requestLimit, outDate, apiKey)
}

func (this *OscoreApiDB) UpdateApiKeyInvokeFre(tx *sqlx.Tx, apiKey string, apiId uint32, usedNum, invokeFre uint64) error {
	var strSql string
	var err error
	if common.IsTestKey(apiKey) {
		// here no need update invokefreq.
		strSql = "update tbl_api_test_key set UsedNum=? where ApiKey=?"
		err = this.Exec(tx, strSql, usedNum, apiKey)
	} else {
		strSql = "update tbl_api_key k,tbl_api_basic_info i set k.UsedNum=?,i.InvokeFrequency=? where k.ApiKey=? and i.ApiId=?"
		err = this.Exec(tx, strSql, usedNum, invokeFre, apiKey, apiId)
	}

	return err
}

func (this *OscoreApiDB) ApiBasicUpateApiState(tx *sqlx.Tx, apiState int32, apiId uint32, oscoreUrlKey string) error {
	strSql := "update tbl_api_basic_info set ApiState=? where ApiId=? and ApiOscoreUrlKey=?"
	if apiState >= tables.API_STATE_LAST {
		return errors.New("wrong api state info.")
	}
	err := this.Exec(tx, strSql, apiState, apiId, oscoreUrlKey)
	return err
}

func (this *OscoreApiDB) ApiBasicUpateApiNotifyDelete(tx *sqlx.Tx, notifyDelete int32, apiId uint32) error {
	strSql := "update tbl_api_basic_info set NotifyDelete=? where ApiId=?"
	err := this.Exec(tx, strSql, notifyDelete, apiId)
	return err
}

func (this *OscoreApiDB) ApiBasicUpateApiStateByOntIdApiId(tx *sqlx.Tx, ontId string, apiState int32, apiId uint32, disableOrderTime int64) error {
	strSql := "update tbl_api_basic_info set ApiState=?,UpdateTime=? where OntId=? and ApiId=?"
	if apiState >= tables.API_STATE_LAST {
		return errors.New("wrong api state info.")
	}
	err := this.Exec(tx, strSql, apiState, disableOrderTime, ontId, apiId)
	return err
}

func (this *OscoreApiDB) InsertApiCryptoInfo(tx *sqlx.Tx, apiCryptoInfo *tables.ApiCryptoInfo) error {
	valueStr := fmt.Sprintf("('%d','%d','%d','%s','%s')", apiCryptoInfo.ApiId, apiCryptoInfo.CryptoType, apiCryptoInfo.CryptoWhere, apiCryptoInfo.CryptoKey, apiCryptoInfo.CryptoValue)
	strSql := `insert into tbl_api_crypto_info (ApiId,CryptoType,CryptoWhere,CryptoKey,CryptoValue) values` + valueStr
	err := this.Exec(tx, strSql)
	return err
}

func (this *OscoreApiDB) QueryApiCryptoInfoByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.ApiCryptoInfo, error) {
	strSql := `select * from tbl_api_crypto_info where ApiId=?`
	ss := make([]*tables.ApiCryptoInfo, 0)
	err := this.Select(tx, &ss, strSql, apiId)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func (this *OscoreApiDB) InsertNotifycations(tx *sqlx.Tx, notifies []*tables.Notification) error {
	sqlStrArr := make([]string, len(notifies))
	for i, info := range notifies {
		sqlStrArr[i] = fmt.Sprintf("('%s','%s','%d','%s','%s','%s','%s','%s','%s','%d')",
			info.Id, info.UserId, info.Type, info.TitleEn, info.TitleZh, info.ContentEn, info.ContentZh, info.KeyWord, info.BusinessId, info.IsRead)
	}

	strSql := `insert into tbl_notification (id,user_id,type,title_en,title_zh,content_en,content_zh,keyword,business_id,is_read) values`
	strSql += strings.Join(sqlStrArr, ",")
	return this.Exec(tx, strSql)
}
