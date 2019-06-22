package varpad

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func generate_bytes(value byte, count int) []byte {
	buffer := make([]byte, 0, count)
	for i := 0; i < count; i++ {
		buffer = append(buffer, value)
	}
	return buffer
}

func assertEncoded(t *testing.T, value int, expectedEncoded []byte) {
	pad := Varpad(value)
	byteCount := int(pad)
	actualEncoded := make([]byte, byteCount)
	err := pad.EncodeTo(actualEncoded)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(expectedEncoded, actualEncoded) {
		t.Errorf("Expected %v but got %v", expectedEncoded, actualEncoded)
	}
	decoded, _, isComplete := DecodeFromBeginning(actualEncoded)
	if !isComplete {
		t.Errorf("Expected decoding to be complete")
		return
	}
	if decoded != pad {
		t.Errorf("Expected decoded value %v but got %v", pad, decoded)
		return
	}
}

func assertEncodingFails(t *testing.T, value int) {
	pad := Varpad(value)
	byteCount := int(pad)
	actualEncoded := make([]byte, byteCount)
	err := pad.EncodeTo(actualEncoded)
	if err == nil {
		t.Error("Expected an error but none occurred")
		return
	}
}

func Test0(t *testing.T) {
	expected := []byte{}
	buffer := make([]byte, 0)
	FillWithPadding(buffer)
	if !bytes.Equal(expected, buffer) {
		t.Errorf("Expected %v but got %v", expected, buffer)
	}
}

func Test1(t *testing.T) {
	assertEncoded(t, 1, generate_bytes(1, 1))
}

func Test2(t *testing.T) {
	assertEncoded(t, 2, generate_bytes(2, 2))
}

func Test3(t *testing.T) {
	assertEncoded(t, 3, generate_bytes(3, 3))
}

func Test4(t *testing.T) {
	assertEncoded(t, 4, generate_bytes(4, 4))
}

func Test7F(t *testing.T) {
	assertEncoded(t, 0x7f, generate_bytes(0x7f, 0x7f))
}

func generatePadding(header []byte, middle byte, middleCount int, footer []byte) []byte {
	padding := append(header, generate_bytes(middle, middleCount)...)
	return append(padding, footer...)
}

func Test80(t *testing.T) {
	assertEncoded(t, 128, generatePadding([]byte{0x81}, 0x00, 0x7e, []byte{0x81}))
}

func Test81(t *testing.T) {
	assertEncoded(t, 129, generatePadding([]byte{0x81}, 0x01, 0x7f, []byte{0x81}))
}

func TestFF(t *testing.T) {
	assertEncoded(t, 255, generatePadding([]byte{0x81}, 0x7f, 0xfd, []byte{0x81}))
}

func Test200000(t *testing.T) {
	assertEncoded(t, 0x200000, generatePadding([]byte{0x81, 0x80, 0x80}, 0x00, 0x200000-6, []byte{0x80, 0x80, 0x81}))
}

func TestPadding(t *testing.T) {
	message := "12345"
	messageLength := len(message)
	paddingModulus := 4
	padding := Padding(messageLength, paddingModulus)
	envelope := make([]byte, messageLength+int(padding))
	copy(envelope, message)
	padding.EncodeTo(envelope[messageLength:])
	expected := []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x03, 0x03, 0x03}
	if !bytes.Equal(expected, envelope) {
		t.Errorf("Expected %v but got %v", expected, envelope)
	}
}

func TestFillWithPadding(t *testing.T) {
	expected := []byte{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
	buffer := make([]byte, 10)
	FillWithPadding(buffer)
	if !bytes.Equal(expected, buffer) {
		t.Errorf("Expected %v but got %v", expected, buffer)
	}
}

func readme_example_simple() {
	buffer := make([]byte, 10)
	FillWithPadding(buffer)
	fmt.Printf("Padding: %v", hex.Dump(buffer))
}

func readme_example_leading_padding() {
	message := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee}
	messageLength := len(message)
	paddingModulus := 4
	padding := Padding(messageLength, paddingModulus)
	envelope := make([]byte, messageLength+int(padding))
	padding.EncodeTo(envelope)
	copy(envelope[int(padding):], message)
	fmt.Printf("Envelope: %v", hex.Dump(envelope))

	decodedPadding, bytesDecoded, isComplete := DecodeFromBeginning(envelope)
	if !isComplete {
		// TODO: Normally this would mean that you need to fetch more bytes
	}
	fmt.Printf("Decoded padding amount %v (length was encoded into %v bytes)\n",
		decodedPadding, bytesDecoded)
}

func readme_example_trailing_padding() {
	message := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee}
	messageLength := len(message)
	paddingModulus := 4
	padding := Padding(messageLength, paddingModulus)
	envelope := make([]byte, messageLength+int(padding))
	copy(envelope, message)
	padding.EncodeTo(envelope[messageLength:])
	fmt.Printf("Envelope: %v", hex.Dump(envelope))

	decodedPadding, bytesDecoded, err := DecodeFromEnd(envelope)
	if err != nil {
		// TODO: Error handling
	}
	fmt.Printf("Decoded padding amount %v (length was encoded into %v bytes)\n",
		decodedPadding, bytesDecoded)
}

func TestExamples(t *testing.T) {
	readme_example_simple()
	readme_example_leading_padding()
	readme_example_trailing_padding()
}
