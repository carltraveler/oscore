package tables

import (
	"database/sql"
)

const (
	ORDER_KIND_API                 = 1
	ORDER_KIND_DATA_PROCESS_WETHER = 2
	ORDER_KIND_API_RENEW           = 3
	ORDER_KIND_DATA_LAST           = 4
)

const (
	ORDER_STATE_WAIT_PAYMENT       int32 = 1
	ORDER_STATE_PAYING             int32 = 2
	ORDER_STATE_COMPLETE           int32 = 3
	ORDER_STATE_CANCEL             int32 = 4
	ORDER_STATE_CANCEL_REFUNDING   int32 = 5
	ORDER_STATE_CANCEL_REFUND_DONE int32 = 6
	ORDER_STATE_DEL_REFUNDING      int32 = 7
	ORDER_STATE_DEL_REFUND_DONE    int32 = 8
	ORDER_STATE_CALLBACK_HANDLING  int32 = 9 // callbackflow handling.
)

const (
	TAKE_ORDER_DEFAULT int32 = 0
	TAKE_ORDER_ADMIN   int32 = 1
)

type Order struct {
	OrderId          string `json:"orderId" db:"OrderId"`
	Title            string `json:"title" db:"Title"`
	ProductName      string `json:"productName" db:"ProductName"`
	OrderType        string `json:"orderType" db:"OrderType"`
	OrderTime        int64  `json:"orderTime" db:"OrderTime"`
	PayTime          int64  `json:"payTime" db:"PayTime"`
	State            int32  `json:"state" db:"State"`
	NotifyPay        int32  `json:"notifyPay" db:"NotifyPay"`
	Amount           string `json:"amount" db:"Amount"`
	OntId            string `json:"ontId" db:"OntId"`
	UserId           string `json:"userId" db:"UserId"`
	UserName         string `json:"userName" db:"UserName"`
	TxHash           string `json:"txHash" db:"TxHash"`
	Price            string `json:"price" db:"Price"`
	ApiId            uint32 `json:"apiId" db:"ApiId"`
	ToolBoxId        uint32 `json:"toolBoxId" db:"ToolBoxId"`
	ApiUrl           string `json:"apiUrl" db:"ApiUrl"`
	SpecificationsId uint32 `json:"specificationsId" db:"SpecificationsId"`
	OrderKind        uint32 `json:"orderKind" db:"OrderKind"`
	Request          string `json:"request" db:"Request"` // this fill
	Result           string `json:"result" db:"Result"`   // this fill
	ApiKey           string `json:"apiKey" db:"ApiKey"`
	Reason           string `json:"reason" db:"Reason"`
}

const (
	API_KEY_TYPE_COUNT    = 1
	API_KEY_TYPE_DURATION = 2
)

type APIKey struct {
	Id           uint32 `json:"id" db:"Id"`
	ApiKey       string `json:"apiKey" db:"ApiKey"`
	OrderId      string `json:"orderId" db:"OrderId"`
	ApiId        uint32 `json:"apiId" db:"ApiId"`
	RequestLimit uint64 `json:"requestLimit" db:"RequestLimit"`
	UsedNum      uint64 `json:"usedNum" db:"UsedNum"`
	OntId        string `json:"ontId" db:"OntId"`
	UserId       string `json:"userId" db:"UserId"`
	OutDate      int64  `json:"outDate" db:"OutDate"`
	ApiKeyType   int32  `json:"apiKeyType" db:"ApiKeyType"`
	CreateTime   int64  `json:"createTime" db:"CreateTime"`
	Layer2Time   int64  `json:"layer2Time"`
}

type QrCode struct {
	QrCodeId     string `json:"id" db:"QrCodeId"`
	Ver          string `json:"ver" db:"Ver"`
	OrderId      string `json:"orderId" db:"OrderId"`
	Requester    string `json:"requester" db:"Requester"`
	Signature    string `json:"signature" db:"Signature"`
	Signer       string `json:"signer" db:"Signer"`
	QrCodeData   string `json:"data" db:"QrCodeData"`
	Callback     string `json:"callback" db:"Callback"`
	Exp          int64  `json:"exp" db:"Exp"`
	Chain        string `json:"chain" db:"Chain"`
	QrCodeDesc   string `json:"desc" db:"QrCodeDesc"`
	ContractType string `json:"contractType" db:"ContractType"`
}

type QrCodeDesc struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Price  string `json:"price"`
}

type QrCodeCallBackData struct {
	QrCodeId     string     `json:"id" db:"QrCodeId"`
	Ver          string     `json:"ver" db:"Ver"`
	OrderId      string     `json:"orderId" db:"OrderId"`
	Requester    string     `json:"requester" db:"Requester"`
	Signature    string     `json:"signature" db:"Signature"`
	Signer       string     `json:"signer" db:"Signer"`
	QrCodeData   string     `json:"data" db:"QrCodeData"`
	Callback     string     `json:"callback" db:"Callback"`
	Exp          int64      `json:"exp" db:"Exp"`
	Chain        string     `json:"chain" db:"Chain"`
	QrCodeDesc   QrCodeDesc `json:"desc" db:"QrCodeDesc"`
	ContractType string     `json:"contractType" db:"ContractType"`
}

type OrderApiComment struct {
	Id          uint32 `json:"id" db:"Id"`
	OrderId     string `json:"orderId" db:"OrderId"`
	ApiId       uint32 `json:"apiId" db:"ApiId"`
	ToolBoxId   uint32 `json:"toolBoxId" db:"ToolBoxId"`
	StarNum     uint32 `json:"starNum" db:"StarNum"`
	Comments    string `json:"comments" db:"Comments"`
	UserName    string `json:"userName" db:"UserName"`
	OntId       string `json:"ontId" db:"OntId"`
	CommentTime int64  `json:"commentTime" db:"CommentTime"`
	State       byte   `json:"state" db:"State"`
}

type UserName struct {
	Id                          string         `json:"id" db:"id"`
	OntId                       sql.NullString `json:"ontId" db:"ont_id"`
	CentOntId                   sql.NullString `json:"centOntId" db:"cent_ont_id"`
	UserName                    sql.NullString `json:"userName" db:"user_name"`
	Email                       sql.NullString `json:"email" db:"email"`
	Password                    sql.NullString `json:"password" db:"password"`
	PrivateKey                  sql.NullString `json:"privateKey" db:"private_key"`
	PhoneNumber                 sql.NullString `json:"phoneNumber" db:"phone_number"`
	AvatarUrl                   sql.NullString `json:"avatarUrl" db:"avatar_url"`
	State                       uint32         `json:"state" db:"state"`
	RegisterAgreement           byte           `json:"registerAgreement" db:"register_agreement"`
	SettlementAgreement         byte           `json:"settlementAgreement" db:"settlement_agreement"`
	AccountType                 uint32         `json:"accountType" db:"account_type"`
	PersonalVerification        byte           `json:"personalVerification" db:"personal_verification"`
	PersonalVerificationTime    sql.NullTime   `json:"personalVerificationTime" db:"personal_verification_time"`
	EnterpriseVerification      byte           `json:"enterpriseVerification" db:"enterprise_verification"`
	EnterpriseVerification_time sql.NullTime   `json:"enterpriseVerificationTime" db:"enterprise_verification_time"`
	CreateTime                  sql.NullTime   `json:"createTime" db:"create_time"`
}

type EnterpriseInfo struct {
	Id                       uint32         `json:"id" db:"id"`
	UserId                   string         `json:"userId" db:"user_id"`
	EnterpriseLegalName      string         `json:"enterpriseLegalName" db:"enterprise_legal_name"`
	RegistrationNumber       string         `json:"registrationNumber" db:"registration_number"`
	LegalRepresentative      string         `json:"legalRepresentative" db:"legal_representative"`
	RegistrationCertificate  sql.NullString `json:"registrationCertificate" db:"registration_certificate"`
	ThirdPartyAuthentication byte           `json:"thirdPartyAuthentication" db:"third_party_authentication"`
	Reason                   string         `json:"reason" db:"reason"`
	State                    uint32         `json:"state" db:"state"`
	CreateTime               sql.NullTime   `json:"createTime" db:"create_time"`
}

type OntIdScore struct {
	OntId           string `json:"ontId" db:"OntId"`
	TotalCommentNum uint64 `json:"totalCommentNum" db:"TotalCommentNum"`
	TotalScore      uint64 `json:"totalScore" db:"TotalScore"`
}

type Notification struct {
	Id         string `json:"id" db:"id"`
	UserId     string `json:"userId" db:"user_id"`
	Type       int32  `json:"type" db:"type"`
	TitleEn    string `json:"titleEn" db:"title_en"`
	TitleZh    string `json:"titleZh" db:"title_zh"`
	ContentEn  string `json:"contentEn" db:"content_en"`
	ContentZh  string `json:"contentZh" db:"content_zh"`
	KeyWord    string `json:"keyWord" db:"keyword"`
	BusinessId string `json:"businessId" db:"business_id"`
	IsRead     int32  `json:"isRead" db:"is_read"`
	CreateTime int64  `json:"createTime" db:"create_time"`
}
