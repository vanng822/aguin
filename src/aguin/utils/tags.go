package utils
import (
	"reflect"
	"fmt"
)

type Tags map[string]string

func (t Tags) Get(name string) string {
	if tag, ok := t[name]; ok {
		return tag
	}
	panic(fmt.Sprintf("There no tag for %s", name))
}

func GetFieldsTag(o interface{}, tagname string) Tags {
	st := reflect.TypeOf(o)
	tags := Tags{}
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		value := field.Tag.Get(tagname)
		if value != "" {
			tags[field.Name] = value
		}
	}
	return tags
}
