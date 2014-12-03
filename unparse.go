package denada

import "io"
import "fmt"
import "bytes"
import "strings"

import "github.com/bitly/go-simplejson"

func Unparse(elems ElementList) string {
	w := bytes.NewBuffer([]byte{})
	UnparseTo(elems, w)
	return w.String()
}

func UnparseTo(elems ElementList, w io.Writer) {
	unparse(elems, "", w)
}

func unparse(elems ElementList, prefix string, w io.Writer) {
	for _, e := range elems {
		unparseElement(*e, prefix, w)
		fmt.Fprintf(w, prefix+"\n")
	}
}

func unparseValue(v *simplejson.Json, prefix string) string {
	enc, err := v.EncodePretty()
	if err != nil {
		panic(err)
	}
	estr := string(enc)
	estr = strings.Replace(estr, "\n", "\n"+prefix, -1)
	return estr
}

func UnparseElement(e Element) string {
	w := bytes.NewBuffer([]byte{})
	unparseElement(e, "", w)
	return w.String()
}

func unparseElement(e Element, prefix string, w io.Writer) {
	fmt.Fprintf(w, prefix)
	for _, q := range e.Qualifiers {
		fmt.Fprintf(w, "%s ", q)
	}
	fmt.Fprintf(w, "%s", e.Name)
	if len(e.Modifications) > 0 {
		first := true
		fmt.Fprintf(w, "(")
		for k, v := range e.Modifications {
			if !first {
				fmt.Fprintf(w, ", ")
			}
			if v != nil {
				estr := unparseValue(v, prefix)
				fmt.Fprintf(w, "%s=%s", k, estr)
			}
			first = false
		}
		fmt.Fprintf(w, ")")
	}
	if e.isDefinition() {
		fmt.Fprintf(w, " ")
		if e.Description != "" {
			fmt.Fprintf(w, "\"%s\" ", strings.Replace(e.Description, "\"", "\\\"", 0))
		}
		fmt.Fprintf(w, "{\n")
		if e.Contents != nil {
			unparse(e.Contents, prefix+"  ", w)
		}
		fmt.Fprintf(w, "%s}", prefix)
	} else {
		if e.Value != nil {
			estr := unparseValue(e.Value, prefix)
			fmt.Fprintf(w, "=%s", estr)
		}
		if e.Description != "" {
			fmt.Fprintf(w, " \"%s\"", strings.Replace(e.Description, "\"", "\\\"", 0))
		}
		fmt.Fprintf(w, ";")
	}
}
