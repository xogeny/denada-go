package denada

import "testing"
import "log"
import "strings"

var sample = `
printer ABC {
   set location = "Mike's desk";
   set model = "HP 8860";
}

printer DEF {
   set location = "Coffee machine";
   set model = "HP 8860";
   set networkName = "PrinterDEF";
}

computer XYZ {
   set location = "Mike's desk";
   set model = "Mac Book Air";
}
`

func Test1(t *testing.T) {
	r := strings.NewReader("x = 5;")
	err := Parse(r)
	if err != nil {
		log.Printf("Error: %x", err)
	}
	_dump()
}
