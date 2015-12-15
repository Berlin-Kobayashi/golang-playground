package fileutil

import (
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadJsonMap() []byte {
	dat, err := ioutil.ReadFile("C:/Users/Daniel Schulz/go-workspace/src/github.com/DanShu93/golang-playground/resources/example2.json")
	check(err)
	return dat
}

func WriteJsonMap(b []byte)  {
	err := ioutil.WriteFile("C:/Users/Daniel Schulz/go-workspace/src/github.com/DanShu93/golang-playground/resources/example2.json", b, 0644)
	check(err)
}