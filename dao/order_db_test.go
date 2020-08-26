package dao

import (
	"testing"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/stretchr/testify/assert"
	"time"
)

var TestDB *OscoreApiDB

func TestMain(m *testing.M) {
	fmt.Println("begin test db.")
	var err error
	TestDB, err = NewOscoreApiDB(oscoreconfig.DefDBConfigMap[oscoreconfig.NETWORK_ID_TRAVIS_NET])
	if err != nil {
		panic(err)
	}
	TestDB.ClearAll()
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "mytestasd",
		ApiProvider:     "",
		ApiUrl:          "",
		Price:           "",
		ApiDesc:         "",
		ApiState:        tables.API_STATE_BUILTIN,
		Specifications:  1,
		Popularity:      0,
		Delay:           0,
		SuccessRate:     0,
		InvokeFrequency: 0,
	}
	err = TestDB.InsertApiBasicInfo(nil, []*tables.ApiBasicInfo{info})
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}
	_, err = TestDB.QueryApiBasicInfoByPage(0, 1, tables.API_STATE_BUILTIN)
	fmt.Printf("Test order begin.")
	if err != nil {
		panic(err)
	}
	m.Run()
	fmt.Printf("all done.")
}

func TestOrderDB_InsertOrder(t *testing.T) {
	tt := time.Now().Unix()
	InsertTestApi()

	basic, err := TestDB.QueryApiBasicInfoByPage(0, 1, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	orderId := "abcedkfy"
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.InsertOrder(nil, order)
	assert.Nil(t, err)
	orderFromDb, err := TestDB.QueryOrderByOrderId(nil, orderId)
	assert.Nil(t, err)
	assert.Equal(t, orderFromDb.OrderId, orderId)

	err = TestDB.UpdateOrderStatus(nil, orderId, oscoreconfig.Canceled)
	assert.Nil(t, err)

	orderId2 := "abcedkfyfgt"
	order2 := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId2,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.InsertOrder(nil, order2)
	assert.Nil(t, err)

	code := &tables.QrCode{
		QrCodeId: "qbcdab",
		OrderId:  orderId,
		Exp:      tt,
	}
	err = TestDB.InsertQrCode(nil, code)
	assert.Nil(t, err)
	code = &tables.QrCode{
		QrCodeId: "qbcdabc",
		OrderId:  orderId,
		Exp:      tt,
	}
	err = TestDB.InsertQrCode(nil, code)
	assert.Nil(t, err)
	code, err = TestDB.QueryQrCodeByOrderId(nil, orderId)
	assert.Nil(t, err)
	fmt.Println(code)
	fmt.Println(code.QrCodeId)
	_, err = TestDB.QueryQrCodeResultByQrCodeId(nil, code.QrCodeId)
	assert.Nil(t, err)

	err = TestDB.DeleteQrCodeByOrderId(nil, orderId)
	assert.Nil(t, err)
	err = TestDB.DeleteOrderByOrderId(nil, orderId)
	assert.Nil(t, err)
	err = TestDB.DeleteOrderByOrderId(nil, orderId2)
	assert.Nil(t, err)
	fmt.Printf("TestOrderDB_InsertOrder done.")
	TestDB.ClearApiBasicDB()
}

func TestOrderDB_QueryOrderByOrderId(t *testing.T) {
	TestDB.ClearOrderDB()
	tt := time.Now().Unix()
	InsertTestApi()
	basic, err := TestDB.QueryApiBasicInfoByPage(0, 1, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	orderId := "abcedkfy"
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.InsertOrder(nil, order)
	assert.Nil(t, err)
	err = TestDB.UpdateTxInfoByOrderId(nil, orderId, "", "", oscoreconfig.Canceled, 0)
	assert.Nil(t, err)

	order, err = TestDB.QueryOrderByOrderId(nil, orderId)
	assert.Nil(t, err)
	assert.Equal(t, order.OrderId, orderId)
	assert.Equal(t, order.OrderStatus, oscoreconfig.Canceled)

	orders, err := TestDB.QueryOrderByPage(nil, 0, 1, "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(orders))

	err = TestDB.UpdateOrderStatus(nil, orderId, oscoreconfig.Completed)
	assert.Nil(t, err)
	status, err := TestDB.QueryOrderStatusByOrderId(nil, orderId)
	assert.Nil(t, err)
	assert.Equal(t, oscoreconfig.Completed, status)

	TestDB.ClearOrderDB()
	TestDB.ClearApiBasicDB()
	fmt.Printf("TestOrderDB_QueryOrderByOrderId done.")
}

func TestOrderDB_DeleteOrderByOrderId(t *testing.T) {
	orderId := "abcedkfyfgtghj"
	InsertTestApi()
	basic, err := TestDB.QueryApiBasicInfoByPage(0, 1, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	tt := time.Now().Unix()
	order2 := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.InsertOrder(nil, order2)
	assert.Nil(t, err)
	err = TestDB.DeleteOrderByOrderId(nil, orderId)
	assert.Nil(t, err)
	TestDB.ClearApiBasicDB()
	fmt.Printf("TestOrderDB_DeleteOrderByOrderId done.")
}

func TestOrderDB_QueryOrderByPage(t *testing.T) {
	order, err := TestDB.QueryOrderByPage(nil, 0, 2, "", "")
	assert.Nil(t, err)
	fmt.Println(order)
	fmt.Printf("TestOrderDB_QueryOrderByPage done.")
}
