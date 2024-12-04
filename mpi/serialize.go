package mpi

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

func init() {
	// Register more common types
	gob.Register([]int{})
	gob.Register([]float64{})
	gob.Register(int(0))
	gob.Register(float64(0))
}

// Serialize serializes data into bytes with improved error handling
func Serialize(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Handle different types for registration
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice, reflect.Array:
		v := reflect.ValueOf(data)
		if v.Len() > 0 {
			// Register the element type
			gob.Register(v.Index(0).Interface())
		}
	}

	if err := enc.Encode(data); err != nil {
		// Instead of fatal log, return an error or panic with more context
		panic(fmt.Sprintf("Serialization error for type %T: %v", data, err))
	}
	return buf.Bytes()
}

// Deserialize deserializes bytes into the provided interface with improved error handling
func Deserialize(data []byte, v interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("cannot deserialize empty byte slice")
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	// Check if the target is a pointer
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer, got %T", v)
	}

	// Actual decoding with error return
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("deserialization error for type %T: %v", v, err)
	}
	return nil
}

// RegisterType allows manual registration of custom types
func RegisterType(v interface{}) {
	gob.Register(v)
}
