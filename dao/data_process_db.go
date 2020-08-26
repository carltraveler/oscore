package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
)

func (this *OscoreApiDB) QueryLocationOfCountryCity(tx *sqlx.Tx, country string) ([]*tables.Location, error) {
	var err error
	res := make([]*tables.Location, 0)
	if country != "ALL" {
		strSql := `select * from tbl_country_city where Country=? order by City`

		err = this.Select(tx, &res, strSql, country)
		if err != nil {
			return nil, err
		}
	} else {
		strSql := `select DISTINCT Country from tbl_country_city order by Country`

		err = this.Select(tx, &res, strSql)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (this *OscoreApiDB) QueryApiBasicInfoByApiTypeKind(tx *sqlx.Tx, apiType string, apiKind int32, apiState int32) ([]*tables.ApiBasicInfo, error) {
	var err error
	res := make([]*tables.ApiBasicInfo, 0)
	StrSql := `select * from tbl_api_basic_info where ApiType=? and ApiKind=? and ApiState=?`
	err = this.Select(tx, &res, StrSql, apiType, apiKind, apiState)
	if err != nil {
		return nil, err
	}
	log.Debugf("apiType: %s, %d, %d, len(res): %d", apiType, apiKind, apiState, len(res))

	return res, nil
}

func (this *OscoreApiDB) QueryApiAlgorithmsByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.ApiAlgorithm, error) {
	var err error
	StrSql := `select * from tbl_api_algorithm where ApiId=? and State=1`
	res := make([]*tables.ApiAlgorithm, 0)
	err = this.Select(tx, &res, StrSql, apiId)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (this *OscoreApiDB) QueryAlgorithmEnvByAlgorithmId(tx *sqlx.Tx, algorithmId uint32) ([]*tables.AlgorithmEnv, error) {
	var err error
	StrSql := `select * from tbl_algorithm_env where AlgorithmId=? and State=1`
	res := make([]*tables.AlgorithmEnv, 0)
	err = this.Select(tx, &res, StrSql, algorithmId)
	if err != nil {
		return nil, err
	}
	log.Debugf("algorithmId: %d, len(res):%d", algorithmId, len(res))

	return res, nil

}

func (this *OscoreApiDB) QueryAlgorithmById(tx *sqlx.Tx, id uint32) (*tables.Algorithm, error) {
	var err error
	var res tables.Algorithm
	StrSql := `select * from tbl_algorithm where Id=? and State=1`
	err = this.Get(tx, &res, StrSql, id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (this *OscoreApiDB) QueryEnvById(tx *sqlx.Tx, id uint32) (*tables.Env, error) {
	var err error
	var res tables.Env
	StrSql := `select * from tbl_env where Id=? and State=1`
	err = this.Get(tx, &res, StrSql, id)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *OscoreApiDB) QueryToolBoxById(tx *sqlx.Tx, id uint32) (*tables.ToolBox, error) {
	var err error
	var res tables.ToolBox
	StrSql := `select * from tbl_tool_box where Id=? and State=1`
	err = this.Get(tx, &res, StrSql, id)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *OscoreApiDB) SearchToolBoxByKey(tx *sqlx.Tx, key string) ([]*tables.ToolBox, error) {
	k := "%" + key + "%"
	strSql := `select * from tbl_tool_box where ToolBoxDesc like ? or Title like ? or Id in (select ToolBoxId from tbl_tool_box_tag where TagId=(select id from tbl_tag where name=?)) limit 30`
	infos := make([]*tables.ToolBox, 0)
	err := this.Select(tx, &infos, strSql, k, k, key)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (this *OscoreApiDB) QueryToolBoxByCategoryId(tx *sqlx.Tx, categoryId, pageNum, pageSize uint32) ([]*tables.ToolBox, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize
	log.Debugf("QueryToolBoxByCategoryId %d, %d, %d", categoryId, start, pageSize)

	var strSql string
	var err error
	res := make([]*tables.ToolBox, 0)
	if categoryId != oscoreconfig.CategoryAllId {
		strSql = `select * from tbl_tool_box where Id in (select ToolBoxId from tbl_tool_box_tag where TagId=(select id from tbl_tag where category_id=?)) limit ?, ?`
		err = this.Select(tx, &res, strSql, categoryId, start, pageSize)
	} else {
		strSql = `select * from tbl_tool_box where Id in (select ToolBoxId from tbl_tool_box_tag where TagId in (select id from tbl_tag)) limit ?, ?`
		err = this.Select(tx, &res, strSql, start, pageSize)
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryToolBoxAll(tx *sqlx.Tx) ([]*tables.ToolBox, error) {
	var err error
	res := make([]*tables.ToolBox, 0)
	StrSql := `select * from tbl_tool_box where State=1`
	err = this.Select(tx, &res, StrSql)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryHeaderValueByApiId(tx *sqlx.Tx, apiId uint32) ([]*tables.ApiHeadValues, error) {
	var err error
	res := make([]*tables.ApiHeadValues, 0)
	StrSql := `select * from tbl_api_header_values where ApiId=?`
	err = this.Select(tx, &res, StrSql, apiId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
