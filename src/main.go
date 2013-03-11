package main

import (
	. "controllers"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	. "github.com/iyf/gotool/browser"
	. "github.com/iyf/gotool/middleware"
	"os/exec"
	"logic"
)

var (
	addr      = flag.String("addr", ":8080", "Server port")
	configDir = flag.String("config", "./config", "Directory of config")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	flag.Parse()
	os.Chdir(filepath.Dir(os.Args[0]))
	fmt.Println("Listen server address: " + *addr)
	fmt.Println("Read configuration directory success, directory: " + filepath.Join(filepath.Dir(os.Args[0]), *configDir))
	cmd := exec.Command("lessc");
	err := cmd.Start()
	if err!=nil{
		Middleware.Add("RunError", fmt.Sprint(err))
	}
	App.Load(*configDir)
	App.AddHeader("Content-Type", "text/html; charset=utf-8")
	b:=Browser{}
	b.OpenBrowserAsync(fmt.Sprint("http://localhost",*addr))
	App.ListenAndServe(*addr, App)
	defer Middleware.Get("Run").(*logic.LessCompile).Close()
}
