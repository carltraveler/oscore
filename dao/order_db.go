package dao

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/oscore/models/tables"
)

const (
	DESC string = "DESC"
	ASC  string = "ASC"
)

func (this *OscoreApiDB) ClearOrderDB() error {
	strSql := "delete from tbl_order"
	_, err := this.DB.Exec(strSql)
	return err
}

func (this *OscoreApiDB) InsertOrder(tx *sqlx.Tx, order *tables.Order) error {
	// use NameExec better.
	strSql := `insert into tbl_order (OrderId,Title, ProductName, OrderType, OrderTime, State,Amount, 
OntId,UserId,UserName,TxHash,Price,ApiId,ToolBoxId,ApiUrl,SpecificationsId,OrderKind,Request,Result,ApiKey) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	err := this.Exec(tx, strSql, order.OrderId, order.Title, order.ProductName, order.OrderType, order.OrderTime, order.State,
		order.Amount, order.OntId, order.UserId, order.UserName, order.TxHash, order.Price, order.ApiId, order.ToolBoxId, order.ApiUrl, order.SpecificationsId, order.OrderKind, order.Request, order.Result, order.ApiKey)
	return err
}

func (this *OscoreApiDB) UpdateTxInfoByOrderId(tx *sqlx.Tx, orderId string, result string, state int32, payTime int64) error {
	strSql := "update tbl_order set Result=?,State=?,PayTime=? where OrderId=?"
	err := this.Exec(tx, strSql, result, state, payTime, orderId)
	return err
}

func (this *OscoreApiDB) UpdateOrderApiKey(tx *sqlx.Tx, txHash, orderId, apiKey string) error {
	strSql := "update tbl_order set TxHash=?,ApiKey=? where OrderId=?"
	err := this.Exec(tx, strSql, txHash, apiKey, orderId)
	return err
}

func (this *OscoreApiDB) QueryOrderStatusByOrderId(tx *sqlx.Tx, orderId string) (int32, error) {
	strSql := `select State from tbl_order where OrderId=?`
	var orderStatus int32
	err := this.Get(tx, &orderStatus, strSql, orderId)
	if err != nil {
		return 0, err
	}
	return orderStatus, nil
}

func (this *OscoreApiDB) QueryOrderByOrderId(tx *sqlx.Tx, orderId string) (*tables.Order, error) {
	strSql := `select * from tbl_order where OrderId=?`
	order := &tables.Order{}
	err := this.Get(tx, order, strSql, orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (this *OscoreApiDB) QueryOrderByQrCodeId(tx *sqlx.Tx, qrCodeId string) (*tables.Order, error) {
	strSql := `select * from tbl_order where OrderId=(select OrderId from tbl_qr_code where  QrCodeId=?)`
	order := &tables.Order{}
	err := this.Get(tx, order, strSql, qrCodeId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (this *OscoreApiDB) QueryOrderSum(tx *sqlx.Tx, ontId string, orderType string) (int, error) {
	strSql := `select count(*) from tbl_order where OntId=? and OrderType=?`
	var sum int
	err := this.Get(tx, &sum, strSql, ontId, orderType)
	if err != nil {
		return 0, nil
	}
	return sum, nil
}

func (this *OscoreApiDB) QueryOrderSumStatus(tx *sqlx.Tx, ontId string, orderType string, state int32) (int, error) {
	strSql := `select count(*) from tbl_order where OntId=? and OrderType=? and State=?`
	var sum int
	err := this.Get(tx, &sum, strSql, ontId, orderType, state)
	if err != nil {
		return 0, nil
	}
	return sum, nil
}

func (this *OscoreApiDB) QueryOrderByPage(tx *sqlx.Tx, start, pageSize int, ontId string, orderType string) ([]*tables.Order, error) {
	strSql := `select * from tbl_order where OntId=? and OrderType=? order by OrderTime desc limit ?, ?`
	res := make([]*tables.Order, 0)
	err := this.Select(tx, &res, strSql, ontId, orderType, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *OscoreApiDB) QueryOrderResultByPage(tx *sqlx.Tx, start, pageSize int, ontId string, orderType string, state int32) ([]*tables.Order, error) {
	strSql := `select * from tbl_order where OntId=? and OrderType=? and State=? order by OrderTime desc limit ?, ?`
	res := make([]*tables.Order, 0)
	err := this.Select(tx, &res, strSql, ontId, orderType, state, start, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (this *OscoreApiDB) QueryOrderRelatedTypeApi(tx *sqlx.Tx, start, pageSize int, orderType string, resApi []*tables.ApiBasicInfo, sorting string, state int32) ([]*tables.Order, int, error) {
	if len(resApi) == 0 {
		tmp := make([]*tables.Order, 0)
		return tmp, 0, nil
	}

	var total int
	apiIds := make([]string, len(resApi))
	for i, api := range resApi {
		apiIds[i] = fmt.Sprintf("%d", api.ApiId)
	}
	apiIdstr := strings.Join(apiIds, ",")
	strSqlCount := fmt.Sprintf("select count(*) from tbl_order where OrderType=? and State=? and ApiId in (%s)", apiIdstr)
	err := this.Get(tx, &total, strSqlCount, orderType, state)
	if err != nil {
		return nil, 0, err
	}

	var strSql string
	if sorting == DESC {
		strSql = fmt.Sprintf("select * from tbl_order where OrderType=? and State=? and ApiId in (%s) order by OrderTime desc limit ?, ?", apiIdstr)
	} else if sorting == ASC {
		strSql = fmt.Sprintf("select * from tbl_order where OrderType=? and State=? and ApiId in (%s) order by OrderTime limit ?, ?", apiIdstr)
	} else {
		return nil, 0, fmt.Errorf("error sorting: %s", sorting)
	}
	res := make([]*tables.Order, 0)
	err = this.Select(tx, &res, strSql, orderType, state, start, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (this *OscoreApiDB) QueryOrderByPageOntIdType(tx *sqlx.Tx, start, pageSize int, orderType, ontId string, sorting string) ([]*tables.Order, int, error) {
	strState := getStrSate([]int32{tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER})
	strSqlApi := fmt.Sprintf("select * from tbl_api_basic_info where OntId=? and ApiState in (%s)", strState)
	resApi := make([]*tables.ApiBasicInfo, 0)
	// only need check the published.
	err := this.Select(tx, &resApi, strSqlApi, ontId)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*tables.Order, 0)
	if len(resApi) == 0 {
		return res, 0, nil
	}

	return this.QueryOrderRelatedTypeApi(tx, start, pageSize, orderType, resApi, sorting, tables.ORDER_STATE_COMPLETE)
}

func (this *OscoreApiDB) QueryOrderByPageOntIdTypeTitle(tx *sqlx.Tx, start, pageSize int, orderType, ontId string, title string, sorting string) ([]*tables.Order, int, error) {
	strState := getStrSate([]int32{tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER})
	strSqlApi := fmt.Sprintf("select * from tbl_api_basic_info where OntId=? and ApiState in (%s) and Title like ?", strState)

	k := "%" + title + "%"
	resApi := make([]*tables.ApiBasicInfo, 0)
	// only need check the published.
	err := this.Select(tx, &resApi, strSqlApi, ontId, k)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*tables.Order, 0)
	if len(resApi) == 0 {
		return res, 0, nil
	}

	return this.QueryOrderRelatedTypeApi(tx, start, pageSize, orderType, resApi, sorting, tables.ORDER_STATE_COMPLETE)
}

func (this *OscoreApiDB) QueryOrderByPageOntIdTypeApiId(tx *sqlx.Tx, start, pageSize int, orderType, ontId string, apiId uint32, sorting string) ([]*tables.Order, int, error) {
	strState := getStrSate([]int32{tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER})
	strSqlApi := fmt.Sprintf("select * from tbl_api_basic_info where OntId=? and ApiState in (%s) and ApiId=?", strState)

	resApi := make([]*tables.ApiBasicInfo, 0)
	// only need check the published.
	err := this.Select(tx, &resApi, strSqlApi, ontId, apiId)
	if err != nil {
		return nil, 0, err
	}
	res := make([]*tables.Order, 0)
	if len(resApi) == 0 {
		return res, 0, nil
	}

	return this.QueryOrderRelatedTypeApi(tx, start, pageSize, orderType, resApi, sorting, tables.ORDER_STATE_COMPLETE)
}

func (this *OscoreApiDB) UpdateOrderStatus(tx *sqlx.Tx, orderId string, state int32) error {
	strSql := "update tbl_order set State=? where OrderId=?"
	err := this.Exec(tx, strSql, state, orderId)
	return err
}

func (this *OscoreApiDB) UpdateOrderStatusSpecIdAmount(tx *sqlx.Tx, orderId string, state int32, specificationsId uint32, amountStr string) error {
	strSql := "update tbl_order set State=?,SpecificationsId=?,Amount=? where OrderId=?"
	err := this.Exec(tx, strSql, state, specificationsId, amountStr, orderId)
	return err
}

func (this *OscoreApiDB) DeleteOrderByOrderId(tx *sqlx.Tx, orderId string) error {
	strSql := `delete from tbl_order where OrderId=?`
	err := this.Exec(tx, strSql, orderId)
	return err
}

func (this *OscoreApiDB) InsertOrderApiCommentByOrderId(tx *sqlx.Tx, comment *tables.OrderApiComment) error {
	strSql := `insert into tbl_order_api_comment (OrderId,ApiId,ToolBoxId,StarNum,Comments,OntId,CommentTime) values (?,?,?,?,?,?,?)`
	err := this.Exec(tx, strSql, comment.OrderId, comment.ApiId, comment.ToolBoxId, comment.StarNum, comment.Comments, comment.OntId, comment.CommentTime)
	return err
}

func (this *OscoreApiDB) QueryOrderApiCommentByOrderId(tx *sqlx.Tx, orderId string) (*tables.OrderApiComment, error) {
	strSql := `select * from tbl_order_api_comment where OrderId=? and State=1`

	res := &tables.OrderApiComment{}
	err := this.Get(tx, res, strSql, orderId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryOrderApiCommentByApiId(tx *sqlx.Tx, start, pageSize int, apiId uint32) ([]*tables.OrderApiComment, int, error) {
	strSqlCount := `select count(*) from tbl_order_api_comment where ApiId=? and State=1`

	var sum int
	err := this.Get(tx, &sum, strSqlCount, apiId)
	if err != nil {
		return nil, 0, nil
	}

	strSql := `select * from tbl_order_api_comment where ApiId=? and State=1 order by CommentTime desc limit ?, ?`

	res := make([]*tables.OrderApiComment, 0)
	err = this.Select(tx, &res, strSql, apiId, start, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return res, sum, nil
}

func (this *OscoreApiDB) DelOrderApiCommentByOrderId(tx *sqlx.Tx, id uint32, ontId string) error {
	strSql := `delete from tbl_order_api_comment where Id=? and OntId=?`
	err := this.Exec(tx, strSql, id, ontId)
	return err
}

func (this *OscoreApiDB) QueryOrderApiCommentByToolBoxId(tx *sqlx.Tx, start, pageSize int, toolBoxId uint32) ([]*tables.OrderApiComment, int, error) {
	strSqlCount := `select count(*) from tbl_order_api_comment where ToolBoxId=? and State=1`

	var sum int
	err := this.Get(tx, &sum, strSqlCount, toolBoxId)
	if err != nil {
		return nil, 0, nil
	}

	strSql := `select * from tbl_order_api_comment where ToolBoxId=? and State=1 order by CommentTime desc limit ?, ?`

	res := make([]*tables.OrderApiComment, 0)
	err = this.Select(tx, &res, strSql, toolBoxId, start, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return res, sum, nil
}

func (this *OscoreApiDB) QueryCommentPage(tx *sqlx.Tx, start, pageSize int) ([]*tables.OrderApiComment, int, error) {
	strSqlCount := `select count(*) from tbl_order_api_comment where State=1`

	var sum int
	err := this.Get(tx, &sum, strSqlCount)
	if err != nil {
		return nil, 0, nil
	}

	strSql := `select * from tbl_order_api_comment where State=1 order by CommentTime desc limit ?, ?`

	res := make([]*tables.OrderApiComment, 0)
	err = this.Select(tx, &res, strSql, start, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return res, sum, nil
}

func (this *OscoreApiDB) QueryCommentById(tx *sqlx.Tx, commentId uint32) (*tables.OrderApiComment, error) {
	strSql := `select * from tbl_order_api_comment where Id=? and State=1`

	res := &tables.OrderApiComment{}
	err := this.Get(tx, res, strSql, commentId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) DelCommentById(tx *sqlx.Tx, commentId uint32) error {
	strSql := `delete from tbl_order_api_comment where Id=?`

	return this.Exec(tx, strSql, commentId)
}

func (this *OscoreApiDB) QueryUserNameByOntId(tx *sqlx.Tx, ontId string) (*tables.UserName, error) {
	strSql := `select * from tbl_user where ont_id=? or cent_ont_id=?`
	res := &tables.UserName{}
	err := this.Get(tx, res, strSql, ontId, ontId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) QueryEnterpriseInfoByUserId(tx *sqlx.Tx, userId string) (*tables.EnterpriseInfo, error) {
	strSql := `select * from tbl_enterprise_info where user_id=?`
	res := &tables.EnterpriseInfo{}
	err := this.Get(tx, res, strSql, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) InsertOntScore(tx *sqlx.Tx, score *tables.OntIdScore) error {
	strSql := `insert into tbl_score (OntId,TotalCommentNum,TotalScore) values (?,?,?)`
	return this.Exec(tx, strSql, score.OntId, score.TotalCommentNum, score.TotalScore)
}

func (this *OscoreApiDB) QueryOntScoreByOntId(tx *sqlx.Tx, ontId string) (*tables.OntIdScore, error) {
	strSql := `select * from tbl_score where OntId=?`
	res := &tables.OntIdScore{}

	err := this.Get(tx, res, strSql, ontId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *OscoreApiDB) UpdateOntIdScore(tx *sqlx.Tx, score *tables.OntIdScore) error {
	strSql := `update tbl_score set TotalCommentNum=?,TotalScore=? where OntId=?`
	return this.Exec(tx, strSql, score.TotalCommentNum, score.TotalScore, score.OntId)
}
