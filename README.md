# Reflectz

Reflectz is a small library which provides deep inspection of struct fields.

[![GoDoc](https://godoc.org/github.com/mcavoyk/reflectz?status.svg)](https://godoc.org/github.com/mcavoyk/reflectz)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

```go
type MyStruct struct {
	ID    int     `json:"id"`
	Email string  `json:"email"`
	Info  Another `json:"info"`
}

type Another struct {
	Address string `reflectz:"-"`
	Details []string `json:"list"`
}

func main() {
	myStruct := MyStruct{
		ID:    10,
		Email: "foo@bar.com",
	}

	opts := &reflectz.Config{
		Recursive:   true,
		NamingTag:   "json",
		EmbeddedSep: "_",
	}

	for _, r := range reflectz.Inspect(&myStruct, opts) {
		fmt.Printf("%+v\n", r)
	} 
	// {FieldName:id Kind:int Type:int Tags:json:"id" Index:[0] Value:10 IsZero:false}
	// {FieldName:email Kind:string Type:string Tags:json:"email" Index:[1] Value:foo@bar.com IsZero:false}
	// {FieldName:info_list Kind:slice Type:[]string Tags:json:"list" Index:[2 0] Value:[] IsZero:true}
}
```
## Configuration
Reflectz offers multiple methods of naming field results, snake_case, PascalCase or camelCase.  The config
option `NamingTag` can also be used to allow field names by a specified struct tag (e.g. `json`)

Reflectz offers the following struct tag options:

- `reflectz:"-"` Does not include field or field children in results.
- `reflectz:"nodive"` Does not dive and inspect field children (Set `Recursive: false` to enable this on all struct subfields.)

The `Index` field returns a slice of integers which can be used to dive
back to an embedded field within the struct.

```Go
type Example struct {
	Data Embedded
}
type Embedded struct {
	Value string
}
// Example setting a field
e := Example{Data: Embedded{Value: "original"}}
structs := reflectz.Inspect(&e)
embeddedField := reflect.ValueOf(&e).Elem().FieldByIndex(structs[0].Index)
embeddedField.Set(reflect.ValueOf("newValue"))

fmt.Println(e)
// {{newValue}}
```