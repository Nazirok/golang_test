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

type DbService interface {
	Set(value *DataForDb) int
	Delete(key int) bool
	Get(key int) (*DataForDb, bool)
	GetAllData() chan *DataForDb // канал тут явно не в тему, почему не слайс просто?
}
