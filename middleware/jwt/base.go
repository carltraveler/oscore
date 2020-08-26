package jwt

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Payload struct {
	Aud     string  `json:"aud"`
	Iss     string  `json:"iss"`
	Exp     int     `json:"exp"`
	Iat     int     `json:"iat"`
	Jti     string  `json:"jti"`
	Content Content `json:"content"`
}

type Content struct {
	Type   string `json:"type"`
	OntId  string `json:"ontId"`
	UserId string `json:"userId"`
	Role   string `json:"role"`
}
