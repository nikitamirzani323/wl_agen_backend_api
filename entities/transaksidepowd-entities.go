package entities

type Model_transdpwd struct {
	Transdpwd_id         string  `json:"transdpwd_id"`
	Transdpwd_date       string  `json:"transdpwd_date"`
	Transdpwd_idcurr     string  `json:"transdpwd_idcurr"`
	Transdpwd_tipedoc    string  `json:"transdpwd_tipedoc"`
	Transdpwd_tipeakun   string  `json:"transdpwd_tipeakun"`
	Transdpwd_idmember   string  `json:"transdpwd_idmember"`
	Transdpwd_notebank   string  `json:"transdpwd_notebank"`
	Transdpwd_amount     float32 `json:"transdpwd_amount"`
	Transdpwd_before     float32 `json:"transdpwd_before"`
	Transdpwd_after      float32 `json:"transdpwd_after"`
	Transdpwd_ipaddress  string  `json:"transdpwd_ipaddress"`
	Transdpwd_note       string  `json:"transdpwd_note"`
	Transdpwd_status     string  `json:"transdpwd_status"`
	Transdpwd_status_css string  `json:"transdpwd_status_css"`
	Transdpwd_create     string  `json:"transdpwd_create"`
	Transdpwd_update     string  `json:"transdpwd_update"`
}
type Controller_transdpwdsave struct {
	Page                string  `json:"page" validate:"required"`
	Sdata               string  `json:"sdata" validate:"required"`
	Transdpwd_id        string  `json:"transdpwd_id"`
	Transdpwd_tipedoc   string  `json:"transdpwd_tipedoc" validate:"required"`
	Transdpwd_idmember  string  `json:"transdpwd_idmember" validate:"required"`
	Transdpwd_bank_in   int     `json:"transdpwd_bank_in" validate:"required"`
	Transdpwd_bank_out  int     `json:"transdpwd_bank_out" validate:"required"`
	Transdpwd_amount    float32 `json:"transdpwd_amount" validate:"required"`
	Transdpwd_ipaddress string  `json:"transdpwd_ipaddress"`
	Transdpwd_note      string  `json:"transdpwd_note"`
	Transdpwd_status    string  `json:"transdpwd_status" validate:"required"`
}
