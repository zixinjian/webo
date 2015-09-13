package main

import (
	"fmt"
	"reflect"
	"time"
	"webo/models/s"
)

func main() {
	fmt.Println(fmt.Sprintf("%v %v", "abc", 1), time.Now().Unix())
	sb := reflect.ValueOf(&s).Elem()
	fmt.Println(sb, s.Name)
}
