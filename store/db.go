package store

const (
	RequestStateNew        = "new"
	RequestStateInProgress = "in_progress"
	RequestStateDone       = "done"
	RequestStateError      = "error"
)

type DataStore interface {
	SetRequest(r *Request) (int, error)
	Delete(id int) error
	GetRequest(id int) (*Request, error)
	GetAllRequests() ([]*Request, error)
	SaveRequest(r *Request) error
}


type ClientRequest struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    interface{}         `json:"body"`
}

type Response struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	BodyLen     int64              `json:"body_len"`
}

type Request struct {
	ID            int            `json:"id"`
	ClientRequest *ClientRequest `json:"request"`
	Response      *Response      `json:"response"`
	Status        *ExecStatus    `json:"status"`
}

type ExecStatus struct {
	State string `json:"state"`
	Err   string `json:"error"`
}

func (r *Request) SetStatus(s string, err string) {
	r.Status.State = s
	r.Status.Err = err
}

func (r *Request) IsNew() bool {
	return r.Status.State == RequestStateNew
}

func (r *Request) IsInProgress() bool {
	return r.Status.State == RequestStateInProgress
}

func (r *Request) IsDone() bool {
	return r.Status.State == RequestStateDone
}

func (r *Request) IsError() bool {
	return r.Status.State == RequestStateError
}
