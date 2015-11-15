package ddb

const (
	Dict="dict"
	List="list"
)

type FieldInfo struct {
	Key  	string		`json:key`
	Type    string		`json:type`
	Require bool		`json:require`
	Enum    []string	`json:enum`
}

type DocInfo struct {
	Path     string
	Name	 string
	Type	 string
	Fields   []FieldInfo
}


type Db struct{
	Name string
	Path string
	Docs []interface{}
}

type OptValue map[string]interface{}
type Options struct {
	Name    string
	Opts    []OptValue
}

var DbMap map[string]Db

func Load(){

}

func Insert(optionsName string, record OptValue){

}
func Append(optionsName string, record OptValue){

}
func Delete(idx int64){

}
func Save(optionsName string){

}


func Open(path string)error{

	return nil
}