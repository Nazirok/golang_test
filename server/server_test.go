package server

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang_test/handler"
	"github.com/golang_test/store"
	"github.com/golang_test/сonstants"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"errors"
	//"time"
)


var s = New()
var mapDb = store.NewMapDataStore()
var w = handler.New(mapDb, &testRequester{})
var wg sync.WaitGroup

type CheckBody struct {
	Checkid int
}

type testRequester struct {
}


func createRequest() int {
	//request := &store.Request{}
	id, _ := mapDb.SetRequest(&store.ClientRequest{})
	return id
}

func (r *testRequester) Do(result *store.ClientRequest) (resp *store.Response, err error) {
	return nil, nil
}


func TestMain(m *testing.M) {
	go w.JobExecutor()
	s.InitHandlers(w)
	os.Exit(m.Run())
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

func Test_GetErrorRequest(t *testing.T) {
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
	assert.Equal(t,сonstants.RequestStateError, r.Status.State, "Bad state of response")
	assert.Equal(t, "error.during.request", r.Status.Err)
}




