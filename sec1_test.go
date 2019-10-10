// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systemcerts

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

var ecKeyTests = []struct {
	derHex            string
	shouldReserialize bool
}{
	// Generated using:
	//   openssl ecparam -genkey -name secp384r1 -outform PEM
	{"3081a40201010430bdb9839c08ee793d1157886a7a758a3c8b2a17a4df48f17ace57c72c56b4723cf21dcda21d4e1ad57ff034f19fcfd98ea00706052b81040022a16403620004feea808b5ee2429cfcce13c32160e1c960990bd050bb0fdf7222f3decd0a55008e32a6aa3c9062051c4cba92a7a3b178b24567412d43cdd2f882fa5addddd726fe3e208d2c26d733a773a597abb749714df7256ead5105fa6e7b3650de236b50", true},
	// This key was generated by GnuTLS and has illegal zero-padding of the
	// private key. See https://golang.org/issues/13699.
	{"3078020101042100f9f43a04b9bdc3ab01f53be6df80e7a7bc3eaf7b87fc24e630a4a0aa97633645a00a06082a8648ce3d030107a1440342000441a51bc318461b4c39a45048a16d4fc2a935b1ea7fe86e8c1fa219d6f2438f7c7fd62957d3442efb94b6a23eb0ea66dda663dc42f379cda6630b21b7888a5d3d", false},
	// This was generated using an old version of OpenSSL and is missing a
	// leading zero byte in the private key that should be present.
	{"3081db0201010441607b4f985774ac21e633999794542e09312073480baa69550914d6d43d8414441e61b36650567901da714f94dffb3ce0e2575c31928a0997d51df5c440e983ca17a00706052b81040023a181890381860004001661557afedd7ac8d6b70e038e576558c626eb62edda36d29c3a1310277c11f67a8c6f949e5430a37dcfb95d902c1b5b5379c389873b9dd17be3bdb088a4774a7401072f830fb9a08d93bfa50a03dd3292ea07928724ddb915d831917a338f6b0aecfbc3cf5352c4a1295d356890c41c34116d29eeb93779aab9d9d78e2613437740f6", false},
}

func TestParseECPrivateKey(t *testing.T) {
	for i, test := range ecKeyTests {
		derBytes, _ := hex.DecodeString(test.derHex)
		key, err := ParseECPrivateKey(derBytes)
		if err != nil {
			t.Fatalf("#%d: failed to decode EC private key: %s", i, err)
		}
		serialized, err := MarshalECPrivateKey(key)
		if err != nil {
			t.Fatalf("#%d: failed to encode EC private key: %s", i, err)
		}
		matches := bytes.Equal(serialized, derBytes)
		if matches != test.shouldReserialize {
			t.Fatalf("#%d: when serializing key: matches=%t, should match=%t: original %x, reserialized %x", i, matches, test.shouldReserialize, serialized, derBytes)
		}
	}
}

const hexECTestPKCS1Key = "3082025c02010002818100b1a1e0945b9289c4d3f1329f8a982c4a2dcd59bfd372fb8085a9c517554607ebd2f7990eef216ac9f4605f71a03b04f42a5255b158cf8e0844191f5119348baa44c35056e20609bcf9510f30ead4b481c81d7865fb27b8e0090e112b717f3ee08cdfc4012da1f1f7cf2a1bc34c73a54a12b06372d09714742dd7895eadde4aa5020301000102818062b7fa1db93e993e40237de4d89b7591cc1ea1d04fed4904c643f17ae4334557b4295270d0491c161cb02a9af557978b32b20b59c267a721c4e6c956c2d147046e9ae5f2da36db0106d70021fa9343455f8f973a4b355a26fd19e6b39dee0405ea2b32deddf0f4817759ef705d02b34faab9ca93c6766e9f722290f119f34449024100d9c29a4a013a90e35fd1be14a3f747c589fac613a695282d61812a711906b8a0876c6181f0333ca1066596f57bff47e7cfcabf19c0fc69d9cd76df743038b3cb024100d0d3546fecf879b5551f2bd2c05e6385f2718a08a6face3d2aecc9d7e03645a480a46c81662c12ad6bd6901e3bd4f38029462de7290859567cdf371c79088d4f024100c254150657e460ea58573fcf01a82a4791e3d6223135c8bdfed69afe84fbe7857274f8eb5165180507455f9b4105c6b08b51fe8a481bb986a202245576b713530240045700003b7a867d0041df9547ae2e7f50248febd21c9040b12dae9c2feab0d3d4609668b208e4727a3541557f84d372ac68eaf74ce1018a4c9a0ef92682c8fd02405769731480bb3a4570abf422527c5f34bf732fa6c1e08cc322753c511ce055fac20fc770025663ad3165324314df907f1f1942f0448a7e9cdbf87ecd98b92156"
const hexECTestPKCS8Key = "30820278020100300d06092a864886f70d0101010500048202623082025e02010002818100cfb1b5bf9685ffa97b4f99df4ff122b70e59ac9b992f3bc2b3dde17d53c1a34928719b02e8fd17839499bfbd515bd6ef99c7a1c47a239718fe36bfd824c0d96060084b5f67f0273443007a24dfaf5634f7772c9346e10eb294c2306671a5a5e719ae24b4de467291bc571014b0e02dec04534d66a9bb171d644b66b091780e8d020301000102818100b595778383c4afdbab95d2bfed12b3f93bb0a73a7ad952f44d7185fd9ec6c34de8f03a48770f2009c8580bcd275e9632714e9a5e3f32f29dc55474b2329ff0ebc08b3ffcb35bc96e6516b483df80a4a59cceb71918cbabf91564e64a39d7e35dce21cb3031824fdbc845dba6458852ec16af5dddf51a8397a8797ae0337b1439024100ea0eb1b914158c70db39031dd8904d6f18f408c85fbbc592d7d20dee7986969efbda081fdf8bc40e1b1336d6b638110c836bfdc3f314560d2e49cd4fbde1e20b024100e32a4e793b574c9c4a94c8803db5152141e72d03de64e54ef2c8ed104988ca780cd11397bc359630d01b97ebd87067c5451ba777cf045ca23f5912f1031308c702406dfcdbbd5a57c9f85abc4edf9e9e29153507b07ce0a7ef6f52e60dcfebe1b8341babd8b789a837485da6c8d55b29bbb142ace3c24a1f5b54b454d01b51e2ad03024100bd6a2b60dee01e1b3bfcef6a2f09ed027c273cdbbaf6ba55a80f6dcc64e4509ee560f84b4f3e076bd03b11e42fe71a3fdd2dffe7e0902c8584f8cad877cdc945024100aa512fa4ada69881f1d8bb8ad6614f192b83200aef5edf4811313d5ef30a86cbd0a90f7b025c71ea06ec6b34db6306c86b1040670fd8654ad7291d066d06d031"

var ecMismatchKeyTests = []struct {
	hexKey        string
	errorContains string
}{
	{hexKey: hexECTestPKCS8Key, errorContains: "use ParsePKCS8PrivateKey instead"},
	{hexKey: hexECTestPKCS1Key, errorContains: "use ParsePKCS1PrivateKey instead"},
}

func TestECMismatchKeyFormat(t *testing.T) {
	for i, test := range ecMismatchKeyTests {
		derBytes, _ := hex.DecodeString(test.hexKey)
		_, err := ParseECPrivateKey(derBytes)
		if !strings.Contains(err.Error(), test.errorContains) {
			t.Errorf("#%d: expected error containing %q, got %s", i, test.errorContains, err)
		}
	}
}
