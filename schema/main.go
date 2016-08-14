package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gorilla/schema"
)

func main() {
	basicModelBinding()
}

//the decoder will cache schema info of the struct
var decoder *schema.Decoder

func basicModelBinding() {
	form := map[string][]string{
		"firstName":         []string{"Tony"},
		"lastName":          []string{"Kung"},
		"updatedAt":         []string{"2016-08-11"},
		"addresses.0.city":  []string{"Hong Kong"},
		"addresses.0.state": []string{"N/A"},
		"addresses.1.city":  []string{"Hong Kong 2"},
		"addresses.1.state": []string{"N/A 2"},
	}

	var p Person

	if decoder == nil {
		decoder = schema.NewDecoder()
		decoder.RegisterConverter(time.Now(), func(val string) reflect.Value {
			result := reflect.Value{}
			if t, err := time.Parse("2006-01-02", val); err == nil {
				result = reflect.ValueOf(t)
			}

			return result
		})
	}

	decoder.Decode(&p, form)
	fmt.Println(p)
}

type Person struct {
	FirstName string    `schema:"firstName"`
	LastName  string    `schema:"-"` //skip this field
	Addresses []Address `schema:"addresses"`
	UpdatedAt time.Time `schema:"updatedAt"`
}

type Address struct {
	City  string `schema:"city"`
	State string `schema:"state"`
}
