package handler

import (
	"github.com/golang_test/requester"
	"github.com/golang_test/store"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"encoding/json"
)

type HandlesrWrapper struct {
	store.DbService
	//requester.Requester
}

func (wrapper *HandlesrWrapper) RequestsForClient(ctx echo.Context) error {
	// метод выдает все сохрааненные просьбы
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(wrapper.GetAllData())
}

func (wrapper *HandlesrWrapper) RequestForClientById(ctx echo.Context) error {
	//метод выдает информацию по просьбе по id
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	request, ok := wrapper.Get(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, request)
}

func (wrapper *HandlesrWrapper) DeleteRequestForClient(ctx echo.Context) error {
	// функция для удаления просьбы
	item := ctx.Param("id")
	tempid, _ := strconv.Atoi(item)
	ok := wrapper.Delete(tempid)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}


func (wrapper *HandlesrWrapper) RequestFromClientHandler(ctx echo.Context) error {
	result := &store.ClientBody{}
	if err := ctx.Bind(result); err != nil {
		return err
	}
	resp, err := requester.RequestIssueExecutor(result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := struct {
		Id      int
		Status  int
		Headers map[string][]string
		Length  int64
	}{
		Headers: resp.Header,
		Status:  resp.StatusCode,
		Length:  resp.ContentLength,
	}
	res.Id = wrapper.Set(result)
	return ctx.JSON(http.StatusOK, res)
}
