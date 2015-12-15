package api

import (
	"fmt"
	"log"
	"net/http"
	"github.com/DanShu93/golang-playground/serialization"
)

var data = make(map[string]int)

type Hello struct{}

func (h Hello) ServeHTTP(
w http.ResponseWriter,
r *http.Request) {
	stock := r.URL.Query().Get("stock")
	data[stock]++
	fmt.Fprint(w, string(serialization.EncodeJsonMap(data)))
}

func Start() {
	var h Hello
	err := http.ListenAndServe("localhost:4000", h)
	if err != nil {
		log.Fatal(err)
	}
}