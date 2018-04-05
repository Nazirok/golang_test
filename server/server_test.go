package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang_test/handler"
	"github.com/golang_test/store"
	"github.com/golang_test/worker"
	"github.com/golang_test/сonstants"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
	id, _ := mapDb.SetRequest(&store.ClientRequest{})
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
	go wr.RequestExecuteLoop()
	s.InitHandlers(w)
	os.Exit(m.Run())
}

func Test_ExecNewRequest(t *testing.T) {
	clientRequest := &store.ClientRequest{
		Method: "GET", URL: "Some url", Body: "body", Headers: map[string][]string{"header": []string{"a", "bv"}},
	}
	temp, err := json.Marshal(clientRequest)
	if err != nil {
		t.Errorf("error during marshal request")
	}
	req := httptest.NewRequest("POST", "/requests", bytes.NewReader(temp))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Bad status code")
	body, err := ioutil.ReadAll(rec.Body)
	r := struct {
		RequestId int `json:"requestId"`
	}{}
	if err = json.Unmarshal(body, &r); err != nil {
		t.Error("error in demarshaling", err)
	}
	id := r.RequestId
	time.Sleep(time.Millisecond * 50)
	res, _ := mapDb.GetRequest(id)
	assert.Equal(t, *clientRequest, *res.ClientRequest)
	assert.Equal(t, сonstants.RequestStateDone, res.Status.State)
}

func Test_ExecNewLongRequest(t *testing.T) {
	requester.SetDelay(1 * time.Second)
	clientRequest := &store.ClientRequest{
		Method: "GET", URL: "Some url", Body: "body", Headers: map[string][]string{"header": []string{"a", "bv"}},
	}
	temp, err := json.Marshal(clientRequest)
	if err != nil {
		t.Errorf("error during marshal request")
	}
	req := httptest.NewRequest("POST", "/requests", bytes.NewReader(temp))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Bad status code")
	body, err := ioutil.ReadAll(rec.Body)
	r := struct {
		RequestId int `json:"requestId"`
	}{}
	if err = json.Unmarshal(body, &r); err != nil {
		t.Error("error in demarshaling", err)
	}
	id := r.RequestId
	res, _ := mapDb.GetRequest(id)
	assert.Equal(t, *clientRequest, *res.ClientRequest)
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, сonstants.RequestStateInProgress, res.Status.State)
	time.Sleep(time.Second * 2)
	res, _ = mapDb.GetRequest(id)
	assert.Equal(t, *clientRequest, *res.ClientRequest)
	assert.Equal(t, сonstants.RequestStateDone, res.Status.State)
}

func Test_GetSuccessResponse(t *testing.T) {
	id := createRequest()
	resp := &store.Response{StatusCode: 200, Body: "someBody"}
	mapDb.SetResponse(id, resp, nil)
	url := fmt.Sprintf("/requests/%d", id)
	req := httptest.NewRequest(echo.GET, url, nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Bad status code")
	body, err := ioutil.ReadAll(rec.Body)
	r := &store.Request{}
	if err = json.Unmarshal(body, &r); err != nil {
		t.Error("error in demarshaling", err)
	}
	assert.Equal(t, resp, r.Response, "Bad response")
	assert.Equal(t, сonstants.RequestStateDone, r.Status.State, "Bad state of response")
}

func Test_GetNotExistedResponse(t *testing.T) {
	url := fmt.Sprintf("/requests/%d", 34568365)
	req := httptest.NewRequest(echo.GET, url, nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, "Bad status code")
}

func Test_GetErrorResponse(t *testing.T) {
	id := createRequest()
	mapDb.SetResponse(id, nil, errors.New("error.during.request"))
	url := fmt.Sprintf("/requests/%d", id)
	req := httptest.NewRequest(echo.GET, url, nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Bad status code")
	body, err := ioutil.ReadAll(rec.Body)
	r := &store.Request{}
	if err = json.Unmarshal(body, &r); err != nil {
		t.Error("error in demarshaling", err)
	}
	assert.Equal(t, сonstants.RequestStateError, r.Status.State, "Bad state of response")
	assert.Equal(t, "error.during.request", r.Status.Err)
}

func Test_GetNotReadyResponse(t *testing.T) {
	doRequest := func(id int) {
		url := fmt.Sprintf("/requests/%d", id)
		req := httptest.NewRequest(echo.GET, url, nil)
		rec := httptest.NewRecorder()
		s.e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code, "Bad status code")
		body, err := ioutil.ReadAll(rec.Body)
		r := struct {
			State string `json:"state"`
		}{}
		if err = json.Unmarshal(body, &r); err != nil {
			t.Error("error in demarshaling", err)
		}
		assert.Equal(t, сonstants.RequestStateInProgress, r.State, "Bad state in response")
	}
	id := createRequest()
	doRequest(id)
	mapDb.ExecRequest(id)
	doRequest(id)
}

func Test_GetAllRequests(t *testing.T) {
	for i := 1; i < 10; i++ {
		createRequest()
	}
	req := httptest.NewRequest("GET", "/requests", nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Bad status code")
	body, err := ioutil.ReadAll(rec.Body)
	r := &store.Request{}
	if err = json.Unmarshal(body, &r); err != nil {
		t.Error("error in demarshaling", err)
	}
	requests := &struct {
		Data []*store.Request `json:"clientRequests"`
	}{}
	if err = json.Unmarshal(body, &requests); err != nil {
		t.Error("error in demarshaling", err)
	}
	var g = assert.Comparison(func() (success bool) {
		return len(requests.Data) >= 10
	})
	assert.Condition(t, g, "Bad len of response")
}

func Test_DeleteRequest(t *testing.T) {
	id := createRequest()
	r, err := mapDb.GetRequest(id)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, r)
	url := fmt.Sprintf("/requests/%d", id)
	req := httptest.NewRequest("DELETE", url, nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Bad status code")
	r, err = mapDb.GetRequest(id)
	assert.NotEqual(t, nil, err)
}

func Test_DeleteNotExistedRequest(t *testing.T) {
	url := fmt.Sprintf("/requests/%d", 8908089)
	req := httptest.NewRequest("DELETE", url, nil)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, "Bad status code")
}
