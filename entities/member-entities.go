package entities

type Model_member struct {
	Member_id         string `json:"member_id"`
	Member_username   string `json:"member_username"`
	Member_timezone   string `json:"member_timezone"`
	Member_ipaddress  string `json:"member_ipaddress"`
	Member_lastlogin  string `json:"member_lastlogin"`
	Member_name       string `json:"member_name"`
	Member_phone      string `json:"member_phone"`
	Member_email      string `json:"member_email"`
	Member_status     string `json:"member_status"`
	Member_status_css string `json:"member_status_css"`
	Member_create     string `json:"member_create"`
	Member_update     string `json:"member_update"`
}

type Controller_membersave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Member_id       string `json:"member_id"`
	Member_username string `json:"member_username" validate:"required"`
	Member_password string `json:"member_password"`
	Member_name     string `json:"member_name" validate:"required"`
	Member_phone    string `json:"member_phone"`
	Member_email    string `json:"member_email"`
	Member_status   string `json:"member_status"`
}
