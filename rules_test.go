package denada

import "log"
import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestSingularRule(t *testing.T) {
	Convey("Testing Singular Rule", t, func() {
		info, err := ParseRuleName("singleton")
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, ElementList{})
		So(info.Name, ShouldEqual, "singleton")
		So(info.Cardinality, ShouldEqual, Cardinality(Singleton))
	})
}

func TestOptionalRule(t *testing.T) {
	Convey("Testing Optional Rule", t, func() {
		info, err := ParseRuleName("optional?")
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, ElementList{})
		So(info.Name, ShouldEqual, "optional")
		So(info.Cardinality, ShouldEqual, Cardinality(Optional))
	})
}

func TestZoMRule(t *testing.T) {
	Convey("Testing Zero-Or-More Rule", t, func() {
		info, err := ParseRuleName("zom*")
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, ElementList{})
		So(info.Name, ShouldEqual, "zom")
		So(info.Cardinality, ShouldEqual, Cardinality(ZeroOrMore))
	})
}

func TestOoMRule(t *testing.T) {
	Convey("Testing One-Or-More Rule", t, func() {
		info, err := ParseRuleName("oom+")
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, ElementList{})
		So(info.Name, ShouldEqual, "oom")
		So(info.Cardinality, ShouldEqual, Cardinality(OneOrMore))
	})
}

func TestRecursiveRule(t *testing.T) {
	Convey("Testing Recursive Rule", t, func() {
		dummy := NewDeclaration("dummy", "dummy*")
		root := ElementList{dummy}
		context := RootContext(root)

		info, err := ParseRule("recur>$root", context)
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, root)
		So(info.Name, ShouldEqual, "recur")
		So(info.Cardinality, ShouldEqual, Cardinality(Singleton))
	})
}

func TestParentRule(t *testing.T) {
	Convey("Testing Parent Rule", t, func() {
		dummy := NewDeclaration("dummy", "dummy*")
		root := ElementList{dummy}
		context := RootContext(root)

		Convey("Check that root scope has no parent", func() {
			log.Printf("TestParentRule.context = %v", context)
			_, err := ParseRule("recur>..", context)
			So(err, ShouldNotBeNil)
		})

		child := ChildContext(ElementList{}, &context)

		Convey("Check that child scope has a parent", func() {
			info, err := ParseRule("recur>..", child)
			So(err, ShouldBeNil)
			So(info.Context.This, ShouldResemble, root)
			So(info.Name, ShouldEqual, "recur")
			So(info.Cardinality, ShouldEqual, Cardinality(Singleton))
		})

		Convey("Check that parent of child of root is root", func() {
			info, err := ParseRule("recur>..", child)
			So(err, ShouldBeNil)
			rinfo, err := ParseRule("recur>..", child)
			So(err, ShouldBeNil)
			So(info.Context, ShouldResemble, rinfo.Context)
		})
	})
}

func TestCurrentRule(t *testing.T) {
	Convey("Testing Current Rule", t, func() {
		dummy := NewDeclaration("dummy", "dummy*")
		root := ElementList{dummy}
		context := RootContext(root)

		info, err := ParseRule("recur>.", context)
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, root)
		So(info.Name, ShouldEqual, "recur")
		So(info.Cardinality, ShouldEqual, Cardinality(Singleton))

		info, err = ParseRule("recur", context)
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, root)
		So(info.Name, ShouldEqual, "recur")
		So(info.Cardinality, ShouldEqual, Cardinality(Singleton))
	})
}

func TestRecursiveComplexRule(t *testing.T) {
	Convey("Testing Complex Recursive Rule", t, func() {
		root := ElementList{new(Element)}
		context := RootContext(root)

		info, err := ParseRule("recur?>$root", context)
		So(err, ShouldBeNil)
		So(info.Context.This, ShouldResemble, root)
		So(info.Name, ShouldEqual, "recur")
		So(info.Cardinality, ShouldEqual, Cardinality(Optional))
	})
}
