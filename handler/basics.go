package handler

import (
	"github.com/golang_test/store"
	"github.com/golang_test/worker"
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
	req, err := w.GetRequest(tempid)
	if err != nil {
		e := errorResponse{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, e)
	}
	if req == nil {
		e := errorResponse{store.RequestNotFound}
		return ctx.JSON(http.StatusNotFound, e)
	}

	if req.IsNew() || req.IsInProgress() {
		s := stateResponse{store.RequestStateInProgress}
		return ctx.JSON(http.StatusAccepted, s)
	} else {
		return ctx.JSON(http.StatusOK, req)
	}
	return nil
}

func (w *HandlersWrapper) DeleteRequestForClient(ctx echo.Context) error {
	// функция для удаления просьбы
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	req, _ := w.GetRequest(tempid)
	if req != nil {
		if req.IsInProgress() {
			s := stateResponse{store.RequestStateInProgress}
			return ctx.JSON(http.StatusAccepted, s)
		}
	}
	req, err := w.Delete(tempid)
	if err != nil {
		e := errorResponse{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, e)
	}
	if req == nil {
		e := errorResponse{store.RequestNotFound}
		return ctx.JSON(http.StatusNotFound, e)
	}
	return ctx.JSON(http.StatusOK, req)
}

func (w *HandlersWrapper) RequestFromClientHandler(ctx echo.Context) error {
	req := &store.Request{
		ClientRequest: &store.ClientRequest{},
		Status:        &store.ExecStatus{State: store.RequestStateNew, Err: ""},
	}
	if err := ctx.Bind(req.ClientRequest); err != nil {
		e := errorResponse{err.Error()}
		return ctx.JSON(http.StatusBadRequest, e)
	}
	id, err := w.SetRequest(req)
	if err != nil {
		e := errorResponse{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, e)
	}
	w.AddRequest(id)
	r := requestIdResponse{id}
	return ctx.JSON(http.StatusOK, r)
}
