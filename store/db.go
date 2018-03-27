package store

type ClientBody struct {
	Method  string              `json:"Method"`
	Url     string              `json:"Url"`
	Headers map[string][]string `json:"Headers"`
	Body    interface{}         `json:"Body"`
}

type ResponseData struct {
	Status  int                 `json:"Status-code"`
	Headers map[string][]string `json:"Headers"`
	Length  int64               `json:"Content-Length"`
}

type ResponseToClient struct {
	Id           int           `json:"id"`
	ResponseData *ResponseData `json:"ResponseData"`
}

type DataForDb struct {
	Id           int           `json:"Id"`
	Request      *ClientBody   `json:"Request"`
	ResponseData *ResponseData `json:"ResponseData"`
}

type Job struct {
	State    string
	Request  *ClientBody
	ToClient *ResponseToClient
	Err      error
}

type DbService interface {
	Set(value *DataForDb) int
	Delete(key int) bool
	Get(key int) (*DataForDb, bool)
	GetAllData() []*DataForDb
}

type JobDbService interface {
	Set(value *Job) int
	Delete(key int) bool
	Get(key int) (*Job, bool)
	ChangeState(key int, s string, t *ResponseToClient, e error) *Job
}
