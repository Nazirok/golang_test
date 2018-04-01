package store

type ClientRequest struct {
	Method  string              `json:"method"`
	Url     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    interface{}         `json:"body"`
}

type Response struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	BodyLength int64               `json:"body_length"`
}

type Request struct {
	ID            int            `json:"id"`
	ClientRequest *ClientRequest `json:"request"`
	Response      *Response      `json:"response"`
	Status        *ExecStatus    `json:"status"`
}

type ExecStatus struct {
	State    string `json:"status"`
	Err      string `json:"error"`
}

type ResponseToClient struct {
	ID           int       `json:"id"`
	ResponseData *Response `json:"response"`
}

type DataStore interface {
	SetRequest(r *ClientRequest) (int, error)
	Delete(id int) bool
	GetRequest(id int) (*Request, bool)
	GetAllRequests() ([]*Request, error)
	ExecRequest(id int) error
}

type JobDbService interface {
	Set(value *ExecStatus) int
	Delete(key int) bool
	Get(key int) (*ExecStatus, bool)
	ChangeState(key int, s string, t *ResponseToClient, e error) *ExecStatus
}
