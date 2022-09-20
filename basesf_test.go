package main

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	content := []byte{128, 182, 109, 169, 39, 17, 65, 10, 93, 201, 88, 143, 79, 5}

	// // Write content to file for testing
	// err := os.WriteFile("test.bin", content, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(content)
	fmt.Println(content[:9])
	fmt.Println(content[9:])

	// Encode as base64.
	encoder := base64.StdEncoding //.WithPadding(base64.NoPadding)
	encoded0 := encoder.EncodeToString(content)
	encoded1 := encoder.EncodeToString(content[:9])
	encoded2 := encoder.EncodeToString(content[9:])

	// Print encoded data to console.
	fmt.Println("ENCODED 0: " + encoded0)
	fmt.Println("ENCODED 1: " + encoded1)
	fmt.Println("ENCODED 2: " + encoded2)
}

func TestCopySlice(t *testing.T) {
	buf0 := make([]byte, 0, 14)
	buf1 := []byte{128, 182, 109, 169, 39, 17, 65, 10, 93, 201, 88, 143, 79, 5}
	buf2 := []byte{65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78}

	fmt.Println(buf0[:14])
	copy(buf0[:14], buf1[:14])
	fmt.Println(buf0[:14])
	copy(buf0[:14], buf2[:14])
	fmt.Println(buf0[:14])
}
