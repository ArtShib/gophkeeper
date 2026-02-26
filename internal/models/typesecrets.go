package models

type Credentials struct {
	URL      string `json:"url"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Notes    string `json:"Notes"`
}

type Card struct {
	Number   string `json:"Number"`
	Holder   string `json:"Holder"`
	ExpiryMM string `json:"ExpiryMM"`
	ExpiryYY string `json:"ExpiryYY"`
	CVV      string `json:"CVV"`
	Bank     string `json:"Bank"`
}

type Note struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}
