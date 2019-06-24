/*
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import php.java.bridge.http.JavaBridgeRunner;

public class TinkDecryptor {
  public static void main(String args[]) throws Exception {
    if (args.length != 1) {
      System.out.println("ERROR   : Please specify the <JavaBridgePortNumber> in the argument.");
      System.out.println("USAGE   : java -jar TinkDecryptor.jar <JavaBridgePortNumber>");
      System.out.println("EXAMPLE : java -jar TinkDecryptor.jar 8081");
      System.exit(1);
    } else {
      System.out.println("PHP/JavaBridge running on Port Number " + args[0]);
      final JavaBridgeRunner runner = JavaBridgeRunner.getInstance(args[0]);
      runner.waitFor();
      System.exit(0);
    }
  }
}
