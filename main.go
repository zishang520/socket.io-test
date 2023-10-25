package main

import (
	"encoding/binary"
	"errors"
)

type State int

const (
	READ_HEADER State = iota
	READ_EXTENDED_LENGTH_16
	READ_EXTENDED_LENGTH_64
	READ_PAYLOAD
)

func totalLength(chunks [][]byte) int {
	total := 0
	for _, chunk := range chunks {
		total += len(chunk)
	}
	return total
}

func concatChunks(chunks [][]byte, size int) ([]byte, [][]byte) {
	if len(chunks[0]) == size {
		return chunks[0], chunks[1:]
	}
	buffer := make([]byte, size)
	j := 0
	for i := 0; i < size; i++ {
		buffer[i] = chunks[0][j]
		j++
		if j == len(chunks[0]) {
			chunks = chunks[1:]
			j = 0
		}
	}
	if len(chunks) > 0 && j < len(chunks[0]) {
		chunks[0] = chunks[0][j:]
	}
	return buffer, chunks
}

func createPacketDecoderStream(maxPayload int) TransformStream {
	var chunks [][]byte
	state := READ_HEADER
	expectedLength := -1
	isBinary := false

	transform := func(chunk []byte, controller TransformStreamController) {
		chunks = append(chunks, chunk)
		for {
			if state == READ_HEADER {
				if totalLength(chunks) < 1 {
					break
				}
				header, remainingChunks := concatChunks(chunks, 1)
				chunks = remainingChunks
				isBinary = (header[0] & 0x80) == 0x80
				expectedLength = int(header[0] & 0x7f)
				if expectedLength < 126 {
					state = READ_PAYLOAD
				} else if expectedLength == 126 {
					state = READ_EXTENDED_LENGTH_16
				} else {
					state = READ_EXTENDED_LENGTH_64
				}
			} else if state == READ_EXTENDED_LENGTH_16 {
				if totalLength(chunks) < 2 {
					break
				}
				header, remainingChunks := concatChunks(chunks, 2)
				chunks = remainingChunks
				expectedLength = int(binary.BigEndian.Uint16(header))
				state = READ_PAYLOAD
			} else if state == READ_EXTENDED_LENGTH_64 {
				if totalLength(chunks) < 8 {
					break
				}
				header, remainingChunks := concatChunks(chunks, 8)
				chunks = remainingChunks
				n := binary.BigEndian.Uint32(header)
				if n > (1<<53)-1 {
					controller.Enqueue(ERROR_PACKET)
					break
				}
				expectedLength = int(n<<32 + binary.BigEndian.Uint32(header[4:]))
				state = READ_PAYLOAD
			} else {
				if totalLength(chunks) < expectedLength {
					break
				}
				data, remainingChunks := concatChunks(chunks, expectedLength)
				chunks = remainingChunks
				var decodedData interface{}
				if isBinary {
					decodedData = data // implement your binary decoding logic here
				} else {
					decodedData = string(data) // implement your text decoding logic here
				}
				controller.Enqueue(decodedData)
				state = READ_HEADER
			}

			if expectedLength == 0 || expectedLength > maxPayload {
				controller.Enqueue(ERROR_PACKET)
				break
			}
		}
	}

	// Implement TransformStream and its methods here

	return transform // Return your TransformStream function here
}

// Implement the remaining parts of the code including TransformStream and its methods.
