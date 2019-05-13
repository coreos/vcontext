package report

import (
	"fmt"

	"github.com/ajeddeloh/vcontext"
	"github.com/ajeddeloh/vcontext/path"
)

type EntryKind interface {
	String() string
	IsFatal() bool
}

type Report struct {
	Entries []Entry
}

func (r *Report) Merge(child Report) {
	r.Entries = append(r.Entries, child.Entries...)
}

func (r Report) IsFatal() bool {
	for _, e := range r.Entries {
		if e.Kind.IsFatal() {
			return true
		}
	}
	return false
}

func (r Report) String() string {
	str := ""
	for _, e := range r.Entries {
		str += e.String() + "\n" 
	}
	return str
}

type Entry struct {
	Kind    EntryKind
	Message string
	Context path.ContextPath
	Marker  vcontext.Marker
}

func (e Entry) String() string {
	at := ""
	switch {
	case e.Marker != nil && e.Context != nil:
		at = fmt.Sprintf("at %s, %s", e.Context.String(), e.Marker.String())
	case e.Marker != nil:
		at = fmt.Sprintf("at %s", e.Marker.String())
	case e.Context != nil:
		at = fmt.Sprintf("at %s", e.Context.String())
	}

	return fmt.Sprintf("%s %s: %s", e.Kind.String(), at, e.Message)
}

type Kind int

const (
	Error Kind = iota
	Warn  Kind = iota
	Info  Kind = iota
)

func (k Kind) String() string {
	switch k {
	case Error:
		return "error"
	case Warn:
		return "warning"
	case Info:
		return "info"
	default:
		return ""
	}
}

func (k Kind) IsFatal() bool {
	return k == Error
}

// Helpers
func FromError(c path.ContextPath, err error) Report {
	return Report{
		Entries: []Entry{
			{
				Message: err.Error(),
				Context: c,
				Kind:    Error,
			},
		},
	}
}
