Varpad
======

A go implementation of [Varpad padding](https://github.com/kstenerud/varpad/blob/master/varpad-specification.md). Varpad embeds the length of the padding inside the padding itself, removing the need for a separate field.


Library Usage
-------------

```golang
func readme_example_simple() {
	buffer := make([]byte, 10)
	varpad.FillWithPadding(buffer)
	fmt.Printf("Padding: %v", hex.Dump(buffer))
}
```

Result:

	Padding: 00000000  0a 0a 0a 0a 0a 0a 0a 0a  0a 0a                    |..........|


```golang
func readme_example_leading_padding() {
	message := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee}
	messageLength := len(message)
	paddingModulus := 4
	padding := varpad.Padding(messageLength, paddingModulus)
	envelope := make([]byte, messageLength+int(padding))
	padding.EncodeTo(envelope)
	copy(envelope[int(padding):], message)
	fmt.Printf("Envelope: %v", hex.Dump(envelope))

	decodedPadding, bytesDecoded, isComplete := varpad.DecodeFromBeginning(envelope)
	if !isComplete {
		// TODO: Normally this would mean that you need to fetch more bytes
	}
	fmt.Printf("Decoded padding amount %v (length was encoded into %v bytes)\n",
		decodedPadding, bytesDecoded)
}
```

Result:

	Envelope: 00000000  03 03 03 aa bb cc dd ee                           |........|
	Decoded padding amount 3 (length was encoded into 1 bytes)


```golang
func readme_example_trailing_padding() {
	message := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee}
	messageLength := len(message)
	paddingModulus := 4
	padding := varpad.Padding(messageLength, paddingModulus)
	envelope := make([]byte, messageLength+int(padding))
	copy(envelope, message)
	padding.EncodeTo(envelope[messageLength:])
	fmt.Printf("Envelope: %v", hex.Dump(envelope))

	decodedPadding, bytesDecoded, err := varpad.DecodeFromEnd(envelope)
	if err != nil {
		// TODO: Error handling
	}
	fmt.Printf("Decoded padding amount %v (length was encoded into %v bytes)\n",
		decodedPadding, bytesDecoded)
}
```

Result:

	Envelope: 00000000  aa bb cc dd ee 03 03 03                           |........|
	Decoded padding amount 3 (length was encoded into 1 bytes)