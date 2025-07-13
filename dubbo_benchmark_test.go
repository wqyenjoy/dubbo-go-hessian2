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

// DubboRequest simulates a typical Dubbo RPC request
type DubboRequest struct {
	MethodName     string
	ParameterTypes []string
	Arguments      []interface{}
	Attachments    map[string]interface{}
}

func init() {
	RegisterPOJO(&DubboRequest{})
}

// JavaClassName implements POJO interface
func (r *DubboRequest) JavaClassName() string {
	return "org.apache.dubbo.rpc.RpcInvocation"
}

// BenchmarkDubboRequestEncodeDecode tests the performance of encoding and decoding a typical Dubbo request
func BenchmarkDubboRequestEncodeDecode(b *testing.B) {
	req := &DubboRequest{
		MethodName:     "echo",
		ParameterTypes: []string{"java.lang.String"},
		Arguments:      []interface{}{"hello world"},
		Attachments: map[string]interface{}{
			"path":      "dubbo-x/dubbo.DubboService",
			"interface": "dubbo.DubboService",
			"version":   "1.0.0",
			"group":     "test",
			"timeout":   "3000",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := NewEncoder()
		err := e.Encode(req)
		if err != nil {
			b.Fatal(err)
		}
		bytes := e.Buffer()

		d := NewDecoder(bytes)
		_, err = d.Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDubboRequestEncodeDecodePooled tests the performance with object pooling
func BenchmarkDubboRequestEncodeDecodePooled(b *testing.B) {
	req := &DubboRequest{
		MethodName:     "echo",
		ParameterTypes: []string{"java.lang.String"},
		Arguments:      []interface{}{"hello world"},
		Attachments: map[string]interface{}{
			"path":      "dubbo-x/dubbo.DubboService",
			"interface": "dubbo.DubboService",
			"version":   "1.0.0",
			"group":     "test",
			"timeout":   "3000",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := GetEncoder()
		err := e.Encode(req)
		if err != nil {
			b.Fatal(err)
		}
		bytes := e.Buffer()
		PutEncoder(e)

		d := GetDecoder(bytes)
		_, err = d.Decode()
		if err != nil {
			b.Fatal(err)
		}
		PutDecoder(d)
	}
}

// BenchmarkDubboRequestEncodeDecodeParallel tests parallel encoding/decoding
func BenchmarkDubboRequestEncodeDecodeParallel(b *testing.B) {
	req := &DubboRequest{
		MethodName:     "echo",
		ParameterTypes: []string{"java.lang.String"},
		Arguments:      []interface{}{"hello world"},
		Attachments: map[string]interface{}{
			"path":      "dubbo-x/dubbo.DubboService",
			"interface": "dubbo.DubboService",
			"version":   "1.0.0",
			"group":     "test",
			"timeout":   "3000",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e := NewEncoder()
			err := e.Encode(req)
			if err != nil {
				b.Fatal(err)
			}
			bytes := e.Buffer()

			d := NewDecoder(bytes)
			_, err = d.Decode()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDubboRequestEncodeDecodePooledParallel tests parallel encoding/decoding with object pooling
func BenchmarkDubboRequestEncodeDecodePooledParallel(b *testing.B) {
	// 在并行测试中，我们需要为每个goroutine创建独立的请求对象
	// 以避免数据竞争和引用问题

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		// 为每个goroutine创建独立的请求对象
		req := &DubboRequest{
			MethodName:     "echo",
			ParameterTypes: []string{"java.lang.String"},
			Arguments:      []interface{}{"hello world"},
			Attachments: map[string]interface{}{
				"path":      "dubbo-x/dubbo.DubboService",
				"interface": "dubbo.DubboService",
				"version":   "1.0.0",
				"group":     "test",
				"timeout":   "3000",
			},
		}

		for pb.Next() {
			e := GetEncoder()
			err := e.Encode(req)
			if err != nil {
				b.Fatal(err)
			}
			bytes := e.Buffer()

			d := GetDecoder(bytes)
			_, err = d.Decode()
			if err != nil {
				b.Fatal(err)
			}
			PutDecoder(d)
			PutEncoder(e)
		}
	})
}

// BenchmarkLargeDubboRequestEncodeDecode tests with a larger payload
func BenchmarkLargeDubboRequestEncodeDecode(b *testing.B) {
	// Create a large attachment map to simulate a more complex request
	attachments := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		key := "key-" + string(rune(i))
		attachments[key] = "value-" + string(rune(i))
	}

	// Create a larger argument list
	args := make([]interface{}, 0, 20)
	for i := 0; i < 20; i++ {
		args = append(args, "argument-"+string(rune(i)))
	}

	req := &DubboRequest{
		MethodName:     "complexMethod",
		ParameterTypes: []string{"java.lang.String", "java.lang.Integer", "java.util.List"},
		Arguments:      args,
		Attachments:    attachments,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := NewEncoder()
		err := e.Encode(req)
		if err != nil {
			b.Fatal(err)
		}
		bytes := e.Buffer()

		d := NewDecoder(bytes)
		_, err = d.Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLargeDubboRequestEncodeDecodePooled tests with a larger payload using object pooling
func BenchmarkLargeDubboRequestEncodeDecodePooled(b *testing.B) {
	// Create a large attachment map to simulate a more complex request
	attachments := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		key := "key-" + string(rune(i))
		attachments[key] = "value-" + string(rune(i))
	}

	// Create a larger argument list
	args := make([]interface{}, 0, 20)
	for i := 0; i < 20; i++ {
		args = append(args, "argument-"+string(rune(i)))
	}

	req := &DubboRequest{
		MethodName:     "complexMethod",
		ParameterTypes: []string{"java.lang.String", "java.lang.Integer", "java.util.List"},
		Arguments:      args,
		Attachments:    attachments,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := GetEncoder()
		err := e.Encode(req)
		if err != nil {
			b.Fatal(err)
		}
		bytes := e.Buffer()
		PutEncoder(e)

		d := GetDecoder(bytes)
		_, err = d.Decode()
		if err != nil {
			b.Fatal(err)
		}
		PutDecoder(d)
	}
}
