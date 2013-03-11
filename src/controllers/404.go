package controllers

import (
	"net/http"
)

type Page404 struct {
	Application
}

func init() {
	App.SetNotFoundController(Page404{})
}

func (p *Page404) Init() {
	p.Application.Init()
	p.Document.GenerateHtml = false
	p.Template = "_notfound/404.html"
	p.RW.WriteHeader(http.StatusNotFound)
}
