package reflectz

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MyStruct struct {
	Address string `foo:"customTag"`
	Email   string
	ID      int
	Info    MoreInfo
	Details []string
	Born    time.Time
}

type MoreInfo struct {
	ZipCode string
	State   string `reflectz:"-"`
	Name    string
}

func Test_RecursiveFalse(t *testing.T) {
	testStruct := &MyStruct{
		Address: "123",
		Email:   "foo@bar.com",
		Info:    MoreInfo{Name: "test"},
		Born:    time.Now(),
	}

	r := Inspect(testStruct, &Config{Recursive: false, NamingScheme: PascalCase})

	assert.Equal(t, 6, len(r), "Number of struct fields")
	address := r[0]
	assert.Equal(t, testStruct.Address, address.Value, "Values are equal")
	assert.Equal(t, "Address", address.FieldName, "FieldName is correct")
	assert.Equal(t, false, address.IsZero, "Field is not zero")
	assert.Equal(t, reflect.String, address.Kind, "Field is kind string")
	id := r[2]
	assert.Equal(t, testStruct.ID, id.Value, "Values are equal")
	assert.Equal(t, "ID", id.FieldName, "FieldName is correct")
	assert.Equal(t, true, id.IsZero, "Field is zero")
	assert.Equal(t, reflect.Int, id.Kind, "Field is kind int")
}

func Test_RecursiveTrue(t *testing.T) {
	testStruct := &MyStruct{
		Address: "123",
		Email:   "foo@bar.com",
		Info:    MoreInfo{Name: "test"},
		Born:    time.Now(),
	}

	r := Inspect(testStruct, &Config{Recursive: true, NamingScheme: SnakeCase, EmbeddedSep: "_"})

	assert.Equal(t, 7, len(r), "Number of struct fields")
	address := r[0]
	assert.Equal(t, testStruct.Address, address.Value, "Values are equal")
	assert.Equal(t, "address", address.FieldName, "FieldName is correct")
	assert.Equal(t, false, address.IsZero, "Field is not zero")
	assert.Equal(t, reflect.String, address.Kind, "Field is kind string")
	assert.Equal(t, []int{0}, address.Index, "Index is correct")
	id := r[2]
	assert.Equal(t, testStruct.ID, id.Value, "Values are equal")
	assert.Equal(t, "id", id.FieldName, "FieldName is correct")
	assert.Equal(t, true, id.IsZero, "Field is zero")
	assert.Equal(t, reflect.Int, id.Kind, "Field is kind int")
	assert.Equal(t, []int{2}, id.Index, "Index is correct")
	name := r[4]
	assert.Equal(t, testStruct.Info.Name, name.Value, "Values are equal")
	assert.Equal(t, "info_name", name.FieldName, "FieldName is correct")
	assert.Equal(t, false, name.IsZero, "Field is zero")
	assert.Equal(t, reflect.String, name.Kind, "Field is kind int")
	assert.Equal(t, []int{3, 1}, name.Index, "Index is correct")

}
