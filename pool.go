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
	"bytes"
	"os"
	"strconv"
	"sync"
)

// EnablePool controls whether to use object pooling
// It can be enabled by setting the HESSIAN_POOL environment variable to 1
var EnablePool = envBool("HESSIAN_POOL", false)

// BufferSizeThreshold is the threshold for buffer size to determine which pool to use
// Buffers larger than this threshold will use LargeBufferPool
const BufferSizeThreshold = 64 << 10 // 64KB

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

	// BufferPool is a pool of medium byte slices
	BufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 4<<10) // 4KB
		},
	}

	// LargeBufferPool is a pool of large byte slices
	LargeBufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, BufferSizeThreshold) // 64KB
		},
	}
)

// envBool returns the boolean value of the environment variable
func envBool(name string, defaultVal bool) bool {
	v := os.Getenv(name)
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}

// GetDecoder gets a Decoder from the pool and resets it with the given byte slice
func GetDecoder(b []byte) *Decoder {
	if !EnablePool {
		return NewDecoder(b)
	}

	decoder := DecoderPool.Get().(*Decoder)

	// For large payloads, use a pooled buffer to reduce GC pressure
	if len(b) > 4096 {
		var buf []byte
		if len(b) > BufferSizeThreshold {
			buf = GetLargeBuffer()
		} else {
			buf = GetBuffer()
		}

		if cap(buf) >= len(b) {
			buf = buf[:len(b)]
			copy(buf, b)
			decoder.reader.Reset(bytes.NewReader(buf))
			decoder.Clean()
			// Store the buffer reference in the decoder for later cleanup
			decoder.pooledBuffer = buf
			return decoder
		}
		// If buffer is too small, put it back
		if len(b) > BufferSizeThreshold {
			PutLargeBuffer(buf)
		} else {
			PutBuffer(buf)
		}
	}

	return decoder.Reset(b)
}

// PutDecoder returns a Decoder to the pool
func PutDecoder(d *Decoder) {
	if !EnablePool {
		return
	}

	// Thoroughly clean the decoder before returning it to the pool
	d.Reset(nil)          // Reset the reader with empty buffer
	d.refHolders = nil    // Release any references
	d.refs = nil          // Release object references
	d.classInfoList = nil // Release class info
	d.typeRefs = nil      // Release type references

	// Return any pooled buffer back to the buffer pool
	if d.pooledBuffer != nil {
		if cap(d.pooledBuffer) > BufferSizeThreshold {
			PutLargeBuffer(d.pooledBuffer)
		} else {
			PutBuffer(d.pooledBuffer)
		}
		d.pooledBuffer = nil
	}

	DecoderPool.Put(d)
}

// GetEncoder gets an Encoder from the pool
func GetEncoder() *Encoder {
	if !EnablePool {
		return NewEncoder()
	}
	return EncoderPool.Get().(*Encoder)
}

// PutEncoder returns an Encoder to the pool
func PutEncoder(e *Encoder) {
	if !EnablePool {
		return
	}

	// Clean the encoder before returning it to the pool
	e.ReuseBufferClean() // Reset the buffer and clean references
	EncoderPool.Put(e)
}

// GetBuffer gets a medium byte buffer from the pool
func GetBuffer() []byte {
	if !EnablePool {
		return make([]byte, 0, 4<<10)
	}
	return BufferPool.Get().([]byte)[:0]
}

// PutBuffer returns a byte buffer to the pool
func PutBuffer(b []byte) {
	if !EnablePool {
		return
	}

	// Only return reasonably sized buffers to the pool
	if cap(b) >= 4<<10 && cap(b) <= BufferSizeThreshold {
		BufferPool.Put(b[:0])
	}
}

// GetLargeBuffer gets a large byte buffer from the pool
func GetLargeBuffer() []byte {
	if !EnablePool {
		return make([]byte, 0, BufferSizeThreshold)
	}
	return LargeBufferPool.Get().([]byte)[:0]
}

// PutLargeBuffer returns a large byte buffer to the pool
func PutLargeBuffer(b []byte) {
	if !EnablePool {
		return
	}

	// Only return reasonably sized large buffers to the pool
	if cap(b) >= BufferSizeThreshold && cap(b) <= BufferSizeThreshold*2 {
		LargeBufferPool.Put(b[:0])
	}
}

// TODO: In Go 1.22+, consider using arena.New() for more efficient memory allocation
