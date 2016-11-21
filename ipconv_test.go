////////////////////////////////////////////////////////////////////////////////
// The MIT License (MIT)
//
// Copyright (c) 2016 Mark LaPerriere
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
////////////////////////////////////////////////////////////////////////////////
package ipconv

import (
	"math/big"
	"net"
	"testing"
)

// TestIPv4Conversion tests that an IPv4 address can be converted to a uint32 and back again
func TestIPv4Conversion(t *testing.T) {
	ipWant := net.ParseIP("192.168.100.100")
	uWant := uint32(3232261220)

	if uGot := IPv4ToUInt(ipWant); uGot != uWant {
		t.Errorf("Failed to convert IPv4 %s to correct uint32 - got %d, want %d", ipWant, uGot, uWant)
	}

	if ipGot := UintToIPv4(uWant); !ipGot.Equal(ipWant) {
		t.Errorf("Failed to convert uint32 %d to IPv4 - got %s, want %s", uWant, ipGot, ipWant)
	}
}

// TestIPv6Conversion tests that an IPv6 address can be converted to two uint64s (net and host)
// and back again
func TestIPv6Conversion(t *testing.T) {
	ipWant := net.ParseIP("fdaf:1285:6f0e:1262:6f0e:8a2e:0370:7334")
	uNetWant, uHostWant := uint64(18279849776823276130), uint64(8002485518114779956)
	bWant, _ := new(big.Int).SetString("337203710538915638676904557932570506036", 10)

	if uNetGot, uHostGot := IPv6ToUInts(ipWant); uNetGot != uNetWant || uHostGot != uHostWant || bWant.Cmp(IPv6ToBig(ipWant)) != 0 {
		t.Errorf("Failed to convert IPv6 %s to correct uint64s - got %d::%d, want %d::%d", ipWant, uNetGot, uHostGot, uNetWant, uHostWant)
	}

	if ipGot := UintsToIPv6(uNetWant, uHostWant); !ipGot.Equal(ipWant) {
		t.Errorf("Failed to convert uint64s %d::%d to IPv6 - got %s, want %s", uNetWant, uHostWant, ipGot, ipWant)
	}
}

// TestIPv4Subnetting tests that an IPv4 address range can be calculated from CIDR notation
func TestIPv4Subnetting(t *testing.T) {
	mBits := 27
	ipWant := net.ParseIP("192.168.100.100")
	sUWant := uint32(3232261216)
	eUWant := uint32(3232261247)

	if sUGot, eUGot := IPv4NetStartEnd(ipWant, mBits); sUGot != sUWant || eUGot != eUWant {
		t.Errorf("Failed to convert IPv4 CIDR %s/%d into correct start and end uint32s - got %d-%d, want %d-%d", ipWant, mBits, sUWant, eUWant, sUGot, eUGot)
	}
}

// TestIPv6NetSubnetting tests that an IPv6 address range can be calculated from CIDR notation
// where the mask only applies to the network bits
func TestIPv6NetSubnetting(t *testing.T) {
	// Mask: 33 -> 18279849774960082944::0-18279849777107566591::18446744073709551615 -> fdaf:1285::-fdaf:1285:7fff:ffff:ffff:ffff:ffff:ffff
	mBits := 33
	ipWant := net.ParseIP("fdaf:1285:6f0e:1262:6f0e:8a2e:0370:7334")
	uSNetWant, uSHostWant := uint64(18279849774960082944), uint64(0)
	uENetWant, uEHostWant := uint64(18279849777107566591), uint64(18446744073709551615)
	bSWant, _ := new(big.Int).SetString("337203710504545790806880554100409237504", 10)
	bEWant, _ := new(big.Int).SetString("337203710544159872064012722897181212671", 10)

	eq := func(sn, sh, en, eh uint64) bool {
		return sn == uSNetWant &&
			sh == uSHostWant &&
			en == uENetWant &&
			eh == uEHostWant &&
			bSWant.Cmp(uint64sToBig(sn, sh)) == 0 &&
			bEWant.Cmp(uint64sToBig(en, eh)) == 0
	}

	if sn, sh, en, eh := IPv6NetStartEnd(ipWant, mBits); !eq(sn, sh, en, eh) {
		t.Errorf("Failed to convert IPv6 CIDR %s/%d into correct start and end uint64s - got %d::%d-%d::%d, want %d::%d-%d::%d", ipWant, mBits, uSNetWant, uSHostWant, uENetWant, uEHostWant, sn, sh, en, eh)
	}
}
