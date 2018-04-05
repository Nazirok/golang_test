package handler

import (
	"github.com/golang_test/store"
	"github.com/golang_test/worker"
	"github.com/golang_test/сonstants"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type HandlersWrapper struct {
	store.DataStore
	worker.RequestsExecutor
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

func New(db store.DataStore, wr worker.RequestsExecutor) *HandlersWrapper {
	return &HandlersWrapper{db, wr}
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
	w.AddRequest(id)
	r := requestIdResponse{id}
	return ctx.JSON(http.StatusOK, r)
}
