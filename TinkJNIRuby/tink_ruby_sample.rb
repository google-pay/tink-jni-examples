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

require 'yaml'
require_relative 'tink_ruby_decrypt'

class TinkRubySample
  SAMPLE_TOKEN = '<place_your_sample_token_here>'

  def main
    # 1. Load and validate the config file
    # Config file can be replaced as a hash or json format accordingly (just maintain as key and value pair)
    config = YAML.load_file('config.yml')
    if !config['Gateway_Name']
      raise 'Setting Gateway_Name must be set.'
    elsif !config['Protocol_Version']
      raise 'Setting Protocol_Version must be set.'
    elsif !config['Base64_PKCS8_Private_key']
      raise 'Setting Base64_PKCS8_Private_key must be set.'
    elsif !config['Trusted_Signing_Keys_Env']
      raise 'Setting Trusted_Signing_Keys_Env must be set.'
    elsif not ["INSTANCE_TEST", "INSTANCE_PRODUCTION"].include?(config['Trusted_Signing_Keys_Env'])
      raise 'Setting Trusted_Signing_Keys_Env must be either INSTANCE_TEST or INSTANCE_PRODUCTION.'
    end

    # 2. Print out payment token sample text
    puts "Sample Token: " + SAMPLE_TOKEN

    # 3. Decrypt and print out plain text
    puts "Plain Text: " + TinkRubyDecrypt.decrypt(SAMPLE_TOKEN, config)
  end
end

sample = TinkRubySample.new
sample.main