package main

import (
	"encoding/json"
	"log"
	"os"
	"tinkjni"
)

func main() {
	// 1. construct decryptor
	decryptor := tinkjni.Decryptor{}

	// 2. load config (file path is relative to the current working directory)
	file, _ := os.Open("tinkjni/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&decryptor); err != nil {
		log.Fatal("failed to open conf.json: ", err)
	}

	// 3. decrypt message
	var sampleToken = "a sample token can be generated using https://developers.google.com/pay/api/processors/guides/test-and-validation/token-generator"
	var clearText = decryptor.Decrypt(string(sampleToken))
	log.Println("Clear text after decryption: " + clearText)
}
