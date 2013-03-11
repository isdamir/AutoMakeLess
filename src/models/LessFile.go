package models 
import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
)
type ModelLess struct{

}
func (p *ModelLess) Get()*Message{
	bi,err:= ioutil.ReadFile("data/data.json")
	fmt.Println(err)
	if err!=nil{
		return &Message{map[string]bool{},false,"less"}
	}
	var ki Message
	err=json.Unmarshal(bi, & ki)
	if err!=nil{
		return &Message{map[string]bool{},false,"less"}
	}
	for k,_:=range ki.Dir{
		fi,err:=os.Stat(k)
		if err==nil&&fi.IsDir(){
			ki.Dir[k]=true
		}else{
			ki.Dir[k]=true
		}
	}
	return &ki
}
func (p *ModelLess) Save(dir map[string]bool,c bool,suffix string){
	m := Message{dir,c,suffix}
	b, _ := json.Marshal(m)
	ioutil.WriteFile( "data/data.json" , b, 0644 )//保存数据到文件	
}
type Message struct {
	Dir map[string]bool
	Compress bool
	Suffix  string
}
