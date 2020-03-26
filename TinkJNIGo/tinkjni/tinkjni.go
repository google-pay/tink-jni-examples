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

func (d Decryptor) Decrypt(encryptedMessage string) string {
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

	// build payment method token recipient
	builder, err := env.NewObject("com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder")
	if err != nil {
		log.Fatal(err)
	}
	builder.CallMethod(env, "addSenderVerifyingKey", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(env, d.GoogleSigningKey))
	builder.CallMethod(env, "recipientId", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(env, d.RecipientId))
	builder.CallMethod(env, "protocolVersion", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(env, d.ProtocolVersion))
	builder.CallMethod(env, "addRecipientPrivateKey", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder", fromGoStr(env, d.RecipientPrivateKey))
	recipient, err := builder.CallMethod(env, "build", "com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient")

	// decrypt and return as output
	clearText, err := recipient.(*jnigi.ObjectRef).CallMethod(env, "unseal", "java/lang/String", fromGoStr(env, encryptedMessage))
	if err != nil {
		log.Fatal(err)
	}
	var clearTextStr = toGoStr(env, clearText.(*jnigi.ObjectRef))
	return clearTextStr
}

func fromGoStr(env *jnigi.Env, str string) *jnigi.ObjectRef {
	jstr, err := env.NewObject("java/lang/String", []byte(str))
	if err != nil {
		log.Fatal(err)
	}
	return jstr
}

func toGoStr(env *jnigi.Env, o *jnigi.ObjectRef) string {
	v, err := o.CallMethod(env, "getBytes", jnigi.Byte|jnigi.Array)
	if err != nil {
		log.Fatal(err)
	}
	return string(v.([]byte))
}
