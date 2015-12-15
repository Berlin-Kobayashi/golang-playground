package serialization

import (
	"encoding/json"
	"bytes"
)

func EncodeJsonMap(m map[string]int) []byte {
	b := new(bytes.Buffer)

	e := json.NewEncoder(b)

	err := e.Encode(m)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}

func DecodeJsonMap(b []byte) map[string]int{
	buffer := bytes.NewBuffer(b)
	var decodedMap map[string]int
	d := json.NewDecoder(buffer)

	err := d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	return decodedMap
}