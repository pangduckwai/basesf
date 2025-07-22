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

	fmt.Println("TestEncode", content)
	fmt.Println("TestEncode", content[:9])
	fmt.Println("TestEncode", content[9:])

	// Encode as base64.
	encoder := base64.StdEncoding //.WithPadding(base64.NoPadding)
	encoded0 := encoder.EncodeToString(content)
	encoded1 := encoder.EncodeToString(content[:9])
	encoded2 := encoder.EncodeToString(content[9:])

	// Print encoded data to console.
	fmt.Println("TestEncode 0: " + encoded0)
	fmt.Println("TestEncode 1: " + encoded1)
	fmt.Println("TestEncode 2: " + encoded2)
}

func TestCopySlice(t *testing.T) {
	src0 := []byte{128, 182, 109, 169, 39, 17, 65, 10, 93, 201, 88, 143, 79, 5}
	src1 := []byte{65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78}
	buf := make([]byte, 0, 14)

	fmt.Println("TestCopySlice", buf[:14])
	copy(buf[:14], src0[:14])
	fmt.Println("TestCopySlice", buf[:14])
	copy(buf[:14], src1[:14])
	fmt.Println("TestCopySlice", buf[:14])
}

func TestAppend(t *testing.T) {
	src0 := []byte{65, 66, 67, 68, 69, 70, 71, 72, 73}
	src1 := []byte{74, 75, 76, 77, 78, 79, 80, 81, 82}
	src2 := []byte{83, 84, 85, 86, 87, 88, 89, 90}
	src3 := []byte{91, 92, 93}
	buf1 := make([]byte, 0, 20)
	cnt1 := 9
	cnt2 := 8
	cnt3 := 6 // 8 - 8%3
	cnt4 := 3 // size of next read

	copy(buf1[:cnt1], src0[:cnt1])

	fmt.Println("TestAppend", buf1[:cnt1])
	buf1 = append(buf1[cnt1:cnt1], src1...)

	fmt.Println("TestAppend", buf1[:cnt1])
	buf1 = append(buf1[cnt1:cnt1], src2...)

	fmt.Println("TestAppend", buf1[:cnt3])
	buf1 = append(buf1[cnt3:cnt2], src3...)

	fmt.Println("TestAppend", buf1[:(cnt2-cnt3+cnt4)])
}
