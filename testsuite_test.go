package denada

import "os"
import "log"
import "fmt"
import "path"
import "testing"
import "strings"

import . "github.com/onsi/gomega"

func CheckFile(name string) error {
	filename := path.Join("testsuite", name)

	elems, err := ParseFile(filename)
	if err != nil {
		return err
	}

	if len(elems) == 0 {
		return fmt.Errorf("Empty file")
	}
	props, elems, err := elems.PopHead()
	if err != nil {
		return fmt.Errorf("Missing properties")
	}

	declsv, exists := props.Modifications["declarations"]
	edecls := 0
	if exists {
		edecls = declsv.MustInt(0)
	}

	defsv, exists := props.Modifications["definitions"]
	edefs := 0
	if exists {
		edefs = defsv.MustInt(0)
	}

	adecls := 0
	adefs := 0
	for _, e := range elems {
		if e.isDeclaration() {
			adecls++
		}
		if e.isDefinition() {
			adefs++
		}
	}

	if adecls != edecls {
		return fmt.Errorf("Expected %d declarations, found %d", edecls, adecls)
	}

	if adefs != edefs {
		return fmt.Errorf("Expected %d declarations, found %d", edecls, adecls)
	}

	return nil
}

func CheckError(name string) {
	err := CheckFile(name)
	Expect(err).ToNot(BeNil())
}

func Test_TestSuite(t *testing.T) {
	RegisterTestingT(t)

	cur, err := os.Open("testsuite")
	Expect(err).To(BeNil())

	files, err := cur.Readdir(0)
	Expect(err).To(BeNil())

	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".dnd") {
			if strings.HasPrefix(name, "case") {
				err := CheckFile(name)
				/*
					if err == nil {
						log.Printf("Case %s: PASSED", name)
					} else {
						log.Printf("Case %s: Failed: %v", name, err)
					}
				*/
				Expect(err).To(BeNil())
				continue
			}
			if strings.HasPrefix(name, "ecase") {
				err := CheckFile(name)
				/*
					if err != nil {
						log.Printf("Error Case %s: PASSED", name)
					} else {
						log.Printf("Error Case %s: FAILED", name)
					}
				*/
				Expect(err).ToNot(BeNil())
				continue
			}
			log.Printf("Unrecognized file type in test suite: %s", name)
		}
	}
}
