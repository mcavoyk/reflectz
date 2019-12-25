package main

import (
	"fmt"
	"reflect"

	"github.com/mcavoyk/reflectz"
)

type MyStruct struct {
	ID    int     `json:"id"`
	Email string  `json:"email,omitempty"`
	Info  Another `json:"info"`
}

type Another struct {
	Address string `reflectz:"-"`
	Extra   int
	Details []string `json:"list"`
}

type Example struct {
	Data Embedded
}
type Embedded struct {
	Value string
}

func main() {
	myStruct := MyStruct{
		ID:    10,
		Email: "foo@bar.com",
	}

	opts := &reflectz.Config{
		NamingTag:   "json",
		Recursive:   true,
		EmbeddedSep: "_",
	}

	for _, r := range reflectz.Inspect(&myStruct, opts) {
		fmt.Printf("%+v\n", r)
	}
	// {FieldName:id Kind:int Type:int Tags:json:"id" Index:[0] Value:10 IsZero:false}
	// {FieldName:email Kind:string Type:string Tags:json:"email" Index:[1] Value:foo@bar.com IsZero:false}
	// {FieldName:info_extra Kind:int Type:int Tags: Index:[2 1] Value:0 IsZero:true}
	// {FieldName:info_list Kind:slice Type:[]string Tags:json:"list" Index:[2 0] Value:[] IsZero:true}

	// Example setting a field
	e := Example{Data: Embedded{Value: "original"}}
	structs := reflectz.Inspect(&e)
	embeddedField := reflect.ValueOf(&e).Elem().FieldByIndex(structs[0].Index)
	embeddedField.Set(reflect.ValueOf("newValue"))

	fmt.Println(e)
}
