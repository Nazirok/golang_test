package server

import (
	"errors"
	"github.com/golang_test/handler"
	"github.com/golang_test/store"
	"github.com/golang_test/worker"
	"net/http"
	"os"
	"testing"
	"time"
	"github.com/stretchr/testify/require"
	"github.com/gavv/httpexpect"
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

func Test_ExecNewRequest(t *testing.T) {
	clientRequest := &store.ClientRequest{
		Method: "GET", URL: "Some url", Body: "body", Headers: map[string][]string{"header": []string{"a", "bv"}},
	}
	obj := tester(t).POST("/requests").
		WithJSON(clientRequest).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("requestId").NotEqual(nil)
	raw := obj.Value("requestId").Raw()
	time.Sleep(time.Millisecond * 50)
	// Request must be transited to state 'done' by worker after ~ 50 milliseconds
	res, err := mapDb.GetRequest(int(raw.(float64)))
	require.NoError(t, err)
	require.Equal(t, *clientRequest, *res.ClientRequest)
	require.Equal(t, store.RequestStateDone, res.Status.State)
}

func Test_ExecNewLongRequest(t *testing.T) {
	requester.SetDelay(1 * time.Second)
	clientRequest := &store.ClientRequest{
		Method: "GET", URL: "Some url", Body: "body", Headers: map[string][]string{"header": []string{"a", "bv"}},
	}
	obj := tester(t).POST("/requests").
		WithJSON(clientRequest).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("requestId").NotEqual(nil)
	raw := obj.Value("requestId").Raw()
	id := int(raw.(float64))
	time.Sleep(time.Millisecond * 50)
	res, err := mapDb.GetRequest(id)
	require.NoError(t, err)
	require.Equal(t, *clientRequest, *res.ClientRequest)
	time.Sleep(time.Millisecond * 50)
	// Request must be transited to state 'in_progress' by worker after ~ 50 milliseconds
	require.Equal(t, store.RequestStateInProgress, res.Status.State)
	time.Sleep(time.Second * 2)
	// Request must be transited to state 'done' by worker after 2 seconds
	res, _ = mapDb.GetRequest(id)
	require.Equal(t, *clientRequest, *res.ClientRequest)
	require.Equal(t, store.RequestStateDone, res.Status.State)
}


func tester(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(s.e),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
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
	obj := tester(t).GET("/requests/{ID}", req.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("response").Equal(resp)
	obj.Value("status").Object().ValueEqual("state", store.RequestStateDone)
}

func Test_GetNotExistedResponse(t *testing.T) {
	obj := tester(t).GET("/requests/{ID}", 12127).
		Expect().
		Status(http.StatusNotFound).
		JSON().Object()
	obj.Value("error").Equal(store.RequestNotFound)
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
	obj := tester(t).GET("/requests/{ID}", req.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("status").Object().ValueEqual("state", store.RequestStateError)
	obj.Value("status").Object().ValueEqual("error", "error.during.request")
	obj.Value("response").Equal(resp)
}

func Test_GetNotReadyResponse(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	req.Status = &store.ExecStatus{State: store.RequestStateNew, Err: ""}
	require.NoError(t, err)
	mapDb.SaveRequest(req)
	doRequestHelper := func() {
		obj := tester(t).GET("/requests/{ID}", req.ID).
			Expect().
			Status(http.StatusAccepted).
			JSON().Object()
		obj.ValueEqual("state", store.RequestStateInProgress)
	}
	// If state is "new" must return message: "in_progress"
	doRequestHelper()
	req.Status = &store.ExecStatus{State: store.RequestStateInProgress, Err: ""}
	// If state is "in_progress" must return message: "in_progress"
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
	obj := tester(t).GET("/requests").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("clientRequests").Array().Length().Ge(10)
}

func Test_DeleteRequest(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	require.NoError(t, err)
	req.Status = &store.ExecStatus{State: store.RequestStateNew, Err: ""}
	mapDb.SaveRequest(req)
	requestHelper := func(status int) *httpexpect.Value {
		obj := tester(t).DELETE("/requests/{ID}", req.ID).
			Expect().
			Status(status).
			JSON()
		return obj
	}
	requestHelper(http.StatusOK).Equal(req)
	requestHelper(http.StatusNotFound).Object().Value("error").Equal(store.RequestNotFound)
}

func Test_DeleteNotExistedRequest(t *testing.T) {
	obj := tester(t).DELETE("/requests/{ID}", 12127).
		Expect().
		Status(http.StatusNotFound).
		JSON().Object()
	obj.Value("error").Equal(store.RequestNotFound)
}

func Test_DeleteInProgressRequest(t *testing.T) {
	id := createRequest()
	req, err := mapDb.GetRequest(id)
	require.NoError(t, err)
	req.Status = &store.ExecStatus{State: store.RequestStateInProgress, Err: ""}
	mapDb.SaveRequest(req)
	obj := tester(t).DELETE("/requests/{ID}", req.ID).
		Expect().
		Status(http.StatusAccepted).
		JSON().Object()
	obj.ValueEqual("state", store.RequestStateInProgress)
}