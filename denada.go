package main

import "fmt"
import "strings"

func (d Denada) String() string {
	ids := strings.Join(d.Ids, ",");
	return fmt.Sprintf("(Denada Ids:%s (%d))", ids, len(d.Ids));
}
