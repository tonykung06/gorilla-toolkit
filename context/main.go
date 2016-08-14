package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
)

const key = "MyKey"

func main() {
	// singleValue()
	// multipleValues()
	purgeContext()
}

func purgeContext() {
	request1 := &http.Request{}
	context.Set(request1, key, "foo")
	fmt.Println(context.Get(request1, key))
	time.Sleep(2 * time.Second)
	context.Purge(1)
	fmt.Println(context.Get(request1, key))
}

func multipleValues() {

	request1 := &http.Request{}
	context.Set(request1, key, "foo")
	context.Set(request1, "second", "bar")

	//returns a copy of the map
	for key, val := range context.GetAll(request1) {
		fmt.Printf("%s: %s\n", key, val)
	}
}

func singleValue() {
	request1 := &http.Request{}
	request2 := &http.Request{}
	context.Set(request1, key, "foo")
	context.Set(request2, key, "bar")

	if val, ok := context.GetOk(request1, key); ok {
		fmt.Println(val)
	}

	_, ok := context.GetOk(request1, "notfound")
	fmt.Println(ok)

	context.Delete(request1, key)

	fmt.Println(context.Get(request1, key))
	fmt.Println(context.Get(request2, key))

	context.Clear(request2)
	fmt.Println(context.Get(request2, key))
}
