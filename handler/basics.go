package handler

import (
	"fmt"
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

var queue = make(chan int)
var quit = make(chan struct{})

type HandlersWrapper struct {
	store.DataStore
}

func (w *HandlersWrapper) RequestsForClient(ctx echo.Context) error {
	// метод выдает все сохрааненные просьбы
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().WriteHeader(http.StatusOK)
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
	request, ok := w.GetRequest(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, request)
}

func (w *HandlersWrapper) DeleteRequestForClient(ctx echo.Context) error {
	// функция для удаления просьбы
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	if ok := w.Delete(tempid); !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if ok := w.Delete(tempid); !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
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
	res := struct {
		CheckID int `json:"checkid"`
	}{id}
	return ctx.JSON(http.StatusOK, res)
}

func (w *HandlersWrapper) CheckResponse(ctx echo.Context) error {
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	request, ok := w.GetRequest(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	switch request.Status.State {

	case "new", "is performing":
		res := struct {
			State string `json:"state"`
		}{"is performing"}
		return ctx.JSON(http.StatusAccepted, res)

	case "perfomed":
		return ctx.JSON(http.StatusOK, request.Response)

	case "error":
			res := struct {
				Error string
			}{fmt.Sprintf("%s", request.Status.Err)}
			return ctx.JSON(http.StatusInternalServerError, res)
	}
	return nil
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
					resp, err := requester.RequestIssueExecutor(clientRequest)
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

func New(db store.DataStore) *HandlersWrapper {
	return &HandlersWrapper{db}
}
