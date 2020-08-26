package common

import (
	"github.com/ontio/oscore/models/tables"
)

type GetOrderResponse struct {
	Result   string `json:"result"`
	UserName string `json:"userName"`
	OntId    string `json:"ontId"`
}

type ApiAttachMent struct {
	*tables.ApiBasicInfo
	SpecId     uint32 `json:"specId"`
	SpecPrice  string `json:"specPrice"`
	SpecAmount uint64 `json:"specAmount"`
	// month num. to ensure the db do not overflow. do not use int to insert.
	SpecEffectiveDuration int32 `json:"specEffectiveDuration"`
	SpecType              int32 `json:"specType"`
}

type ResponseParamDesc struct {
	ParamName string `json:"paramName"`
	ParamType string `json:"paramType"`
	ParamDesc string `json:"paramDesc"`
}

type PublishErrorCode struct {
	Code string `json:"code"`
	Desc string `json:"description"`
}

type ApiDetailResponse struct {
	ApiId               uint32                   `json:"apiId"`
	Mark                string                   `json:"mark"`
	ResponseParam       string                   `json:"responseParam"`
	ResponseType        string                   `json:"responseType"`
	ResponseExample     string                   `json:"responseExample"`
	DataDesc            string                   `json:"dataDesc"`
	DataSource          string                   `json:"dataSource"`
	ApplicationScenario string                   `json:"applicationScenario"`
	RequestParams       []*tables.RequestParam   `json:"requestParams"`
	ErrorCodes          []*PublishErrorCode      `json:"errorCodes"`
	Specifications      []*tables.Specifications `json:"specifications"`
	ResponseParamDescs  []*ResponseParamDesc     `json:"responseParamDescs"`
	ApiBasicInfo        *tables.ApiBasicInfo     `json:"apiBasicInfo"`
}

type OrderResult struct {
	Title          string                   `json:"title"`
	Total          uint32                   `json:"total"`
	OrderId        string                   `json:"orderId"`
	Amount         string                   `json:"amount"`
	CreateTime     int64                    `json:"createTime"`
	PayTime        int64                    `json:"payTime"`
	ApiId          uint32                   `json:"apiId"`
	ApiUrl         string                   `json:"apiUrl"`
	State          int32                    `json:"state"`
	ApiKey         string                   `json:"apiKey"`
	Price          string                   `json:"price"`
	TxHash         string                   `json:"txHash"`
	Type           string                   `json:"type"`
	ApiState       int32                    `json:"apiState"`
	RequestLimit   uint64                   `json:"requestLimit"`
	UsedNum        uint64                   `json:"usedNum"`
	Spec           *tables.Specifications   `json:"spec"`
	Specifications []*tables.Specifications `json:"specifications"`
	Comment        *tables.OrderApiComment  `json:"comment"`
}

type DataProcessOrderResult struct {
	Title     string                  `json:"title"`
	OrderId   string                  `json:"orderId"`
	OrderTime int64                   `json:"orderTime"`
	PayTime   int64                   `json:"payTime"`
	ApiId     uint32                  `json:"apiId"`
	State     int32                   `json:"state"`
	Price     string                  `json:"price"`
	TxHash    string                  `json:"txHash"`
	OrderKind uint32                  `json:"orderKind"`
	Request   string                  `json:"request"`
	Result    string                  `json:"result"`
	Type      string                  `json:"type"`
	Icon      string                  `json:"icon"`
	Comment   *tables.OrderApiComment `json:"comment"`
}

type OrderResultResponse struct {
	Order  *tables.Order
	ApiKey *tables.APIKey
}

type OrderDetailResponse struct {
	Res interface{}
}

type WetherOrderDetail struct {
	TargetDate int64                `json:"targetDate"`
	Location   *tables.Location     `json:"location"`
	ToolBox    *tables.ToolBox      `json:"toolBoxId"`
	ApiSource  *tables.ApiBasicInfo `json:"apiSourceId"`
	Algorithm  *tables.Algorithm    `json:"algorithmId"`
	Env        *tables.Env          `json:"envId"`
	Result     string               `json:"result"`
	State      int32                `json:"state"`
}

type CommentDetail struct {
	tables.OrderApiComment
	Title string `json:"title"`
}

type ApiKeysDetail struct {
	tables.APIKey
	WillInvalid    bool                     `json:"willInvalid"`
	Icon           string                   `json:"icon"`
	Title          string                   `json:"title"`
	SpecType       int32                    `json:"specType"`
	ApiState       int32                    `json:"apiState"`
	Specifications []*tables.Specifications `json:"specifications"`
}
