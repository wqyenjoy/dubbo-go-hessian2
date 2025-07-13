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
	"sync"
)

var (
	// DecoderPool is a pool of Decoder objects
	DecoderPool = sync.Pool{
		New: func() interface{} {
			return NewCheapDecoderWithSkip(nil)
		},
	}

	// EncoderPool is a pool of Encoder objects
	EncoderPool = sync.Pool{
		New: func() interface{} {
			return NewEncoder()
		},
	}
)

// GetDecoder gets a Decoder from the pool and resets it with the given byte slice
func GetDecoder(b []byte) *Decoder {
	decoder := DecoderPool.Get().(*Decoder)
	return decoder.Reset(b)
}

// PutDecoder returns a Decoder to the pool
func PutDecoder(d *Decoder) {
	DecoderPool.Put(d)
}

// GetEncoder gets an Encoder from the pool
func GetEncoder() *Encoder {
	return EncoderPool.Get().(*Encoder)
}

// PutEncoder returns an Encoder to the pool
func PutEncoder(e *Encoder) {
	e.Buffer() // Reset the buffer
	EncoderPool.Put(e)
}
