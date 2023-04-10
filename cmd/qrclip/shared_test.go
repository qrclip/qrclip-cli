package main

import (
	"strings"
	"testing"
)

func TestEncodeBase64UrlNoPad(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8"},
		{[]byte("world"), "d29ybGQ"},
	}

	for _, tc := range testCases {
		actual := EncodeBase64UrlNoPad(tc.input)
		if actual != tc.expected {
			t.Errorf("EncodeBase64UrlNoPad(%s) = %s; expected %s", tc.input, actual, tc.expected)
		}
	}
}

func TestDecodeBase64Url(t *testing.T) {
	testCases := []struct {
		input          string
		expected       []byte
		expectError    bool
		errorSubstring string
	}{
		{"aGVsbG8", []byte("hello"), false, ""},
		{"d29ybGQ", []byte("world"), false, ""},
		{"invalid_base64&", nil, true, "illegal base64 data at input byte 14"},
	}

	for _, tc := range testCases {
		actual, err := DecodeBase64Url(tc.input)
		if err != nil && !tc.expectError {
			t.Errorf("DecodeBase64Url(%s) returned an unexpected error: %v", tc.input, err)
		}

		if err == nil && tc.expectError {
			t.Errorf("DecodeBase64Url(%s) did not return the expected error", tc.input)
		}

		if err != nil && tc.expectError && !strings.Contains(err.Error(), tc.errorSubstring) {
			t.Errorf("DecodeBase64Url(%s) returned an error with an unexpected message: %v", tc.input, err)
		}

		if string(actual) != string(tc.expected) {
			t.Errorf("DecodeBase64Url(%s) = %s; expected %s", tc.input, actual, tc.expected)
		}
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	tKeyBase64 := GenerateEncryptionKey()
	tKey, tErr := DecodeBase64Url(tKeyBase64)
	if tErr != nil {
		t.Errorf("Invalid base64 key")
	}
	if len(tKey) != 32 {
		t.Errorf("Wrong key size")
	}
}
