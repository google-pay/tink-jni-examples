<?php 
<<<<<<< HEAD
    require_once("http://localhost:8081/JavaBridge/java/Java.inc");

    class TinkDecryptor {
        // load ini file into key-value pairs
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

        // initialize settings from ini files
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
=======
    class TinkDecryptor {
        // validate settings from ini files
        public function check_config() {
            if (!get_cfg_var("trusted_signing_keys_env")) {
                throw new Exception('Setting trusted_signing_keys_env must be set.');
            }
            else if (get_cfg_var("trusted_signing_keys_env") <> 'INSTANCE_TEST' && get_cfg_var("trusted_signing_keys_env") <> 'INSTANCE_PRODUCTION') {
                throw new Exception('Setting trusted_signing_keys_env must be either INSTANCE_TEST or INSTANCE_PRODUCTION.');
            }
            else if (!get_cfg_var("gateway_name")) {
                throw new Exception('Setting gateway_name must be set.');
            }
            else if (!get_cfg_var("base64_pkcs8_private_key")) {
                throw new Exception('Setting base64_pkcs8_private_key must be set.');
            }
            else if (!get_cfg_var("protocol_version")) {
                throw new Exception('Setting protocol_version must be set.');
            }
            else if (!get_cfg_var("java_bridge_servlet_path")) {
                throw new Exception('Setting java_bridge_servlet_path must be set.');
            }
            else if (!get_cfg_var("java_bridge_local_path")) {
                throw new Exception('Setting java_bridge_local_path must be set.');
            }
        }

        // Initialize java bridge library and load into local library file
        // Note: This method should be called once during server start up only
        // For subsequent call, it is recommended to load the library from cache instead
        public function load_java_bridge_lib() {
            $servlet = get_cfg_var("java_bridge_servlet_path");
            $local = get_cfg_var("java_bridge_local_path");

            if (file_exists($local) && !is_writable($local)) {
                throw new Exception("java_bridge_local_path: {$local} write access denied.");
            }
            else {
                $remote_contents = file_get_contents($servlet);
                file_put_contents($local, $remote_contents);
            }
            require_once($local);
>>>>>>> 1432c37e256aa98b63d754d8db6c165594c75b89
        }

        // decrypt the encrypted payload
        public function decrypt($encryptedMessage) {
<<<<<<< HEAD
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
=======
            $recipientBuilder = new java('com.google.crypto.tink.apps.paymentmethodtoken.PaymentMethodTokenRecipient$Builder');
            $keyManager = new javaclass('com.google.crypto.tink.apps.paymentmethodtoken.GooglePaymentsPublicKeysManager');
            
            if (get_cfg_var("trusted_signing_keys_env") == "INSTANCE_TEST") {
                $keyManager = $keyManager->INSTANCE_TEST;
            }
            else if (get_cfg_var("trusted_signing_keys_env") == "INSTANCE_PRODUCTION") {
                $keyManager = $keyManager->INSTANCE_PRODUCTION;
            }

            return $recipientBuilder->fetchSenderVerifyingKeysWith($keyManager)
                ->recipientId(get_cfg_var("gateway_name"))
                ->protocolVersion(get_cfg_var("protocol_version"))
                ->addRecipientPrivateKey(get_cfg_var("base64_pkcs8_private_key"))
>>>>>>> 1432c37e256aa98b63d754d8db6c165594c75b89
                ->build()
                ->unseal($encryptedMessage);
        }
    }
    

    // test decryption with sample payload
    try {
<<<<<<< HEAD
        $encryptedMessage = file_get_contents("SamplePayload.json");
        $tinkDecryptor = new TinkDecryptor;
=======
        $encryptedMessage = file_get_contents("tokendata.json");
        $tinkDecryptor = new TinkDecryptor;
        $tinkDecryptor->check_config();
        $tinkDecryptor->load_java_bridge_lib();
>>>>>>> 1432c37e256aa98b63d754d8db6c165594c75b89
        echo $tinkDecryptor->decrypt($encryptedMessage);
    }
    catch (Exception $e) {
        echo $e;
    }
?>
