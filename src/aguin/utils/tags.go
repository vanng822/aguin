package utils
import (
	"reflect"
	"fmt"
)

type Tags map[string]string

func (t Tags) Get(name string) string {
	tag, ok := t[name]
	if !ok {
		panic(fmt.Sprintf("There no tag for %s", name))
	}
	return tag
}

func GetFieldsTag(o interface{}, tagname string) Tags {
	st := reflect.TypeOf(o)
	tags := make(map[string]string, st.NumField())	
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tags[field.Name] = field.Tag.Get(tagname)
	}
	return tags
}
