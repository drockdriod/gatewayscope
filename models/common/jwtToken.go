package common


// Binding from JSON
type JwtToken struct {
	Sub		 	string `json:"sub"`
	Aud     	string `json:"aud"`
	Iss 	 	string `json:"iss"`
	Exp 		string `json:"exp"`
}