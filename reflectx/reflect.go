package reflectx

import (
	"reflect"
)

func GetTag(field reflect.StructField, tags ...string) (content string) {
	for content == "" && len(tags) > 0 {
		tag := tags[0]
		tags = tags[1:]
		content = field.Tag.Get(tag)
	}

	return
}
