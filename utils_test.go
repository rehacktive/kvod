package kvod

import (
	"testing"
)

func TestSerializeDeserialize(t *testing.T) {
	s := "this is just a string to test serialization/deserialization"
	s1, _ := serialize(s)

	var s2 string
	deserialize(s1, &s2)

	if s != string(s2) {
		t.Errorf("expected %x found %x", s, s2)
	}
}

func BenchmarkSerialize(b *testing.B) {
	s := "this is just a string to test serialization/deserialization"
	for n := 0; n < b.N; n++ {
		serialize(s)
	}
}

func BenchmarkDeserialize(b *testing.B) {
	s := "this is just a string to test serialization/deserialization"
	s1, _ := serialize(s)
	var s2 string

	for n := 0; n < b.N; n++ {
		deserialize(s1, &s2)
	}
}
