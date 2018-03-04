package store

type ClientBody struct {
	Method      string
	Url         string
	ContentType string `json:"content-type"`
	Body        interface{}
}

type DbService interface {
	Set(key int, value ClientBody)
	Delete(key int) bool
	Get(key int) (ClientBody, bool)
	GetAllDataJson() ([]byte, error)
}
