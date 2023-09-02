package helpers

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Record  interface{} `json:"record"`
	Time    string      `json:"time"`
}
type Responsepaging struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	Record      interface{} `json:"record"`
	Perpage     int         `json:"perpage"`
	Totalrecord int         `json:"totalrecord"`
	Time        string      `json:"time"`
}
type Responsemember struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Record   interface{} `json:"record"`
	Listbank interface{} `json:"listbank"`
	Time     string      `json:"time"`
}
type Responsemembercredit struct {
	Status        int     `json:"status"`
	Message       string  `json:"message"`
	Member_id     string  `json:"member_id"`
	Member_name   string  `json:"member_name"`
	Member_credit float64 `json:"member_credit"`
	Time          string  `json:"time"`
}
type Responsemaster struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Record   interface{} `json:"record"`
	Listcurr interface{} `json:"listcurr"`
	Listbank interface{} `json:"listbank"`
	Time     string      `json:"time"`
}
type Responseprovider struct {
	Status       int         `json:"status"`
	Message      string      `json:"message"`
	Record       interface{} `json:"record"`
	Listprovider interface{} `json:"listprovider"`
	Time         string      `json:"time"`
}
type ResponseAdmin struct {
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Record   interface{} `json:"record"`
	Listrule interface{} `json:"listruleadmin"`
	Time     string      `json:"time"`
}
type ResponseEmployee struct {
	Status          int         `json:"status"`
	Message         string      `json:"message"`
	Record          interface{} `json:"record"`
	Listdepartement interface{} `json:"listdepartement"`
	Time            string      `json:"time"`
}
type ErrorResponse struct {
	Field string
	Tag   string
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
