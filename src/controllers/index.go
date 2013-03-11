package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/iyf/gotool/jsonstate"
	. "github.com/iyf/gotool/middleware"
	"logic"
	"net/http"
	"os/exec"
	"strconv"
	"os"
)

type PageIndex struct {
	Application
}

func init() {
	App.RegisterController("index/", PageIndex{})
}

func (p *PageIndex) Init() {
	p.Application.Init()
}

func (p *PageIndex) Index(w http.ResponseWriter, r *http.Request) {
	rd := Data{}
	RunError := Middleware.Get("RunError")
	if RunError != nil {
		if t := p.GET["try"]; t != "" {
			cmd := exec.Command("lessc")
			err := cmd.Start()
			if err != nil {
				Middleware.Add("RunError", fmt.Sprint(err))
			} else {
				Middleware.Del("RunError")
				http.Redirect(w, r, "index.html", http.StatusFound)
				return
			}
		}
		rd.IsRun = false
		rd.Error = RunError.(string)
	} else {
		rd.IsRun = true

		trun := Middleware.Get("Run")
		if trun == nil {
			run := logic.NewLessCompile()
			run.Start()
			Middleware.Add("Run", run)
			trun = run
		}
		run := trun.(*logic.LessCompile)
		rd.Dir = run.Dir
		rd.File = map[string]*FileData{}
		for k, v := range run.File {
			fd:=FileData{}
			fd.Time = v.Time.Format("2006-01-02 15:04:05")
			if v.Error==""{
				fd.HasError=false
			}else{
				fd.HasError=true
			}
			fd.Error=v.Error
			rd.File[k]=&fd
		}
		rd.IsCompress = run.IsCompress
	}
	p.Body = rd
}
func (p *PageIndex) Add(w http.ResponseWriter, r *http.Request) {
	p.Hide = true
	st := jsonstate.BoolString{false, "参数错误"}
	if path := p.GET["path"]; path != "" {

		trun := Middleware.Get("Run")
		if trun == nil {
			st.S = false
			st.T = "未启动编译进程"
		} else {
			run := trun.(*logic.LessCompile)
			err := run.Add(path)
			if err != nil {
				st.S = false
				st.T = fmt.Sprint(err)
			} else {
				st.S = true
				st.T = "增加目录成功"
			}
		}
	}
	b, err := json.Marshal(st)
	if err == nil {
		w.Write(b)
	}
}
func (p *PageIndex) Del(w http.ResponseWriter, r *http.Request) {
	p.Hide = true
	st := jsonstate.BoolString{false, "参数错误"}
	if path := p.GET["path"]; path != "" {
		trun := Middleware.Get("Run")
		if trun == nil {
			st.S = false
			st.T = "未启动编译进程"
		} else {
			run := trun.(*logic.LessCompile)
			err := run.Del(path)
			if err != nil {
				st.S = false
				st.T = fmt.Sprint(err)
			} else {
				st.S = true
				st.T = "删除目录成功"
			}
		}
	}
	b, err := json.Marshal(st)
	if err == nil {
		w.Write(b)
	}
}
func (p *PageIndex) Set(w http.ResponseWriter, r *http.Request) {
	p.Hide = true
	st := jsonstate.BoolString{false, "参数错误"}
	if compress := p.GET["compress"]; compress != "" {
		b,_ := strconv.ParseBool(compress)
		trun := Middleware.Get("Run")
		if trun == nil {
			return
		}
		run := trun.(*logic.LessCompile)
		run.IsCompress = b
		run.Save()
		st.S=b
	}
	b, err := json.Marshal(st)
	if err == nil {
		w.Write(b)
	}
}
func (p *PageIndex) ScanCompile() {
	trun := Middleware.Get("Run")
		if trun == nil {
			return
		}
		run := trun.(*logic.LessCompile)
		run.FindAll()
		p.Hide=true
}
func (p *PageIndex) Compile() {
	trun := Middleware.Get("Run")
		if trun == nil {
			return
		}
		run := trun.(*logic.LessCompile)
		run.CompileAll()
		p.Hide=true

}
func (p *PageIndex) Close() {
	os.Exit(0)
}

type Data struct {
	IsRun      bool
	Dir        map[string]bool
	File       map[string]*FileData
	IsCompress bool
	Error      string
}
type FileData struct{
	Time	string
	HasError bool
	Error	string
}
