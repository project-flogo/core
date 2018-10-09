package support

//Copyright © 2014, Roger Peppe
//All rights reserved.
//
//Redistribution and use in source and binary forms, with or without modification,
//are permitted provided that the following conditions are met:
//
//* Redistributions of source code must retain the above copyright notice,
//this list of conditions and the following disclaimer.
//* Redistributions in binary form must reproduce the above copyright notice,
//this list of conditions and the following disclaimer in the documentation
//and/or other materials provided with the distribution.
//* Neither the name of this project nor the names of its contributors
//may be used to endorse or promote products derived from this software
//without specific prior written permission.
//
//THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
//TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
//PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
//LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
//NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
//SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package fastuuid provides fast UUID generation of 192 bit
// universally unique identifiers. It does not provide
// formatting or parsing of the identifiers (it is assumed
// that a simple hexadecimal or base64 representation
// is sufficient, for which adequate functionality exists elsewhere).
//
// Note that the generated UUIDs are not unguessable - each
// UUID generated from a Generator is adjacent to the
// previously generated UUID.
//
// It ignores RFC 4122.
//
// fm: Modified to return 128 bit UUID as string

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync/atomic"
)

// Generator represents a UUID generator that
// generates UUIDs in sequence from a random starting
// point.
type Generator struct {
	seed    [24]byte
	counter uint64
}

// NewGenerator returns a new Generator.
// It can fail if the crypto/rand read fails.
func NewGenerator() (*Generator, error) {
	var g Generator
	_, err := rand.Read(g.seed[:])
	if err != nil {
		return nil, errors.New("cannot generate random seed: " + err.Error())
	}
	return &g, nil
}

// Next returns the next UUID from the generator.
// Only the first 8 bytes can differ from the previous
// UUID, so taking a slice of the first 16 bytes
// is sufficient to provide a somewhat less secure 128 bit UUID.
//
// It is OK to call this method concurrently.
func (g *Generator) Next() [24]byte {
	x := atomic.AddUint64(&g.counter, 1)
	var counterBytes [8]byte
	binary.LittleEndian.PutUint64(counterBytes[:], x)

	uuid := g.seed
	for i, b := range counterBytes {
		uuid[i] ^= b
	}
	return uuid
}

// NextAsString returns the next UUID from the generator as a string.
func (g *Generator) NextAsString() string {
	x := atomic.AddUint64(&g.counter, 1)
	var counterBytes [8]byte
	binary.LittleEndian.PutUint64(counterBytes[:], x)

	uuid := g.seed
	for i, b := range counterBytes {
		uuid[i] ^= b
	}

	buf := make([]byte, 32)
	hex.Encode(buf[0:], uuid[0:16]) //truncate to 128bit UUID

	return string(buf)
}
