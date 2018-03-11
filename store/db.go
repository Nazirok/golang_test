package store

type ClientBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

type DbService interface {
	Set(value *ClientBody) int
	Delete(key int) bool
	Get(key int) (*ClientBody, bool)
	GetAllData() (map[int]*ClientBody)
}

