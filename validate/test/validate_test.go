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
	"errors"
	"reflect"
	"testing"

	"github.com/ajeddeloh/vcontext/path"
	"github.com/ajeddeloh/vcontext/report"
	"github.com/ajeddeloh/vcontext/validate"
)

// Why are these tests in their own package?
// Go doesn't let you call reflect.Value.Interface() on non-public types but we
// need to test that. This package prevents those types from needing to be
// included in the exported package that people use.

var (
	dummy = errors.New("dummy")
)

func fromSingleError(c []interface{}, err error) (r report.Report) {
	ctx := path.ContextPath{
		Path: c,
	}
	r.AddOnError(ctx, err)
	return
}

type Test1 struct{}

func (t Test1) Validate(c path.ContextPath) (r report.Report) {
	return
}

type Test2 struct{}

func (t Test2) Validate(c path.ContextPath) (r report.Report) {
	r.AddOnError(c, dummy)
	return
}

type Test3 struct {
	Foo Test2 `yaml:"yaml,garbage" json:"json,garbage"`
}

type Test4 struct {
	Test2 `yaml:",inline"`
}

type Test5 struct {
	Bar Test3 `yaml:"yaml2,garbage" json:"json2,garbage"`
}

type Test6 struct {
	Bar *Test3 `yaml:"yaml2,garbage" json:"json2,garbage"`
}

type Test7 struct {
	Bar []*Test2 `yaml:"yaml" json:"json"`
}

func TestValidate(t *testing.T) {
	type test struct {
		in  interface{}
		src string
		out report.Report
	}
	tests := []test{
		{
			in: nil,
		},
		{
			in: struct{}{},
		},
		{
			in: []int{},
		},
		{
			in: Test1{},
		},
		{
			in:  Test2{},
			out: fromSingleError(nil, dummy),
		},
		{
			in:  Test3{},
			out: fromSingleError([]interface{}{"Foo"}, dummy),
		},
		{
			in:  Test3{},
			src: "json",
			out: fromSingleError([]interface{}{"json"}, dummy),
		},
		{
			in:  Test3{},
			src: "yaml",
			out: fromSingleError([]interface{}{"yaml"}, dummy),
		},
		{
			in:  Test4{},
			out: fromSingleError(nil, dummy),
		},
		{
			in:  Test4{},
			src: "json",
			out: fromSingleError(nil, dummy),
		},
		{
			in:  Test4{},
			src: "yaml",
			out: fromSingleError(nil, dummy),
		},
		{
			in:  Test5{},
			out: fromSingleError([]interface{}{"Bar", "Foo"}, dummy),
		},
		{
			in:  Test5{},
			src: "json",
			out: fromSingleError([]interface{}{"json2", "json"}, dummy),
		},
		{
			in:  Test5{},
			src: "yaml",
			out: fromSingleError([]interface{}{"yaml2", "yaml"}, dummy),
		},
		{
			in: Test6{},
		},
		{
			in:  Test6{Bar: new(Test3)},
			out: fromSingleError([]interface{}{"Bar", "Foo"}, dummy),
		},
		{
			in:  Test6{Bar: new(Test3)},
			src: "json",
			out: fromSingleError([]interface{}{"json2", "json"}, dummy),
		},
		{
			in:  Test6{Bar: new(Test3)},
			src: "yaml",
			out: fromSingleError([]interface{}{"yaml2", "yaml"}, dummy),
		},
		{
			in: Test7{},
		},
		{
			in:  Test7{Bar: []*Test2{new(Test2)}},
			out: fromSingleError([]interface{}{"Bar", 0}, dummy),
		},
		{
			in:  Test7{Bar: []*Test2{new(Test2)}},
			src: "json",
			out: fromSingleError([]interface{}{"json", 0}, dummy),
		},
		{
			in:  Test7{Bar: []*Test2{new(Test2)}},
			src: "yaml",
			out: fromSingleError([]interface{}{"yaml", 0}, dummy),
		},
	}
	for i, test := range tests {
		expected := test.out
		for i, _ := range test.out.Entries {
			test.out.Entries[i].Context.Tag = test.src
		}

		actual := validate.Validate(test.in, test.src)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Fail %da: expected %+v got %+v", i, expected, actual)
		}

		// now test on the pointer version for good measure
		actual = validate.Validate(&test.in, test.src)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Fail %db: expected %+v got %+v", i, expected, actual)
		}
	}
}
