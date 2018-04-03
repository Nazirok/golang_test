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
	GetRequest(id int) (*Request, error)
	GetAllRequests() ([]*Request, error)
	ExecRequest(id int) (*ClientRequest, error)
	SetResponse(id int, response *Response, err error) error
}