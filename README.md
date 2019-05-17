## vcontext: Validation with context

vcontext is a library supporting validation of config files parsed from yaml or json with support for giving context as
to where errors occur. It is composed of multiple semi-indepedent packages:

 - report: a structure for containing multiple errors, warnings, etc. It also contains context (from the path package)
   for determining where in the config the reports came from (e.g. `$.foo.baz.4.quux`) 
 - validate: a package for composing a report by validating go structs.
 - tree: a structure for containing metadata about the location (line/column) of objects in the source of the config
 - json, yaml: packages for generating trees from json or yaml
 - path: a structure for defining how to find json/yaml elements

### Usage:

Validating a config generally involves:
1) Unmarshal the yaml or json to a go struct, handle any syntax/type errors
1) Generate a report by running validate.Validate(yourConfigStruct)
1) Generate a tree of line/column metadata with [json|yaml].UnmarshalToContext()
1) Correlate the report to the metadata tree with report.Correlate(). This fills in line and column information from the tree.

The `validate` package does not require the structs came from json or yaml and will generate reports with context for each
error with a path like `$.foo.baz.4.quux`. The validate package walks the structs using reflection, aggregating the results
of calling any `Validate(c path.ContextPath) report.Report` functions defined on the types it is walking.

To write validation functions for your types, implement this interface:
```
Validate(path.ContextPath) report.Report
```
on any types you wish to be validated. The ContextPath passed to that function respresents the path to the type being
validated. ContextPath is a `[]interface{}` under the hood. You can append strings or ints to it to be more precise about
where problems occur. 

Example:
```go
type MyStruct struct {
	MustBePositive int `json:"mustBePosistive"`
}

func (m MyStruct) Validate(c path.ContextPath) (r report.Report) {
	if m.MustBePositive <= 0 {
		// append the json tag so the report specifies that field is invalid, not the whole struct
		r.AddOnError(append(c, "mustBePositive"), errors.New("mustBePositive was not positive"))
	}
	return
}
```

### Notes:

* This project under development and may undergo breaking changes
