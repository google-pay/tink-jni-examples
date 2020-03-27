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

package tinkjni

import (
	"github.com/timob/jnigi"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

// A facade for payment method token recipient based on Tink library https://github.com/google/tink.
// It uses timob/jnigi to connect between Go and Java with Java Native Interface.
// Please ensure Java 8 or later is installed on your server.
type Decryptor struct {
	JVMLibraryPath      string
	GoogleSigningKey    string
	ProtocolVersion     string
	RecipientId         string
	RecipientPrivateKey string
}

var PaymentMethodTokenRecipient *jnigi.ObjectRef // static variable
var JavaEnvironment *jnigi.Env                   // static variable

func InitJVM(d Decryptor) { // static method
	// load JVM library
	if err := jnigi.LoadJVMLib(d.JVMLibraryPath); err != nil {
		log.Fatal(err)
		log.Println(d.JVMLibraryPath)
	}

	// construct classpath for JVM
	_, curFilename, _, _ := runtime.Caller(1)
	wd := path.Dir(curFilename)
	java_libs := path.Join(wd, "tinkjni", "java_libs")
	var classpath string
	files, err := ioutil.ReadDir(java_libs)
	if err != nil {
		log.Fatal(err)
	}
	var classpathSeparator = ":" // Linux and Mac
	if string(os.PathSeparator) == "\\" {
		classpathSeparator = ";" // Windows
	}
	for _, f := range files {
		classpath += path.Join(java_libs, f.Name()) + classpathSeparator
	}

	// create JVM
	_, env, err := jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION, []string{"-Xcheck:jni", "-Xbootclasspath/a:" + classpath}))
	if err != nil {
		log.Fatal(err)
	}
	JavaEnvironment = env

	// build payment method token recipient
	builder, err := env.NewObject("com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder")
	if err != nil {
		log.Fatal(err)
	}
	builder.CallMethod(env, "addSenderVerifyingKey", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(d.GoogleSigningKey))
	builder.CallMethod(env, "recipientId", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(d.RecipientId))
	builder.CallMethod(env, "protocolVersion", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(d.ProtocolVersion))
	builder.CallMethod(env, "addRecipientPrivateKey", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(d.RecipientPrivateKey))
	recipient, err := builder.CallMethod(env, "build", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient")
	PaymentMethodTokenRecipient = recipient.(*jnigi.ObjectRef)
}

func (d Decryptor) Decrypt(encryptedMessage string) string {
	// decrypt and return as output
	clearText, err := PaymentMethodTokenRecipient.CallMethod(JavaEnvironment, "unseal", "java/lang/String", fromGoStr(encryptedMessage))
	if err != nil {
		log.Fatal(err)
	}
	var clearTextStr = toGoStr(clearText.(*jnigi.ObjectRef))
	return clearTextStr
}

func fromGoStr(str string) *jnigi.ObjectRef {
	jstr, err := JavaEnvironment.NewObject("java/lang/String", []byte(str))
	if err != nil {
		log.Fatal(err)
	}
	return jstr
}

func toGoStr(o *jnigi.ObjectRef) string {
	v, err := o.CallMethod(JavaEnvironment, "getBytes", jnigi.Byte|jnigi.Array)
	if err != nil {
		log.Fatal(err)
	}
	return string(v.([]byte))
}
