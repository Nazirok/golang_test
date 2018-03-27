package handler

import (
	"fmt"
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

const (
	newjb     = "new"
	perfoming = "is performing"
	perfomed  = "perfomed"
)

var queue = make(chan int)

type HandlersWrapper struct {
	db  store.DbService
	jdb store.JobDbService
}

func (w *HandlersWrapper) RequestsForClient(ctx echo.Context) error {
	// метод выдает все сохрааненные просьбы
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().WriteHeader(http.StatusOK)
	allRequests := &struct {
		Data []*store.DataForDb `json:"Request"`
	}{w.db.GetAllData()}
	return ctx.JSON(http.StatusOK, allRequests)
}

func (w *HandlersWrapper) RequestForClientById(ctx echo.Context) error {
	//метод выдает информацию по просьбе по id
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	request, ok := w.db.Get(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, request)
}

func (w *HandlersWrapper) DeleteRequestForClient(ctx echo.Context) error {
	// функция для удаления просьбы
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	if _, ok := w.db.Get(tempid); !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	ok := w.db.Delete(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

func (w *HandlersWrapper) RequestFromClientHandler(ctx echo.Context) error {
	result := &store.ClientBody{}
	if err := ctx.Bind(result); err != nil {
		return err
	}
	id := w.toQueue(result)
	res := struct {
		Checkid int `json:"Checkid"`
	}{id}
	return ctx.JSON(http.StatusOK, res)
}

func (w *HandlersWrapper) CheckResponse(ctx echo.Context) error {
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	job, ok := w.jdb.Get(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	switch job.State {
	case newjb, perfoming:
		res := struct {
			State string `json:"State"`
		}{perfoming}
		return ctx.JSON(http.StatusAccepted, res)
	case perfomed:
		defer w.jdb.Delete(tempid)
		if job.Err != nil {
			res := struct {
				Error string
			}{fmt.Sprintf("%s", job.Err)}
			return ctx.JSON(http.StatusInternalServerError, res)
		}
		return ctx.JSON(http.StatusOK, job.ToClient)
	}
	return nil
}

func (w *HandlersWrapper) toQueue(r *store.ClientBody) int {
	job := &store.Job{newjb, r, nil, nil}
	id := w.jdb.Set(job)
	go func() { queue <- id }()
	return id
}

func (w *HandlersWrapper) JobExecutor() {
	for {
		i := <-queue
		data := w.jdb.ChangeState(i, perfoming, nil, nil)
		go func(i int) {
			resp, err := requester.RequestIssueExecutor(data.Request)
			if err != nil {
				w.jdb.ChangeState(i, perfomed, nil, err)
				return
			}
			responseToClient := &store.ResponseToClient{
				ResponseData: resp,
			}
			dataFoDb := &store.DataForDb{
				Request:      data.Request,
				ResponseData: resp,
			}
			responseToClient.Id = w.db.Set(dataFoDb)
			w.jdb.ChangeState(i, perfomed, responseToClient, nil)
		}(i)
	}
}

func New(db store.DbService, jdb store.JobDbService) *HandlersWrapper {
	return &HandlersWrapper{db, jdb}
}
