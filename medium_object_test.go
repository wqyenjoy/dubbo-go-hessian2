/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hessian

import (
	"testing"
)

// MediumObject represents a typical RPC object around 32KB in size
type MediumObject struct {
	ID        int64
	Name      string
	Tags      []string
	Metadata  map[string]interface{}
	Items     []Item
	Config    map[string]interface{}
	Data      []byte
	Relations []Relation
}

// Item represents a sub-object in MediumObject
type Item struct {
	ID    int64
	Value string
	Data  map[string]interface{}
}

// Relation represents a relationship in MediumObject
type Relation struct {
	Type   string
	Target int64
	Props  map[string]string
}

func init() {
	RegisterPOJO(&MediumObject{})
	RegisterPOJO(&Item{})
	RegisterPOJO(&Relation{})
}

// JavaClassName implements POJO interface
func (m *MediumObject) JavaClassName() string {
	return "com.test.MediumObject"
}

// JavaClassName implements POJO interface
func (i *Item) JavaClassName() string {
	return "com.test.Item"
}

// JavaClassName implements POJO interface
func (r *Relation) JavaClassName() string {
	return "com.test.Relation"
}

// createMediumTestObject creates a test object around 32KB in size
func createMediumTestObject() *MediumObject {
	// Create a medium-sized object (approximately 32KB)
	obj := &MediumObject{
		ID:   12345,
		Name: "Medium Test Object",
		Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		Metadata: map[string]interface{}{
			"created":  "2023-01-01",
			"modified": "2023-06-30",
			"version":  "1.0.0",
			"active":   true,
			"counts":   map[string]int{"visits": 100, "downloads": 50, "shares": 25},
		},
		Items:  make([]Item, 0, 50),
		Config: make(map[string]interface{}),
		Data:   make([]byte, 8*1024), // 8KB of binary data
	}

	// Add 50 items (approximately 10KB)
	for i := 0; i < 50; i++ {
		item := Item{
			ID:    int64(i),
			Value: "Item value with some reasonable length to simulate real data",
			Data: map[string]interface{}{
				"prop1": i,
				"prop2": "value",
				"prop3": i%2 == 0,
			},
		}
		obj.Items = append(obj.Items, item)
	}

	// Add 30 relations (approximately 6KB)
	for i := 0; i < 30; i++ {
		relation := Relation{
			Type:   "related-to",
			Target: int64(1000 + i),
			Props: map[string]string{
				"strength": "high",
				"since":    "2023-01-01",
				"notes":    "This is a test relation with some data to make it more realistic",
			},
		}
		obj.Relations = append(obj.Relations, relation)
	}

	// Add configuration (approximately 8KB)
	for i := 0; i < 100; i++ {
		key := "config-key-" + string(rune(i))
		obj.Config[key] = "config-value-" + string(rune(i)) + "-with-some-additional-data-to-make-it-more-realistic"
	}

	return obj
}

// BenchmarkMediumObjectDecode tests decoding performance for medium objects
func BenchmarkMediumObjectDecode(b *testing.B) {
	// Create and encode a medium object
	obj := createMediumTestObject()
	encoder := NewEncoder()
	err := encoder.Encode(obj)
	if err != nil {
		b.Fatal(err)
	}
	bytes := encoder.Buffer()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decoder := NewDecoder(bytes)
		_, err := decoder.Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMediumObjectDecodePooled tests decoding performance with object pooling
func BenchmarkMediumObjectDecodePooled(b *testing.B) {
	// Create and encode a medium object
	obj := createMediumTestObject()
	encoder := GetEncoder()
	err := encoder.Encode(obj)
	if err != nil {
		b.Fatal(err)
	}
	bytes := encoder.Buffer()
	PutEncoder(encoder)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decoder := GetDecoder(bytes)
		_, err := decoder.Decode()
		if err != nil {
			b.Fatal(err)
		}
		PutDecoder(decoder)
	}
}

// BenchmarkMediumObjectEncode tests encoding performance for medium objects
func BenchmarkMediumObjectEncode(b *testing.B) {
	// Create a medium object
	obj := createMediumTestObject()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder := NewEncoder()
		err := encoder.Encode(obj)
		if err != nil {
			b.Fatal(err)
		}
		_ = encoder.Buffer()
	}
}

// BenchmarkMediumObjectEncodePooled tests encoding performance with object pooling
func BenchmarkMediumObjectEncodePooled(b *testing.B) {
	// Create a medium object
	obj := createMediumTestObject()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder := GetEncoder()
		err := encoder.Encode(obj)
		if err != nil {
			b.Fatal(err)
		}
		_ = encoder.Buffer()
		PutEncoder(encoder)
	}
}

// BenchmarkMediumObjectRoundTrip tests full round-trip performance
func BenchmarkMediumObjectRoundTrip(b *testing.B) {
	// Create a medium object
	obj := createMediumTestObject()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Encode
		encoder := NewEncoder()
		err := encoder.Encode(obj)
		if err != nil {
			b.Fatal(err)
		}
		bytes := encoder.Buffer()

		// Decode
		decoder := NewDecoder(bytes)
		_, err = decoder.Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMediumObjectRoundTripPooled tests full round-trip performance with object pooling
func BenchmarkMediumObjectRoundTripPooled(b *testing.B) {
	// Create a medium object
	obj := createMediumTestObject()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Encode
		encoder := GetEncoder()
		err := encoder.Encode(obj)
		if err != nil {
			b.Fatal(err)
		}
		bytes := encoder.Buffer()

		// Decode
		decoder := GetDecoder(bytes)
		_, err = decoder.Decode()
		if err != nil {
			b.Fatal(err)
		}
		PutDecoder(decoder)
		PutEncoder(encoder)
	}
}
