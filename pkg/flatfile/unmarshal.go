package flatfile

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Unmarshaler interface {
	UnmarshalText([]byte) error
}

/*Unmarshal will read data and convert it into a struct based on a schema/map defined by struct tags

Struct tags are in the form `flatfile:"start,end"`. start and end should be integers > 0. End is optional
and if not specified, it will be set to the remaining fields.

*/
func Unmarshal(data []byte, v interface{}) error {
	//init parserTag for later use
	parserTag := &flatfileTag{}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("flatfile.Unmarshal: Unmarshal not complete. %s is not a pointer", reflect.TypeOf(v))
	}
	//Get underlying type
	vType := reflect.TypeOf(v).Elem()

	//Only process if kind is Struct
	if vType.Kind() != reflect.Struct {
		return nil
	}
	//Dereference pointer to struct
	vStruct := reflect.ValueOf(v).Elem()
	maxField := vStruct.NumField()
	//Loop through struct fields/properties
	for i := 0; i < maxField; i++ {
		//Get underlying type of field
		fieldType := vStruct.Field(i).Type()
		fieldTag, tagFlag := vType.Field(i).Tag.Lookup("flatfile")
		if !tagFlag {
			continue
		}

		if tagParseErr := parseFlatfileTag(fieldTag, parserTag); tagParseErr != nil {
			return fmt.Errorf("flatfile.Unmarshal: Failed to parse field tag %s: %w", fieldTag, tagParseErr)
		}
		//extract byte slice from byte data
		splitData := strings.Split(string(data), " ")
		lowerBound := parserTag.start - 1
		upperBound := len(splitData)
		if parserTag.end != "*" {
			upperBound, _ = strconv.Atoi(parserTag.end)
		}
		fieldData := []byte(strings.Join(splitData[lowerBound:upperBound], " "))
		err := assignBasedOnKind(fieldType.Kind(), vStruct.Field(i), fieldData, parserTag)
		if err != nil {
			return fmt.Errorf("flatfile.Unmarshal: Failed to unmarshal: %w", err)
		}
	}
	return nil
}
