package report

import (
	"fmt"

	"github.com/ajeddeloh/vcontext/path"
	"github.com/ajeddeloh/vcontext/tree"
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

// getDeepestNode 
func getDeepestNode(n tree.Node, c path.ContextPath) tree.Node {
	if child, err := n.Get(c); err != nil {
		return getDeepestNode(n, c[0:len(c)-1])
	} else {
		return child
	}
}

// Correlate takes a node tree and populates the markers
func (r *Report) Correlate(n tree.Node) {
	for i, e := range r.Entries {
		r.Entries[i].Marker = getDeepestNode(n, e.Context).GetMarker()
	}
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
	Marker  tree.Marker
}

func (e Entry) String() string {
	at := ""
	switch {
	case e.Marker.StartP != nil && e.Context != nil:
		at = fmt.Sprintf("at %s, %s", e.Context.String(), e.Marker.String())
	case e.Marker.StartP != nil:
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

func (r *Report) add(c path.ContextPath, err error, k Kind) {
	if err == nil {
		return
	}
	r.Entries = append(r.Entries, Entry{
		Message: err.Error(),
		Context: c,
		Kind:    k,
	})
}

func (r *Report) AddOnError(c path.ContextPath, err error) {
	r.add(c, err, Error)

}
func (r *Report) AddOnWarn(c path.ContextPath, err error) {
	r.add(c, err, Warn)
}

func (r *Report) AddOnInfo(c path.ContextPath, err error) {
	r.add(c, err, Info)
}
