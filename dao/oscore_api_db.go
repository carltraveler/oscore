package dao

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"time"
)

var DefOscoreApiDB *OscoreApiDB

type OscoreApiDB struct {
	DB *sqlx.DB
}

func NewOscoreApiDB(dbConfig *oscoreconfig.DBConfig) (*OscoreApiDB, error) {
	dbx, dberr := sqlx.Open("mysql",
		dbConfig.ProjectDBUser+
			":"+dbConfig.ProjectDBPassword+
			"@tcp("+dbConfig.ProjectDBUrl+
			")/"+dbConfig.ProjectDBName+
			"?charset=utf8&parseTime=true")
	if dberr != nil {
		return nil, dberr
	}

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	err := dbx.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	dbx.SetMaxIdleConns(256)

	return &OscoreApiDB{
		DB: dbx,
	}, nil
}

func (this *OscoreApiDB) Exec(tx *sqlx.Tx, query string, args ...interface{}) error {
	var err error
	if tx != nil {
		_, err = tx.Exec(query, args...)
	} else {
		_, err = this.DB.Exec(query, args...)
	}
	return err
}

func (this *OscoreApiDB) Select(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	var err error
	if tx != nil {
		err = tx.Select(dest, query, args...)
	} else {
		err = this.DB.Select(dest, query, args...)
	}
	return err
}

func (this *OscoreApiDB) Get(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	var err error
	if tx != nil {
		err = tx.Get(dest, query, args...)
	} else {
		err = this.DB.Get(dest, query, args...)
	}
	return err
}

/////////////////////
func (this *OscoreApiDB) InsertApiTag(tx *sqlx.Tx, apiTag *tables.ApiTag) error {
	sqlStr := `insert into tbl_api_tag (ApiId, TagId, State) values (?,?,?)`
	err := this.Exec(tx, sqlStr, apiTag.ApiId, apiTag.TagId, apiTag.State)
	return err
}

func (this *OscoreApiDB) QueryTagByNameId(tx *sqlx.Tx, categoryId uint32, name string) (*tables.Tag, error) {
	var res tables.Tag
	strSql := `select * from tbl_tag where category_id=? and name=?`
	err := this.Get(tx, &res, strSql, categoryId, name)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *OscoreApiDB) QueryCategoryByName(tx *sqlx.Tx, NameEn string) (*tables.Category, error) {
	var res tables.Category
	sqlStr := `select * from tbl_category where name_en=?`
	err := this.Get(tx, &res, sqlStr, NameEn)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *OscoreApiDB) QueryCategoryById(tx *sqlx.Tx, categoryId uint32) (*tables.Category, error) {
	var res tables.Category
	sqlStr := `select * from tbl_category where id=?`
	err := this.Get(tx, &res, sqlStr, categoryId)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *OscoreApiDB) Close() error {
	return this.DB.Close()
}
