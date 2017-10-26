package api

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/karlockhart/gorunremote/pkg/executor"
	"github.com/labstack/echo"
)

type GoRunRemoteApi struct {
	ech *echo.Echo
	exe *executor.SerialExecutor
}

func NewGoRunRemoteApi() *GoRunRemoteApi {
	api := GoRunRemoteApi{}
	api.ech = echo.New()

	api.exe = executor.NewSerialExecutor()
	api.ech.Static("/", "app/dist")
	api.ech.GET("api/go/load/", api.Load)
	api.ech.POST("api/go/fmt", api.Format)
	api.ech.POST("api/go/run", api.Run)

	return &api
}
func (a *GoRunRemoteApi) Start(wg *sync.WaitGroup) {

	a.ech.Logger.Error(a.ech.Start(":1323"))
	wg.Done()
}

func (a *GoRunRemoteApi) Load(c echo.Context) error {
	hash := c.FormValue("hash")

	f, e := a.exe.Load(hash)
	if e != nil {
		return c.String(http.StatusInternalServerError, e.Error())
	}
	return c.JSON(http.StatusOK, f)
}

func (a *GoRunRemoteApi) Format(c echo.Context) error {
	b, e := ioutil.ReadAll(c.Request().Body)
	if e != nil {
		return c.String(http.StatusInternalServerError, e.Error())
	}
	f, e := a.exe.Format(b)
	if e != nil {
		return c.String(http.StatusInternalServerError, e.Error())
	}
	return c.JSON(http.StatusOK, f)
}

func (a *GoRunRemoteApi) Run(c echo.Context) error {
	b, e := ioutil.ReadAll(c.Request().Body)
	if e != nil {
		return c.String(http.StatusInternalServerError, e.Error())
	}
	f, e := a.exe.Run(b)
	if e != nil {
		return c.String(http.StatusInternalServerError, e.Error())
	}
	return c.JSON(http.StatusOK, f)
}
