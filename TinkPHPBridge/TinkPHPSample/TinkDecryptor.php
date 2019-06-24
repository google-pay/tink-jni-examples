<?php 
    require_once("http://localhost:8081/JavaBridge/java/Java.inc");

    class TinkDecryptor {
        private function load_ini_files() {
            $output = array();
            if ($filelist = php_ini_scanned_files()) {
                if (strlen($filelist) > 0) {
                    $files = explode(',', $filelist);
                    foreach ($files as $file) {
                        $output = array_merge($output, parse_ini_file(trim($file)));
                    }
                }
            }
            return $output;
        }

        private function get_config() {
            $iniParams = $this->load_ini_files();
            if (!isset($iniParams['trusted_signing_keys_env'])) {
                throw new Exception('Setting trusted_signing_keys_env must be set.');
            }
            else if ($iniParams['trusted_signing_keys_env'] <> 'INSTANCE_TEST' && $iniParams['trusted_signing_keys_env'] <> 'INSTANCE_PRODUCTION') {
                throw new Exception('Setting trusted_signing_keys_env must be either INSTANCE_TEST or INSTANCE_PRODUCTION.');
            }
            else if (!isset($iniParams['gateway_name'])) {
                throw new Exception('Setting gateway_name must be set.');
            }
            else if (!isset($iniParams['base64_pkcs8_private_key'])) {
                throw new Exception('Setting base64_pkcs8_private_key must be set.');
            }
            else if (!isset($iniParams['protocol_version'])) {
                throw new Exception('Setting protocol_version must be set.');
            }
            return $iniParams;
        }

        public function decrypt($encryptedMessage) {
            $iniParams = $this->get_config();
            $recipientBuilder = new java('com.google.crypto.tink.apps.paymentmethodtoken.PaymentMethodTokenRecipient$Builder');
            
            if ($iniParams['trusted_signing_keys_env'] == "INSTANCE_TEST") {
                $keyManager = (new java('com.google.crypto.tink.apps.paymentmethodtoken.GooglePaymentsPublicKeysManager'))->INSTANCE_TEST;
            }
            else if ($iniParams['trusted_signing_keys_env'] == "INSTANCE_PRODUCTION") {
                $keyManager = (new java('com.google.crypto.tink.apps.paymentmethodtoken.GooglePaymentsPublicKeysManager'))->INSTANCE_PRODUCTION;
            }

            return $recipientBuilder->fetchSenderVerifyingKeysWith($keyManager)
                ->recipientId($iniParams['gateway_name'])
                ->protocolVersion($iniParams['protocol_version'])
                ->addRecipientPrivateKey($iniParams['base64_pkcs8_private_key'])
                ->build()
                ->unseal($encryptedMessage);
        }
    }

    try {
        $encryptedMessage = file_get_contents("tokendata.json");
        $tinkDecryptor = new TinkDecryptor;
        echo $tinkDecryptor->decrypt($encryptedMessage);
    }
    catch (Exception $e) {
        echo $e;
    }
?>