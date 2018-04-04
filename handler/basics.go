package handler

import (
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"github.com/golang_test/сonstants"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

var queue = make(chan int)
var quit = make(chan struct{})

type HandlersWrapper struct {
	store.DataStore
	r requester.Requester
}

type errorResponse struct {
	Error string `json:"error"`
}

type requestIdResponse struct {
	RequestId int `json:"requestId"`
}

type stateResponse struct {
	State string `json:"state"`
}

func (w *HandlersWrapper) RequestsForClient(ctx echo.Context) error {
	// метод выдает все сохрааненные просьбы
	requests, err := w.GetAllRequests()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	requestsJson := &struct {
		Data []*store.Request `json:"clientRequests"`
	}{requests}
	return ctx.JSON(http.StatusOK, requestsJson)
}

func (w *HandlersWrapper) RequestForClientById(ctx echo.Context) error {
	//метод выдает информацию по просьбе по id
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	request, err := w.GetRequest(tempid)
	if err != nil {
		e := errorResponse{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, e)
	}
	if request == nil {
		e := errorResponse{сonstants.RequestNotFound}
		return ctx.JSON(http.StatusNotFound, e)
	}

	switch request.Status.State {
	case сonstants.RequestStateNew, сonstants.RequestStateInProgress:
		s := stateResponse{сonstants.RequestStateInProgress}
		return ctx.JSON(http.StatusAccepted, s)

	case сonstants.RequestStateDone, сonstants.RequestStateError:
		return ctx.JSON(http.StatusOK, request)
	}
	return nil
}

func (w *HandlersWrapper) DeleteRequestForClient(ctx echo.Context) error {
	// функция для удаления просьбы
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	if err := w.Delete(tempid); err != nil {
		e := errorResponse{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, e)
	}
	return ctx.JSON(http.StatusOK, "OK")
}

func (w *HandlersWrapper) RequestFromClientHandler(ctx echo.Context) error {
	result := &store.ClientRequest{}
	if err := ctx.Bind(result); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	id, _ := w.SetRequest(result)
	w.toQueue(id)
	r := requestIdResponse{id}
	return ctx.JSON(http.StatusOK, r)
}

func (w *HandlersWrapper) toQueue(id int) {
	go func() { queue <- id }()
}

func (w *HandlersWrapper) JobExecutor() {
	for {
		select {
		case id := <-queue:
			clientRequest, err := w.ExecRequest(id)
			if err != nil {
				continue
			}
			go func(id int) {
				resp, err := w.r.Do(clientRequest)
				if err != nil {
					w.SetResponse(id, resp, err)
					return
				}
				w.SetResponse(id, resp, nil)
			}(id)
		case <-quit:
			return
		}
	}
}

func (w *HandlersWrapper) StopJobExecutor() {
	quit <- struct{}{}
}

func New(db store.DataStore, r requester.Requester) *HandlersWrapper {
	return &HandlersWrapper{db, r}
}
