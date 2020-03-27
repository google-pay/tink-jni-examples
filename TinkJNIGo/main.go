// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	tinkjni.InitJVM(decryptor) // NOTE: load JVM once only after application has started instead of loading JVM for each and every decryption

	// 3. decrypt message
	var sampleToken = "a sample token can be generated using https://developers.google.com/pay/api/processors/guides/test-and-validation/token-generator"
	var clearText = decryptor.Decrypt(string(sampleToken))
	log.Println("Clear text after decryption: " + clearText)
}
