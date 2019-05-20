// Copyright 2019 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.)

package validate

import (
	"reflect"
	"strings"

	"github.com/ajeddeloh/vcontext/path"
	"github.com/ajeddeloh/vcontext/report"
)

type validator interface {
	Validate(path.ContextPath) report.Report
}

func Validate(thing interface{}, tag string) report.Report {
	if thing == nil {
		return report.Report{}
	}
	v := reflect.ValueOf(thing)
	ctx := path.ContextPath{Tag: tag}
	return validate(ctx, v)
}

func validate(context path.ContextPath, v reflect.Value) (r report.Report) {
	if !v.IsValid() {
		return
	}
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		} else {
			v = makeConcrete(v)
		}
	}
	// first check if this object has Validate(context) defined, but only on value
	// recievers. Both pointer and value receivers satisfy a value receiver interface
	// so ensure we're not a pointer too.
	if obj, ok := v.Interface().(validator); ok && v.Kind() != reflect.Ptr {
		r.Merge(obj.Validate(context))
	}

	switch v.Kind() {
	case reflect.Struct:
		r.Merge(validateStruct(context, v))
	case reflect.Slice:
		r.Merge(validateSlice(context, v))
	case reflect.Ptr:
		if !v.IsNil() {
			r.Merge(validate(context, v.Elem()))
		}
	}

	return
}

type structField struct {
	reflect.StructField
	Value reflect.Value
}

// makeConcrete takes a value and if it is a value of an interface returns the
// value of the actual underlying type implementing that interface. If the value
// is already concrete, it returns the same value
func makeConcrete(v reflect.Value) reflect.Value {
	return reflect.ValueOf(v.Interface())
}

func getFields(v reflect.Value) []structField {
	ret := []structField{}
	if v.Kind() != reflect.Struct {
		return ret
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if !field.Anonymous {
			ret = append(ret, structField{
				StructField: field,
				Value:       v.Field(i),
			})
		} else {
			concrete := makeConcrete(v.Field(i))
			ret = append(ret, getFields(concrete)...)
		}
	}
	return ret
}

func fieldName(s structField, tag string) string {
	if tag == "" {
		return s.Name
	}
	tag = s.Tag.Get(tag)
	return strings.Split(tag, ",")[0]
}

func validateStruct(context path.ContextPath, v reflect.Value) (r report.Report) {
	fields := getFields(v)
	for _, field := range fields {
		fieldContext := context.Append(fieldName(field, context.Tag))
		r.Merge(validate(fieldContext, field.Value))
	}
	return
}

func validateSlice(context path.ContextPath, v reflect.Value) (r report.Report) {
	for i := 0; i < v.Len(); i++ {
		childContext := context.Append(i)
		r.Merge(validate(childContext, v.Index(i)))
	}
	return
}
