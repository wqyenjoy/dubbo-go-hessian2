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

// BenchmarkEncoderLargeMap tests encoding performance for large maps
func BenchmarkEncoderLargeMap(b *testing.B) {
	// Generate a large map for testing
	data := generateLargeMap(2, 5) // About 300KB

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder := NewEncoder()
		err := encoder.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = encoder.Buffer()
	}
}

// BenchmarkEncoderLargeMapPooled tests encoding performance with object pooling
func BenchmarkEncoderLargeMapPooled(b *testing.B) {
	// Generate a large map for testing
	data := generateLargeMap(2, 5) // About 300KB

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder := GetEncoder()
		err := encoder.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = encoder.Buffer()
		PutEncoder(encoder)
	}
}

// BenchmarkEncoderMediumMap tests encoding performance for medium maps
func BenchmarkEncoderMediumMap(b *testing.B) {
	// Generate a medium map for testing (32KB)
	data := generateLargeMap(1, 10) // About 30KB

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder := NewEncoder()
		err := encoder.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = encoder.Buffer()
	}
}

// BenchmarkEncoderMediumMapPooled tests encoding performance for medium maps with pooling
func BenchmarkEncoderMediumMapPooled(b *testing.B) {
	// Generate a medium map for testing (32KB)
	data := generateLargeMap(1, 10) // About 30KB

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoder := GetEncoder()
		err := encoder.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = encoder.Buffer()
		PutEncoder(encoder)
	}
}

// BenchmarkEncoderParallel tests parallel encoding performance
func BenchmarkEncoderParallel(b *testing.B) {
	// Generate test data
	data := generateLargeMap(2, 5) // About 300KB

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encoder := NewEncoder()
			err := encoder.Encode(data)
			if err != nil {
				b.Fatal(err)
			}
			_ = encoder.Buffer()
		}
	})
}

// BenchmarkEncoderParallelPooled tests parallel encoding performance with object pooling
func BenchmarkEncoderParallelPooled(b *testing.B) {
	// Generate test data
	data := generateLargeMap(2, 5) // About 300KB

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encoder := GetEncoder()
			err := encoder.Encode(data)
			if err != nil {
				b.Fatal(err)
			}
			_ = encoder.Buffer()
			PutEncoder(encoder)
		}
	})
}

// BenchmarkEncoderComplexObject tests encoding performance for complex objects
func BenchmarkEncoderComplexObject(b *testing.B) {
	// Create a complex object with nested structures
	obj := createComplexTestObject()

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

// BenchmarkEncoderComplexObjectPooled tests encoding performance for complex objects with pooling
func BenchmarkEncoderComplexObjectPooled(b *testing.B) {
	// Create a complex object with nested structures
	obj := createComplexTestObject()

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

// Helper function to create a complex test object
func createComplexTestObject() map[string]interface{} {
	// Create a complex object with nested structures
	obj := map[string]interface{}{
		"id":   12345,
		"name": "Complex Test Object",
		"tags": []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		"metadata": map[string]interface{}{
			"created":  "2023-01-01",
			"modified": "2023-06-30",
			"version":  "1.0.0",
			"active":   true,
			"counts":   map[string]int{"visits": 100, "downloads": 50, "shares": 25},
		},
		"items": []interface{}{
			map[string]interface{}{"id": 1, "value": "Item 1"},
			map[string]interface{}{"id": 2, "value": "Item 2"},
			map[string]interface{}{"id": 3, "value": "Item 3"},
			map[string]interface{}{"id": 4, "value": "Item 4"},
			map[string]interface{}{"id": 5, "value": "Item 5"},
		},
		"config": map[string]interface{}{
			"timeout":  30,
			"retries":  3,
			"debug":    false,
			"features": []string{"feature1", "feature2", "feature3"},
			"limits": map[string]interface{}{
				"maxConnections": 100,
				"maxRequests":    1000,
				"maxUsers":       50,
			},
		},
	}
	return obj
}
