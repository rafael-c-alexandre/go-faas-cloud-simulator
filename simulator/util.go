package main

import (
	"encoding/gob"
	"log"
	"math/rand"
)

// RandStringBytes generates a random string of length n
func RandStringBytes(n int) string {

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	//rand.Seed(time.Now().Unix())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RemoveFromList[I any](s []I, i int) []I {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// InterfaceEncode encodes and saves the value to the encoder
func InterfaceEncode[I any](enc *gob.Encoder, i I) {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.
	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type.  If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&i)
	if err != nil {
		log.Fatalf("Caught error:%s", err)
	}
}

// InterfaceDecode decodes the value of the interface and returns
func InterfaceDecode[I any](dec *gob.Decoder, i I) I {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	err := dec.Decode(&i)
	if err != nil {
		log.Fatalf("Caught error:%s", err)
	}
	return i
}

func MergeMaps[I any](m1 map[string]I, m2 map[string]I) map[string]I {
	for k, v := range m2 {
		m1[k] = v
	}

	return m1
}
