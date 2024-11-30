package mpi

import (
	"bytes"
	"encoding/gob"
	"log"
)

func init() {
	gob.Register([]int{})
	gob.Register(int(0))
}

// Serialize serializes data into bytes
func Serialize(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		log.Fatalf("Serialization error: %v", err)
	}
	return buf.Bytes()
}

// Deserialize deserializes bytes into the provided interface
func Deserialize(data []byte, v interface{}) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(v); err != nil {
		log.Fatalf("Deserialization error: %v", err)
	}
}
