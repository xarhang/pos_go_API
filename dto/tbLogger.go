package dto

type TbLogger struct {
	Url         string `json:"url"`
	Method      string `json:"method"`
	Request     string `json:"request"`
	RequestBody string `json:"reuest_body"`
	Code        int    `json:"code"`
	Response    string `json:"response"`
	Accesstime  string `json:"accesstime"`
	Handletime  string `json:"handletime"`
}
