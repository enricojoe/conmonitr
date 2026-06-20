package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	plain := strings.NewReplacer("+", "", "/", "", "=", "").Replace(base64.StdEncoding.EncodeToString(b))
	if len(plain) > 40 {
		plain = plain[:40]
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Printf("CONMONITR_PASSWORD=%s\n", plain)
	fmt.Printf("CONMONITR_PASSWORD_HASH=%s\n", string(hash))
}
