package denada

import "os"
import "log"
import "fmt"
import "path"
import "testing"
import "strings"

import . "github.com/smartystreets/goconvey/convey"

func ReparseFile(name string) error {
	filename := path.Join("testsuite", name)

	elems, err := ParseFile(filename)
	if err != nil {
		return err
	}

	str := Unparse(elems, false)

	relems, err := ParseString(str)
	if err != nil {
		return fmt.Errorf("Error in unparsed code: %v", err)
	}

	err = elems.Equals(relems)
	if err != nil {
		return fmt.Errorf("Inequality in reparsing of %s: %v", filename, err)
	}
	return nil
}

func CheckFile(name string) error {
	filename := path.Join("testsuite", name)

	elems, err := ParseFile(filename)
	if err != nil {
		return err
	}

	if len(elems) == 0 {
		return fmt.Errorf("Empty file")
	}

	props := elems[0]

	declsv, exists := props.Modifications["declarations"]
	var edecls int = 0
	if exists {
		edecls = declsv.MustInt(0)
	}

	defsv, exists := props.Modifications["definitions"]
	var edefs int = 0
	if exists {
		edefs = defsv.MustInt(0)
	}

	var adecls int = 0
	var adefs int = 0
	for _, e := range elems {
		if e.IsDeclaration() {
			adecls++
		}
		if e.IsDefinition() {
			adefs++
		}
	}

	if adecls != edecls {
		return fmt.Errorf("Expected %d declarations, found %d", edecls, adecls)
	}

	if adefs != edefs {
		return fmt.Errorf("Expected %d definitions, found %d", edefs, adefs)
	}

	grmv, exists := props.Modifications["grammar"]
	if exists {
		gfile := grmv.MustString()
		g, err := ParseFile(path.Join("testsuite", gfile))
		if err != nil {
			return err
		}
		err = Check(elems, g, false)
		// Check if descriptions on input elements matche expected rule names
		if err == nil {
			for _, e := range elems.AllElements() {
				if e.Description == "" {
					err = fmt.Errorf("Input element %v in %s didn't have a description",
						e, name)
					return err
				} else if e.RulePath() == "" {
					err = fmt.Errorf("Input element %v in %s didn't seem to match anything",
						e, name)
					return err
				} else {
					if e.RulePath() != e.Description {
						err = fmt.Errorf("Input element %v matched rule %s but description implies a match with %s", e, e.RulePath(), e.Description)
						return err
					}
				}
			}
		}
		return err
	}

	return nil
}

func CheckError(name string) {
	err := CheckFile(name)
	So(err, ShouldBeNil)
}

func Test_TestSuite(t *testing.T) {
	Convey("Running TestSuite", t, func() {
		cur, err := os.Open("testsuite")
		So(err, ShouldBeNil)

		files, err := cur.Readdir(0)
		So(err, ShouldBeNil)

		for _, f := range files {
			name := f.Name()
			if !strings.HasSuffix(name, ".dnd") {
				continue
			}
			Convey("Processing "+name, func() {
				if strings.HasPrefix(name, "case") {
					err := CheckFile(name)
					if err != nil {
						log.Printf("Case %s: Failed: %v", name, err)
					}
					So(err, ShouldBeNil)
					err = ReparseFile(name)
					if err != nil {
						log.Printf("Case %s: Reparse failed: %v", name, err)
					}
					So(err, ShouldBeNil)
					return
				}
				if strings.HasPrefix(name, "ecase") {
					err := CheckFile(name)
					if err == nil {
						log.Printf("Error Case %s: FAILED", name)
					}
					So(err, ShouldNotBeNil)
					return
				}
				log.Printf("Unrecognized file type in test suite: %s", name)
			})
		}
	})
}
