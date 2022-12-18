package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
)

func main() {
	// encrypt
	password := "nokia123"
	plaintext := []byte(password)

	key := make([]byte, 32)
	rand.Read(key)
	fmt.Println(key, len(key))

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	bs := c.BlockSize()
	iv := make([]byte, bs)
	rand.Read(iv)
	fmt.Println(iv, len(iv))

	ciphertext := make([]byte, len(plaintext))
	cfb := cipher.NewCFBEncrypter(c, iv)
	cfb.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("%s => %x\n", plaintext, ciphertext)

	// decrypt
	text := make([]byte, len(ciphertext))
	cfbdec := cipher.NewCFBEncrypter(c, iv)
	cfbdec.XORKeyStream(text, ciphertext)
	fmt.Printf("%x => %s\n", ciphertext, text)
}
