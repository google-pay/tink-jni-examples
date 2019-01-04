using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using java.io;
using java.lang;
using java.util;
using net.sf.jni4net;
using net.sf.jni4net.adaptors;
using System.Configuration;
using net.sf.jni4net.jni;

namespace TinkDotNetSample
{
    /// <summary>
    /// A facade for payment method token recipient based on Tink library https://github.com/google/tink. 
    /// It uses Java Native Interface to invoke Tink Java library. Please ensure Java 8 or later is installed on your server.
    /// </summary>
    class TinkDecryptor
    {
        private JNIEnv _jniEnv = null;
        private string _trustedSigningKeysJson = null;

        /// <summary>
        /// Create TinkDecryptor by creating JNI .Net bridge and preload trusted signing keys from Google servers
        /// </summary>
        public TinkDecryptor()
        {
            // create bridge using jni4net.j.jar in the same folder as jni4net.n.dll
            var bridgeSetup = new BridgeSetup();
            bridgeSetup.AddAllJarsClassPath("./tink_jars/"); // load libs
            _jniEnv = Bridge.CreateJVM(bridgeSetup); // create jvm

            // preload trusted signing keys from Google servers, cache into memory before performing any transactions
            Class googlePaymentsPublicKeysManager = _jniEnv.FindClass("com/google/crypto/tink/apps/paymentmethodtoken/GooglePaymentsPublicKeysManager");
            var env = ConfigurationManager.AppSettings["Trusted_Signing_Keys_Env"]; // environment of trusted signing keys from configs
            if (env != "INSTANCE_TEST" && env != "INSTANCE_PRODUCTION")
                throw new ConfigurationErrorsException("Setting Trusted_Signing_Keys_Env must be either INSTANCE_TEST or INSTANCE_PRODUCTION.");

            // initialize public key manager and load signing keys
            var publicKeyManager = googlePaymentsPublicKeysManager.GetFieldValue<java.lang.Object>(env, "Lcom/google/crypto/tink/apps/paymentmethodtoken/GooglePaymentsPublicKeysManager;");
            publicKeyManager.Invoke("refreshInBackground", "()V");
            _trustedSigningKeysJson = publicKeyManager.Invoke<java.lang.String>("getTrustedSigningKeysJson", "()Ljava/lang/String;");
        }

        /// <summary>
        /// Decrypt the given cipher text by performing the necessary signature verification and * decryption (if required) steps based on the protocolVersion
        /// </summary>
        /// <param name="cipherText">cipher text</param>
        /// <returns>plain text</returns>
        public string Decrypt(string cipherText)
        { 
            // load configs
            var gatewayName = ConfigurationManager.AppSettings["Gateway_Name"];
            if (string.IsNullOrWhiteSpace(gatewayName))
                throw new ConfigurationErrorsException("Setting Gateway_name must be set.");
            var privateKey = ConfigurationManager.AppSettings["Base64_PKCS8_Private_key"];
            if (string.IsNullOrWhiteSpace(privateKey))
                throw new ConfigurationErrorsException("Setting Base64_PKCS8_Private_key must be set.");

            // build payment method token recipient
            Class paymentMethodTokenRecipientBuilder = _jniEnv.FindClass("com/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder");
            var recipientBuilder = paymentMethodTokenRecipientBuilder.newInstance();
            recipientBuilder.Invoke<java.lang.Object>("senderVerifyingKeys",
                "(Ljava/lang/String;)Lcom/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder;", _trustedSigningKeysJson);
            recipientBuilder.Invoke<java.lang.Object>("recipientId", "(Ljava/lang/String;)Lcom/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder;",
                "gateway:" + gatewayName);
            recipientBuilder.Invoke<java.lang.Object>("addRecipientPrivateKey", "(Ljava/lang/String;)Lcom/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient$Builder;",
                privateKey);
            var recipient = recipientBuilder.Invoke<java.lang.Object>("build", "()Lcom/google/crypto/tink/apps/paymentmethodtoken/PaymentMethodTokenRecipient;");
            
            // decrypt message
            var plainText = recipient.Invoke<java.lang.String>("unseal", "(Ljava/lang/String;)Ljava/lang/String;", cipherText);
            return plainText;
        }
    }
}
