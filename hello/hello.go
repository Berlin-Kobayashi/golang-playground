package main

import (
	"fmt"

	"github.com/DanShu93/golang-playground/stringutil"
	"github.com/DanShu93/golang-playground/serialization"
	"github.com/DanShu93/golang-playground/fileutil"

//	"github.com/golang/example/stringutil"
	"github.com/DanShu93/golang-playground/api"
)

var m = map[string]int{"one":1, "two":2, "three":3}

func main() {
	fmt.Println(stringutil.Reverse("Hello GO /:)"))

	encodedMap := serialization.EncodeJsonMap(m)
	fmt.Println("encodedMap: ", string(encodedMap))

	decodedMap := serialization.DecodeJsonMap(encodedMap)
	fmt.Println("decodedMap: ", decodedMap)

	rawData := fileutil.ReadJsonMap()
	fileutil.ReadJsonMap()

	fmt.Println("data in file: ", string(rawData))

	decodedRawData := serialization.DecodeJsonMap(rawData)

	fmt.Println("decoded data in file: ", decodedRawData)

	fileutil.WriteJsonMap(encodedMap)

	fmt.Println("just writtten json into a file")

	fmt.Println("starting service now...")

	api.Start()
}