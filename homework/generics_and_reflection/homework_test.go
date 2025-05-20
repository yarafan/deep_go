package main

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	result := make([]string, 0)
	st := reflect.TypeOf(person)
	sv := reflect.ValueOf(person)

	for i := range st.NumField() {
		f := st.Field(i)
		v := sv.Field(i)

		tagStr, ok := f.Tag.Lookup("properties")
		if !ok {
			continue
		}

		parts := strings.Split(tagStr, ",")
		if slices.Contains(parts, "omitempty") && v.IsZero() {
			continue
		}

		valueAsString := ""
		switch f.Type.Kind() {
		case reflect.String:
			valueAsString = v.String()
		case reflect.Int:
			valueAsString = strconv.Itoa(int(v.Int()))
		case reflect.Bool:
			valueAsString = strconv.FormatBool(v.Bool())
		}

		result = append(result, fmt.Sprintf("%s=%s", parts[0], valueAsString))
	}

	return strings.Join(result, "\n")
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
