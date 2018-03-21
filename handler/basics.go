package handler

import (
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type HandlersWrapper struct {
	store.DbService
}

func (w *HandlersWrapper) RequestsForClient(ctx echo.Context) error {
	// метод выдает все сохрааненные просьбы
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().WriteHeader(http.StatusOK)
	allRequests := &struct {
		Data []*store.DataForDb `json:"Data"`
	}{w.GetAllData()}
	return ctx.JSON(http.StatusOK, allRequests)
}

func (w *HandlersWrapper) RequestForClientById(ctx echo.Context) error {
	//метод выдает информацию по просьбе по id
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	request, ok := w.Get(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, request)
}

func (w *HandlersWrapper) DeleteRequestForClient(ctx echo.Context) error {
	// функция для удаления просьбы
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	if _, ok := w.Get(tempid); !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	ok := w.Delete(tempid)
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
	resp, err := requester.RequestIssueExecutor(result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseToClient := &store.ResponseToClient{
		ResponseData: resp,
	}

	dataFoDb := &store.DataForDb{
		Request:      result,
		ResponseData: resp,
	}
	responseToClient.Id = w.Set(dataFoDb)
	return ctx.JSON(http.StatusOK, responseToClient)
}
