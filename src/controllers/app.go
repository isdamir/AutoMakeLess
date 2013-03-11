package controllers

import (
	"github.com/iyf/gotool/web"
	"net/http"
)

type Application struct {
	web.Page
	OffLogin bool
	OffRight bool
	RW       http.ResponseWriter
	R        *http.Request
}

var App = &Application{
	Page: web.NewPage(web.PageParam{TimerDuration: "2h"}),
}

func (a *Application) Init() {
	a.Page.Init(a.RW, a.R)
}
