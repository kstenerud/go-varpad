// Varpad is a padding encoding scheme that embeds the padding length inside the
// padding itself.
package varpad

// Copyright 2019 Karl Stenerud
// All rights reserved.
// Distributed under MIT license.

import (
	"fmt"

	"github.com/kstenerud/go-vlq"
)

type Varpad vlq.Rvlq

// Calculate the appropriate padding for a given length and modulus.
// A modulus less than 2 will always result in a padding of 0.
func Padding(length int, paddingModulus int) Varpad {
	if paddingModulus < 2 {
		return 0
	}
	modulus := int(paddingModulus)
	remainder := length % modulus
	return Varpad(modulus - remainder)
}

// Encode padding to a buffer. Returns an error if the buffer isn't big enough.
//
// Note: A padding of zero will write nothing at all, which means that there
// will be no length field to be read by the other side. Be sure this is what
// you want, or else guard against it.
func (this Varpad) EncodeTo(buffer []byte) error {
	if this < 1 {
		return nil
	}
	if len(buffer) < int(this) {
		return fmt.Errorf("Not enough bytes in buffer to store %v padding bytes (buffer size is %v)", this, len(buffer))
	}
	buffer = buffer[:int(this)]
	length := vlq.Rvlq(this)
	bytesWritten, err := length.EncodeTo(buffer)
	if err != nil {
		return err
	}
	_, err = length.EncodeReversedTo(buffer)
	if err != nil {
		return err
	}
	middle := buffer[bytesWritten-1]
	buffer = buffer[bytesWritten : len(buffer)-bytesWritten+1]
	for i, _ := range buffer {
		buffer[i] = middle
	}

	return nil
}

// Fill the selected buffer with padding.
func FillWithPadding(buffer []byte) {
	Varpad(len(buffer)).EncodeTo(buffer)
}

// Decode the padding length from a buffer. Returns true for isComplete once the
// length is fully decoded (i.e. it has encountered a byte with the high bit cleared).
// This allows for progressive decoding of the length across multiple buffers.
func (this *Varpad) DecodeFromBeginning(buffer []byte) (bytesDecoded int, isComplete bool) {
	length := vlq.Rvlq(*this)
	bytesDecoded, isComplete = length.DecodeFrom(buffer)
	*this = Varpad(length)
	return bytesDecoded, isComplete
}

// Decode the padding length from a buffer. Returns true for isComplete once the
// length is fully decoded (i.e. it has encountered a byte with the high bit cleared).
// This allows for progressive decoding of the length across multiple buffers.
func DecodeFromBeginning(buffer []byte) (value Varpad, bytesDecoded int, isComplete bool) {
	bytesDecoded, isComplete = value.DecodeFromBeginning(buffer)
	return value, bytesDecoded, isComplete
}

// Decode the padding length from the end of the buffer. Unlike DecodeFrom(),
// this decode function requires the entire padding sequence to be present in
// the buffer.
func (this *Varpad) DecodeFromEnd(buffer []byte) (bytesDecoded int, err error) {
	length := vlq.Rvlq(*this)
	bytesDecoded, err = length.DecodeReversedFrom(buffer)
	*this = Varpad(length)
	return bytesDecoded, err
}

// Decode the padding length from the end of the buffer. Unlike DecodeFrom(),
// this decode function requires the entire padding sequence to be present in
// the buffer.
func DecodeFromEnd(buffer []byte) (value Varpad, bytesDecoded int, err error) {
	bytesDecoded, err = value.DecodeFromEnd(buffer)
	return value, bytesDecoded, err
}
