package path

import (
	"fmt"
	"strings"
)

type ContextPath []interface{}

func (c ContextPath) String() string {
	strs := []string{"$"}
	for _, e := range c {
		strs = append(strs, fmt.Sprintf("%v", e))
	}
	return strings.Join(strs, ".")
}
