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

require 'rjb'

# A facade for payment method token recipient based on Tink library https://github.com/google/tink.
# It uses Ruby Java Bridge library to connect between Ruby and Java with Java Native Interface.
# Please ensure Java 8 or later is installed on your server.
class TinkRubyDecrypt
  def self.decrypt(encrypted_message, config)
    # create bridge and load all the required jar files
    Rjb::add_classpath(Dir.glob(__dir__ + "/java_libs/*.jar").join(':'))
    Rjb::load()

    # build key manager based on environment
    key_manager_class = Rjb::import('com.google.crypto.tink.apps.paymentmethodtoken.GooglePaymentsPublicKeysManager')
    key_manager = case config['Trusted_Signing_Keys_Env']
                  when 'INSTANCE_TEST' then key_manager_class.INSTANCE_TEST
                  when 'INSTANCE_PRODUCTION' then key_manager_class.INSTANCE_PRODUCTION
                  end

    # build payment method token recipient
    builder = Rjb::import('com.google.crypto.tink.apps.paymentmethodtoken.PaymentMethodTokenRecipient$Builder').new
    builder.fetchSenderVerifyingKeysWith(key_manager)
    builder.recipientId('gateway:' + config['Gateway_Name'])
    builder.protocolVersion(config['Protocol_Version'])
    builder.addRecipientPrivateKey(config['Base64_PKCS8_Private_key'])
    recipient = builder.build

    # decrypt and return as output
    recipient.unseal(encrypted_message)
  end
end