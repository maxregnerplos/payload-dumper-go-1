package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Payload struct {
	Offset             int32
	OriginalSize       int32
	CompressedSize     int32
	CompressionType    int32
	EncryptionType     int32
	Reserved           int32
	EncryptionKeyValue int32
	Padding            [8]byte
	Data               []byte
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: payload-dumper-go <payload file>")
	}

	payloadFile := os.Args[1]
	payloadBytes, err := ioutil.ReadFile(payloadFile)
	if err != nil {
		log.Fatalf("Failed to read payload file %s: %v", payloadFile, err)
	}

	fmt.Println("Payload file:", payloadFile)

	headerSize := binary.Size(Payload{})
	if len(payloadBytes) < headerSize {
		log.Fatalf("Invalid payload file size: %d bytes", len(payloadBytes))
	}

	buf := bytes.NewBuffer(payloadBytes)
	var payload Payload
	err = binary.Read(buf, binary.LittleEndian, &payload)
	if err != nil {
		log.Fatalf("Failed to read payload header: %v", err)
	}

	fmt.Printf("Offset: %08X\n", payload.Offset)
	fmt.Printf("Original Size: %08X\n", payload.OriginalSize)
	fmt.Printf("Compressed Size: %08X\n", payload.CompressedSize)
	fmt.Printf("Compression Type: %08X\n", payload.CompressionType)
	fmt.Printf("Encryption Type: %08X\n", payload.EncryptionType)
	fmt.Printf("Reserved: %08X\n", payload.Reserved)
	fmt.Printf("Encryption Key Value: %08X\n", payload.EncryptionKeyValue)

	data := payloadBytes[headerSize:]
	if len(data) < int(payload.CompressedSize) {
		log.Fatalf("Invalid payload file size: %d bytes", len(payloadBytes))
	}

	fmt.Printf("Payload Data:\n% X\n", data)
}
