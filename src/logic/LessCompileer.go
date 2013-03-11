package logic

type LessCompileer interface{
	Start()
	Close()
	Add(dir string) error
	Del(dir string) error
	FindALL()
	CompileAll()
}
