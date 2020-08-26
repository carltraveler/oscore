package tables

type Location struct {
	Id      uint32 `json:"id" db:"Id"`
	Country string `json:"country" db:"Country"`
	City    string `json:"city" db:"City"`
	Lat     string `json:"lat" db:"Lat"`
	Lng     string `json:"lng" db:"Lng"`
}

type Algorithm struct {
	Id           uint32 `json:"id" db:"Id"`
	AlgName      string `json:"algName" db:"AlgName"`
	Provider     string `json:"provider" db:"Provider"`
	Description  string `json:"description" db:"Description"`
	Price        string `json:"price" db:"Price"`
	Coin         string `json:"coin" db:"Coin"`
	ResourceId   string `json:"resourceId" db:"ResourceId"`
	TokenHash    string `json:"tokenHash" db:"TokenHash"`
	OwnerAddress string `json:"ownerAddress" db:"OwnerAddress"`
	State        byte   `json:"state" db:"State"`
	CreateTime   Time   `json:"createTime" db:"CreateTime"`
}

type Env struct {
	Id           uint32 `json:"id" db:"Id"`
	EnvName      string `json:"envName" db:"EnvName"`
	Provider     string `json:"provider" db:"Provider"`
	Description  string `json:"description" db:"Description"`
	Price        string `json:"price" db:"Price"`
	Coin         string `json:"coin" db:"Coin"`
	ServiceUrl   string `json:"ServiceUrl" db:"ServiceUrl"`
	ResourceId   string `json:"resourceId" db:"ResourceId"`
	TokenHash    string `json:"tokenHash" db:"TokenHash"`
	OwnerAddress string `json:"ownerAddress" db:"OwnerAddress"`
	State        byte   `json:"state" db:"State"`
	CreateTime   Time   `json:"createTime" db:"CreateTime"`
}

type ApiAlgorithm struct {
	Id          uint32 `json:"id" db:"Id"`
	ApiId       uint32 `json:"apiId" db:"ApiId"`
	AlgorithmId uint32 `json:"algorithmId" db:"AlgorithmId"`
	State       byte   `json:"state" db:"State"`
	CreateTime  Time   `json:"createTime" db:"CreateTime"`
}

type AlgorithmEnv struct {
	Id          uint32 `json:"id" db:"Id"`
	AlgorithmId uint32 `json:"algorithmId" db:"AlgorithmId"`
	EnvId       uint32 `json:"envId" db:"EnvId"`
	State       byte   `json:"state" db:"State"`
	CreateTime  Time   `json:"createTime" db:"CreateTime"`
}

type ToolBox struct {
	Id          uint32 `json:"id" db:"Id"`
	Title       string `json:"title" db:"Title"`
	ToolBoxDesc string `json:"toolBoxDesc" db:"ToolBoxDesc"`
	Icon        string `json:"icon" db:"Icon"`
	Coin        string `json:"coin" db:"Coin"`
	Price       string `json:"price" db:"Price"`
	ToolBoxType string `json:"type" db:"ToolBoxType"`
	State       byte   `json:"state" db:"State"`
	CreateTime  Time   `json:"createTime" db:"CreateTime"`
}

// no need sql tables. merge to request params.
type ApiHeadValues struct {
	HeaderKey   string `json:"headerKey" db:"HeaderKey"`
	HeaderValue string `json:"headerValue" db:"HeaderValue"`
}
