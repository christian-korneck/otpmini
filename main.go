package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
		fmt.Println("OTPmini - a toy TOTP generator")
		fmt.Println("Usage: otpmini <secret>")
		fmt.Println("Example: otpmini JBSAY3DGEHPK2PXP")
		os.Exit(1)
	}

	secret := removeWhitespace(strings.ToUpper(os.Args[1]))

	otp, err := generateTOTP(secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating OTP: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OTP: %06d\n", otp)
}

func generateTOTP(secret string) (int, error) {
	// Decode base32 secret
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return 0, fmt.Errorf("invalid base32 secret: %v", err)
	}

	// get current time step (30 seconds)
	timeStep := time.Now().Unix() / 30

	// convert time step to bytes
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timeStep))

	// generate HMAC-SHA1
	h := hmac.New(sha1.New, key)
	h.Write(timeBytes)
	hash := h.Sum(nil)

	// truncate
	offset := hash[len(hash)-1] & 0x0F
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF

	// modulo -> 6-digit OTP
	otp := int(truncated % 1000000)

	return otp, nil
}

func removeWhitespace(s string) string {
	var b []rune
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			b = append(b, r)
		}
	}
	return string(b)
}
