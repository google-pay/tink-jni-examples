# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os.path

class Decryptor:
    @staticmethod
    def init(config):
        # Validate the JARs folder
        # Path is relative to the current working directory
        if not os.path.isdir('./java_libs/'):
            raise Exception("Java library path could not be located.")

        # Construct classpath for JVM
        import jnius_config
        jnius_config.add_options('-Xrs', '-Xmx4m')  # JVM options in comma separated
        jnius_config.set_classpath('./java_libs/*') # Path for the Java library
        import jnius

        # Build Payment Method Token Recipient
        Builder = jnius.autoclass('com.google.crypto.tink.apps.paymentmethodtoken.PaymentMethodTokenRecipient$Builder')
        builder = Builder()
        builder.addSenderVerifyingKey(config['GoogleSigningKey'])
        builder.recipientId(config['RecipientId'])
        builder.protocolVersion(config['ProtocolVersion'])
        builder.addRecipientPrivateKey(config['RecipientPrivateKey'])
        Decryptor.paymentMethodTokenRecipient = builder.build()

    @staticmethod
    def decrypt(encryptedMessage):
        # Validate the PaymentMethodTokenRecipient has been built
        if not hasattr(Decryptor, "paymentMethodTokenRecipient"):
            raise Exception("PaymentMethodTokenRecipient is not available.")

        # Decrypt and return as output
        return Decryptor.paymentMethodTokenRecipient.unseal(encryptedMessage)
