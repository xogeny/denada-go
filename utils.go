package denada

import "github.com/bitly/go-simplejson"

func getString(v *simplejson.Json) (string, bool) {
	str, err := v.String()
	return str, err == nil
}

func getStringArray(v *simplejson.Json) ([]string, bool) {
	str, err := v.StringArray()
	return str, err == nil
}

/*
   A special utility function to quickly get the value (if it exists) of constructs
   like this:

   fmu SomeName {
     parameter := "Contents";
   }
*/
func GetStringParameter(app *Element, key string) (string, bool) {
	e := app.Contents.Declarations().FirstNamed(key)
	if e == nil {
		return "", false
	}
	return getString(e.Value)
}

func GetStringValue(app *Element) (string, bool) {
	val := app.Value
	if val == nil {
		return "", false
	}
	return getString(val)
}

func GetStringModification(app *Element, key string) (string, bool) {
	val, exists := app.Modifications[key]
	if !exists {
		return "", false
	}
	return getString(val)
}

func GetStringArrayParameter(app *Element, key string) ([]string, bool) {
	e := app.Contents.Declarations().FirstNamed(key)
	if e == nil {
		return []string{}, false
	}
	return getStringArray(e.Value)
}

func GetStringArrayValue(app *Element) ([]string, bool) {
	val := app.Value
	if val == nil {
		return []string{}, false
	}
	return getStringArray(val)
}

func GetStringArrayModification(app *Element, key string) ([]string, bool) {
	val, exists := app.Modifications[key]
	if !exists {
		return []string{}, false
	}
	return getStringArray(val)
}
