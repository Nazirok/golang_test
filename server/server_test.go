package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang_test/handler"
	"github.com/golang_test/store"
	"github.com/golang_test/worker"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"github.com/stretchr/testify/require"
	"github.com/labstack/echo"
	"fmt"
)

var s = New()
var mapDb = store.NewMapDataStore()
var requester = &testRequester{0}
var wr = worker.NewRequestsExecutorByChan(mapDb, requester)
var w = handler.New(mapDb, wr)

type testRequester struct {
	delay time.Duration
}

func createRequest() int {
	id, _ := mapDb.SetRequest(&store.Request{})
	return id
}

func (r *testRequester) SetDelay(t time.Duration) {
	r.delay = t
}

func (r *testRequester) Do(result *store.ClientRequest) (resp *store.Response, err error) {
	if r.delay != 0 {
		time.Sleep(r.delay)
	}
	return nil, nil
}

func TestMain(m *testing.M) {
	go wr.RequestsExecuteLoop()
	s.InitHandlers(w)
	os.Exit(m.Run())
}

func doJSONRequest(t *testing.T, method, url string, requestBody interface{}, expectedHttpStatus int, response interface{}) {
	var data []byte
	if requestBody != nil {
		var err error
		data, err = json.Marshal(requestBody)
		require.NoError(t, err)
	}
	req := httptest.NewRequest(method, url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, expectedHttpStatus, rec.Code, "Bad status code")
	body, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)
	e := json.Unmarshal(body, &response)
	require.NoError(t, e)
}


func Test_ExecNewRequest(t *testing.T) {
	clientRequest := &store.ClientRequest{
		Method: "GET", URL: "Some url", Body: "body", Headers: map[string][]string{"header": []string{"a", "bv"}},
	}
	r := struct {
		RequestId int `json:"requestId"`
	}{}
	doJSONRequest(t, echo.POST, "/requests", clientRequest, http.StatusOK, &r)
	time.Sleep(time.Millisecond * 50)
	res, err := mapDb.GetRequest(r.RequestId)
	require.NoError(t, err)
	assert.Equal(t, *clientRequest, *res.ClientRequest)
	assert.Equal(t, store.RequestStateDone, res.Status.State)
}

func Test_ExecNewLongRequest(t *testing.T) {
	requester.SetDelay(1 * time.Second)
	clientRequest := &store.ClientRequest{
		Method: "GET", URL: "Some url", Body: "body", Headers: map[string][]string{"header": []string{"a", "bv"}},
	}
	r := struct {
		RequestId int `json:"requestId"`
	}{}
	doJSONRequest(t, echo.POST, "/requests", clientRequest, http.StatusOK, &r)
	res, err := mapDb.GetRequest(r.RequestId)
	require.NoError(t, err)
	assert.Equal(t, *clientRequest, *res.ClientRequest)
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, store.RequestStateInProgress, res.Status.State)
	time.Sleep(time.Second * 2)
	res, _ = mapDb.GetRequest(r.RequestId)
	assert.Equal(t, *clientRequest, *res.ClientRequest)
	assert.Equal(t, store.RequestStateDone, res.Status.State)
}

func Test_GetSuccessResponse(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	require.NoError(t, err)
	resp := &store.Response{StatusCode: 200, BodyLen: 4}
	req.Response = resp
	req.ClientRequest =  &store.ClientRequest{}
	req.Status = &store.ExecStatus{State: store.RequestStateDone, Err: ""}
	mapDb.SaveRequest(req)
	response := store.Request{}
	url := fmt.Sprintf("/requests/%d", id)
	doJSONRequest(t, echo.GET, url, nil, http.StatusOK, &response)
	assert.Equal(t, resp, response.Response, "Bad response")
	assert.Equal(t, store.RequestStateDone, response.Status.State, "Bad state of response")
}

func Test_GetNotExistedResponse(t *testing.T) {
	url := fmt.Sprintf("/requests/%d", 34568365)
	doJSONRequest(t, echo.GET, url, nil, http.StatusNotFound, nil)
}

func Test_GetErrorResponse(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	require.NoError(t, err)
	resp := &store.Response{StatusCode: 200, BodyLen: 4}
	req.Response = resp
	req.ClientRequest =  &store.ClientRequest{}
	req.Status = &store.ExecStatus{State: store.RequestStateError, Err: errors.New("error.during.request").Error()}
	mapDb.SaveRequest(req)
	response := store.Request{}
	url := fmt.Sprintf("/requests/%d", id)
	doJSONRequest(t, echo.GET, url, nil, http.StatusOK, &response)
	assert.Equal(t, store.RequestStateError, response.Status.State, "Bad state of response")
	assert.Equal(t, "error.during.request", response.Status.Err)
}

func Test_GetNotReadyResponse(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	req.Status = &store.ExecStatus{State: store.RequestStateNew, Err: ""}
	require.NoError(t, err)
	mapDb.SaveRequest(req)
	r := struct {
		State string `json:"state"`
	}{}
	doRequestHelper := func(){
		url := fmt.Sprintf("/requests/%d", id)
		doJSONRequest(t, echo.GET, url, nil, http.StatusAccepted, &r)
		assert.Equal(t, store.RequestStateInProgress, r.State, "Bad state in response")
	}
	doRequestHelper()
	req.Status = &store.ExecStatus{State: store.RequestStateInProgress, Err: ""}
	doRequestHelper()
}

func Test_GetAllRequests(t *testing.T) {
	for i := 1; i < 10; i++ {
		id := createRequest()
		req, err := mapDb.GetRequest(id)
		require.NoError(t, err)
		req.Status = &store.ExecStatus{State: store.RequestStateNew, Err: ""}
		mapDb.SaveRequest(req)
	}
	requests := struct {
		Data []*store.Request `json:"clientRequests"`
	}{}
	doJSONRequest(t, echo.GET, "/requests", nil, http.StatusOK, &requests)
	var g = assert.Comparison(func() (success bool) {
		return len(requests.Data) >= 10
	})
	assert.Condition(t, g, "Bad len of response")
}

func Test_DeleteRequest(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	require.NoError(t, err)
	req.Status = &store.ExecStatus{State: store.RequestStateNew, Err: ""}
	mapDb.SaveRequest(req)
	response := store.Request{}
	url := fmt.Sprintf("/requests/%d", id)
	doJSONRequest(t, echo.DELETE, url, nil, http.StatusOK, &response)
	assert.Equal(t, *req, response)
	req, _ = mapDb.GetRequest(id)
	if req != nil {
		t.Error("req must be nil")
	}
}

func Test_DeleteNotExistedRequest(t *testing.T) {
	url := fmt.Sprintf("/requests/%d", 34568365)
	doJSONRequest(t, echo.DELETE, url, nil, http.StatusNotFound, nil)
}
