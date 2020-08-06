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

from tinkjni import Decryptor
import json

# 1. Load config file (Path is relative to the current working directory)
with open('./config.json') as f:
    config = json.load(f)

# 2. Construct the decryptor with config
# NOTE: load JVM once only after application has started instead of loading JVM for each and every decryption
Decryptor.init(config)

# 3. Decrypt message and display result
sampleToken = "<place_your_token_here>"
print("Clear text after decryption:")
print(Decryptor.decrypt(sampleToken))
