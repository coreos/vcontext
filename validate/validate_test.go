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
	"testing"
)

func TestGetFields(t *testing.T) {
	// basic case
	type Test1 struct {
		A int
		B string
	}
	test1 := Test1{
		A: 1,
		B: "one",
	}

	// test embedded structs
	type Test2 struct {
		C int
		Test1
	}
	test2 := Test2{
		C:     5,
		Test1: test1,
	}

	// test doublely embedded structs
	type Test3 struct {
		D int
		Test2
	}
	test3 := Test3{
		D:     3,
		Test2: test2,
	}
	// test structs embedded via an alias to interface{}
	type Anything interface{}

	test4 := struct {
		E int
		Anything
	}{
		E:        7,
		Anything: test3,
	}

	// test normally contained structs don't cause problems
	test5 := struct {
		E int
		F Test3
	}{
		E: 2,
		F: test3,
	}

	// test non-structs embedded via an alias to interface{} don't cause panics
	test6 := struct {
		E int
		Anything
	}{
		E:        5,
		Anything: 65,
	}

	// test embedded nils
	test7 := struct {
		E int
		Anything
	}{
		E: 5,
	}

	tests := []struct {
		in  reflect.Value
		out []StructField
	}{
		{
			reflect.ValueOf(test1),
			[]StructField{
				{
					StructField: reflect.TypeOf(test1).Field(0),
					Value:       reflect.ValueOf(test1.A),
				},
				{
					StructField: reflect.TypeOf(test1).Field(1),
					Value:       reflect.ValueOf(test1.B),
				},
			},
		},
		{
			reflect.ValueOf(test2),
			[]StructField{
				{
					StructField: reflect.TypeOf(test2).Field(0),
					Value:       reflect.ValueOf(test2.C),
				},
				{
					StructField: reflect.TypeOf(test1).Field(0),
					Value:       reflect.ValueOf(test1.A),
				},
				{
					StructField: reflect.TypeOf(test1).Field(1),
					Value:       reflect.ValueOf(test1.B),
				},
			},
		},
		{
			reflect.ValueOf(test3),
			[]StructField{
				{
					StructField: reflect.TypeOf(test3).Field(0),
					Value:       reflect.ValueOf(test3.D),
				},
				{
					StructField: reflect.TypeOf(test2).Field(0),
					Value:       reflect.ValueOf(test2.C),
				},
				{
					StructField: reflect.TypeOf(test1).Field(0),
					Value:       reflect.ValueOf(test1.A),
				},
				{
					StructField: reflect.TypeOf(test1).Field(1),
					Value:       reflect.ValueOf(test1.B),
				},
			},
		},
		{
			reflect.ValueOf(test4),
			[]StructField{
				{
					StructField: reflect.TypeOf(test4).Field(0),
					Value:       reflect.ValueOf(test4.E),
				},
				{
					StructField: reflect.TypeOf(test3).Field(0),
					Value:       reflect.ValueOf(test3.D),
				},
				{
					StructField: reflect.TypeOf(test2).Field(0),
					Value:       reflect.ValueOf(test2.C),
				},
				{
					StructField: reflect.TypeOf(test1).Field(0),
					Value:       reflect.ValueOf(test1.A),
				},
				{
					StructField: reflect.TypeOf(test1).Field(1),
					Value:       reflect.ValueOf(test1.B),
				},
			},
		},
		{
			reflect.ValueOf(test5),
			[]StructField{
				{
					StructField: reflect.TypeOf(test5).Field(0),
					Value:       reflect.ValueOf(test5.E),
				},
				{
					StructField: reflect.TypeOf(test5).Field(1),
					Value:       reflect.ValueOf(test5.F),
				},
			},
		},
		{
			reflect.ValueOf(test6),
			[]StructField{
				{
					StructField: reflect.TypeOf(test6).Field(0),
					Value:       reflect.ValueOf(test6.E),
				},
				{
					StructField: reflect.TypeOf(test6).Field(1),
					Value:       reflect.ValueOf(65),
				},
			},
		},
		{
			reflect.ValueOf(test7),
			[]StructField{
				{
					StructField: reflect.TypeOf(test7).Field(0),
					Value:       reflect.ValueOf(test7.E),
				},
				{
					StructField: reflect.TypeOf(test7).Field(1),
					Value:       reflect.ValueOf(nil),
				},
			},
		},
	}

	for i, test := range tests {
		fields := GetFields(test.in)
		// We cannot use reflect.DeepEqual because reflect.DeepEqual(reflect.ValueOf(someinstance),reflect.ValueOf(someinstance))
		// will always return false. We must manually loop over it and convert reflect.Value's to interface{}'s as we go
		for idx, f := range fields {
			if !reflect.DeepEqual(f.Type, test.out[idx].Type) {
				t.Errorf("#%d: bad error with type: want \n%+v, got \n%+v", i, fields, test.out)
			}
			if !reflect.DeepEqual(f.Value.Interface(), test.out[idx].Value.Interface()) {
				t.Errorf("#%d: bad error: want \n%+v, got \n%+v", i, fields, test.out)
			}
		}
	}
}
