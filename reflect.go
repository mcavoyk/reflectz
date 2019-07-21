package reflectz

import (
	"reflect"

	"github.com/iancoleman/strcase"
)

type StructField struct {
	FieldName string
	Kind      reflect.Kind
	Type      reflect.Type
	Tags      reflect.StructTag
	Index     []int
	Value     interface{}
	IsZero    bool
}

type Config struct {
	// If Recursive is false Inspect will not dive embedded structs
	// true by default
	Recursive bool

	// EmbeddedSep controls the separator between embedded structs
	// underscore by default
	EmbeddedSep string

	// NamingScheme controls the naming convention of struct fields
	// Options include snake_case, PascalCase or camelCase
	// snake_case by default
	NamingScheme NameType

	// NamingTag controls the name of struct tag which should be used
	// for field names, for example 'json'. If the tag is not present
	// on a field, the name will fallback to the specified NamingScheme
	NamingTag string
}

type NameType int

const (
	SnakeCase  NameType = 0
	CamelCase  NameType = 1
	PascalCase NameType = 2
)

// Tag is the name of struct tag used for reflectz directives
// "-" and "nodive" are the current tag options available
var Tag = "reflectz"

func Inspect(model interface{}, opt ...*Config) []StructField {
	var config *Config
	if len(opt) == 0 {
		config = defaultOptions()
	} else {
		config = opt[0]
	}
	return inspect(model, config)
}

func defaultOptions() *Config {
	return &Config{
		Recursive:    true,
		EmbeddedSep:  "_",
		NamingScheme: SnakeCase,
	}
}

func inspect(model interface{}, opt *Config) []StructField {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Inspect must be passed a struct or pointer to a struct")
	}
	r := make([]StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		fieldT := t.Field(i)
		// Skipping unexported fields
		if fieldT.PkgPath != "" || fieldT.Tag.Get(Tag) == "-" {
			//fmt.Printf("Skipping - %+v\n", fieldT)
			continue
		}
		kind := fieldT.Type.Kind()
		name := fieldName(fieldT, opt)
		fieldV := v.Field(i)
		value := fieldV.Interface()
		index := []int{i}

		// Dive embedded structs
		noDive := fieldT.Tag.Get(Tag) == "nodive"
		if kind == reflect.Struct && opt.Recursive && !noDive {
			embedded := inspect(value, opt)

			for i, e := range embedded {
				e.FieldName = name + opt.EmbeddedSep + e.FieldName
				e.Index = append(index, i)
				r = append(r, e)
			}
			if len(embedded) != 0 {
				continue
			}
		}

		isZero := reflect.DeepEqual(reflect.Zero(fieldT.Type).Interface(), value)
		r = append(r, StructField{
			FieldName: name,
			Type:      fieldT.Type,
			Tags:      fieldT.Tag,
			Kind:      kind,
			Value:     value,
			IsZero:    isZero,
			Index:     index,
		})
	}
	return r
}

func fieldName(field reflect.StructField, opt *Config) string {
	if tagName, tagExists := field.Tag.Lookup(opt.NamingTag); tagExists {
		return tagName
	}
	name := field.Name
	switch opt.NamingScheme {
	case SnakeCase:
		return strcase.ToSnake(name)
	case CamelCase:
		return strcase.ToLowerCamel(name)
	case PascalCase:
		return strcase.ToCamel(name)
	default:
		panic("Unknown NamingScheme: " + string(opt.NamingScheme))
	}
}
