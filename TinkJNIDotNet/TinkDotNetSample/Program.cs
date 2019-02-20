// Copyright 2018 Google LLC
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

using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace TinkDotNetSample
{
    class Program
    {
        /// <summary>
        /// Sample token generated by Google Pay test web or app.
        /// Token can be retrieved in "paymentData.paymentMethodData.tokenizationData.token" (refer to example in https://developers.google.com/pay/api/web/guides/tutorial#full-example)
        /// </summary>
        private const string _paymentSampleToken = "{\"signature\":\"MEUCIAfhO2itOXRj8bbbDUckv6PikE39OLnmAYmzKIh7s2D1AiEA5U9Tnsm7yRiMbUjtKR20dPV7KHnTmp+LBPCby4LGCN8\\u003d\",\"protocolVersion\":\"ECv1\",\"signedMessage\":\"{\\\"encryptedMessage\\\":\\\"VdUqx6QEUuuPXc3dGVFp0tdtMkKu7kpRGao9i8l5rHhKD4/eoyIXU7DqqhBlZqcRrRFNeJHFAPSetOeDfMhDuo8HAjkK470qIp/kxhBXYU9pxOk92ansUywHdxZxoblIXcMOoRRk7ynADnPljhRPyUsCJTTMzT3WyfqY+/8RH6lLJKeu6vXsMMIFFzb8G9xytm/P+SGn7z6s7fWm3+djqbv9fQ4blAVH1kiSZcSM/IZAGbXzurH5VQJjKMq52MKXUvC8qGbOukJpYFE25u3dInHX5BfHNuIic7P8TBP58G6V6yixx9uMAGJ42SK68zYYXpsyfadJmDRO5Lj/IV6u6389OmsB55jA+/PfFTZH3w3+TgFx1E3zPikWtyQE/QGjtnRTQp9JpDYxaUaK9SXLjED88oT/Jpbkgt087byTZ5R9Igzo/hFJtX7BY4PP741sGwZmUWaiTA\\\\u003d\\\\u003d\\\",\\\"ephemeralPublicKey\\\":\\\"BLwlsXwoTI2OLo2k90MRyPZMyTbB9LIPSAe074CFJuVgjCv0yb514oXF27wVFhMUJ45f9+914Z53qFLQ8hXEMRA\\\\u003d\\\",\\\"tag\\\":\\\"SMW1ffGKDXoGFHVkoyIw2ACxpi5vclYYtHvVN6cR3vI\\\\u003d\\\"}\"}";

        /// <summary>
        /// Main method
        /// </summary>
        /// <param name="args">arguments</param>
        static void Main(string[] args)
        {
            // initialilze decryptor object to make trusted signing keys available in library
            var decryptor = new TinkDecryptor();
            // print out payment sample text
            Console.Write("Sample Token: ");
            Console.WriteLine(_paymentSampleToken);
            Console.WriteLine();
            // decrypt and print out plain text
            Console.Write("Plain Text: ");
            Console.WriteLine(decryptor.Decrypt(_paymentSampleToken)); // "expired payload" exception will be thrown if the message is expired
        }
    }
}
