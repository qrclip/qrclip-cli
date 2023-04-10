package main

import (
	"testing"
)

func TestCalculatePasswordKey(t *testing.T) {
	tAccessKey, tPasswordKey, tErr := calculatePasswordKey("MySuperPassword!", "abcdefghijklmnop")
	if tAccessKey != "QA1-DT2MtDM76soobi6zfQ" {
		t.Errorf("Accesskey error")
	}

	if EncodeBase64UrlNoPad(tPasswordKey) != "2NC5-mqDGnjjg79BHg-L9hXPJ7yX66Qk7zZKLRO990VADX4NPYy0MzvqyihuLrN9rPdgxck4-5qrlTegJ0CX6g" {
		t.Errorf("Password key error")
	}

	if tErr != nil {
		t.Errorf("Error calculating password key:" + tErr.Error())
	}
}

func TestGenerateCombinedKey(t *testing.T) {
	_, tPasswordKey, _ := calculatePasswordKey("MySuperPassword!", "abcdefghijklmnop")

	if len(tPasswordKey) != 64 {
		t.Errorf("Password key length error")
	}

	tEncryptionKey := "Dlib7B7CSUzCL7Gz4NFVRSB8C4xVZbFzq24X7WWJqL0"

	tCombinedKey, tErr := generateCombinedKey(tEncryptionKey, tPasswordKey)

	if tErr != nil {
		t.Errorf("Error generating combined key:" + tErr.Error())
	}

	if tCombinedKey != "eOB8gQvrBNYIORReHotrVHelxb5gWgQl-K9uBqaVbCE" {
		t.Errorf("Combined key error")
	}

}
