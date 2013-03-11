package logic

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	. "models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"regexp"
)

type LessCompile struct {
	Dir        map[string]bool                  //要监听的目录
	File       map[string]*LessFileState //文件
	IsCompress bool                      //是否压缩 -x
	w          *fsnotify.Watcher
	Suffix     string //扩展名
	c          chan int
	watch      map[string]bool
	reg			*regexp.Regexp
}
type LessFileState struct {
	Time  time.Time
	Error string
}

func (p *LessCompile) Close() {
	p.w.Close()
	p.c <- 1
}
func (p *LessCompile) Start() {
	go p.start()
}
func (p *LessCompile) start() {
	for k, _ := range p.Dir {
		log.Println("Find Path:", k)
		p.add(k)
	}
	p.doCompiler()
}
func (p *LessCompile) Add(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("不是一个有效的目录")
	}
	if _,ok:=p.Dir[dir];!ok {
		p.Dir[dir]=true
		err := p.add(dir)
		if err != nil {
			return err
		}
		p.Save()
	}else{
		return errors.New("已经存在")
	}
	return nil
}
func (p *LessCompile) Save() {
	model := ModelLess{}
	model.Save(p.Dir, p.IsCompress, p.Suffix)
}
func (p *LessCompile) add(dir string) (err error) {
	p.findFile(dir)
	return nil
}
func (p *LessCompile) Del(dir string) (err error) {
	err = p.deleteDir(dir)
	delete(p.Dir, dir)
	return nil
}
func (p *LessCompile) findFile(dir string) {
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			if v, ok := p.watch[path]; !ok || !v {
				log.Println("Add Watch:", path)
				err = p.w.Watch(path)
				p.watch[path] = true
			}
		} else {
			if strings.HasSuffix(path, p.Suffix) {
				p.File[path] = &LessFileState{}
				log.Println("Found File:",path)
				p.compiler(path)
			}
		}
		return nil
	})
}
func (p *LessCompile) doCompiler() {
	for {
		select {
		case v := <-p.w.Event:
			{
				log.Println("event:", v)
				fi, err := os.Stat(v.Name)
				if err == nil && fi.IsDir() && v.IsCreate() {
					p.add(v.Name)
				}
				if strings.HasSuffix(v.Name, p.Suffix) {
					log.Println("event:", v)
					if v.IsModify()||v.IsCreate() {
						if _,ok:=p.File[v.Name];!ok{
							p.File[v.Name] = &LessFileState{}
						}
						p.compiler(v.Name)
					}
					if v.IsRename()||v.IsDelete() {
						p.deleteFile(v.Name)
					}
				} else {
					if v.IsDelete() || v.IsRename() {
						log.Println("event:", v)
						p.deleteDir(v.Name)
						if v.IsRename(){
							log.Println("rename:",v.Name)
							fi,err:=os.Stat(v.Name)
							if err==nil&&fi.IsDir(){
								//支持windows中新建目录问题,以及目录改名
								p.add(v.Name)
							}
						}
					}
				}
			}
		case err := <-p.w.Error:
			log.Println("error:", err)
		case <-p.c:
			break
		}
	}
}
func (p *LessCompile) deleteDir(path string) (err error) {
	err = p.w.RemoveWatch(path)
	if v, ok := p.watch[path]; ok && v {
		if _,ok:=p.Dir[path];ok {
			p.Dir[path]=false
			p.Save()
		}
		for k, _ := range p.File {
			if strings.HasPrefix(k, path) {
				log.Println("delete file:", k)
				delete(p.File, k)
			}
		}
		p.watch[path] = false
		log.Println("delete dir:",path)
	}
	return
}
func (p *LessCompile) deleteFile(path string) {
	log.Println("delete file:", path)
	os.Remove(strings.Replace(path, fmt.Sprint(".", p.Suffix), ".css", -1))
	delete(p.File, path)
}
func (p *LessCompile) compiler(file string) {
	log.Println("build:", file)
	str := ""
	if p.IsCompress {
		//当有压缩参数的时候,文件的状态是两次修改
		str = "-x"
	}
	cmd := exec.Command("lessc", str, file, strings.Replace(file, fmt.Sprint(".", p.Suffix), ".css", -1))
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Run()
	st:=out.String()
	if st!=""{
		p.File[file].Error =strings.Replace(p.reg.ReplaceAllLiteralString(st, ""),file,"当前文件",-1)
	}
	p.File[file].Time = time.Now()
}
func (p *LessCompile) FindAll() {
	p.File = map[string]*LessFileState{}
	for k, _ := range p.Dir {
		p.findFile(k)
	}
}
func (p *LessCompile) CompileAll() {
	for k, _ := range p.File {
		p.compiler(k)
	}
}
func NewLessCompile() *LessCompile {
	model := ModelLess{}
	m := model.Get()
	le := LessCompile{}
	le.Dir = m.Dir
	le.File = map[string]*LessFileState{}
	le.IsCompress = m.Compress
	le.w, _ = fsnotify.NewWatcher()
	le.Suffix = m.Suffix
	le.watch = map[string]bool{}
	le.reg,_=regexp.Compile(`\[\d+?m`)
	return &le
}
