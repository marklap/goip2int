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
)

// Version is the version
const Version = "0.1.0"

var (
	minIPv4 = uint32(0)
	maxIPv4 = uint32(1<<32 - 1) // copied from standard lib math package constants (math.MaxUint32)
	minIPv6 = uint64(0)
	maxIPv6 = uint64(1<<64 - 1) // copied from standard lib math package constants (math.MaxUint64)
)

// IPv4ToUInt converts an IPv4 address to the equivalent integer
func IPv4ToUInt(ip net.IP) uint32 {
	n := []byte(ip.To4())
	return (uint32(n[0]) << 24) | (uint32(n[1]) << 16) | (uint32(n[2]) << 8) | uint32(n[3])
}

// UintToIPv4 converts an integer into it's equivalent IPv4 address
func UintToIPv4(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

// IPv6ToUInts converts an IPv6 address to the equivalent set of integers
func IPv6ToUInts(ip net.IP) (network, host uint64) {
	n := []byte(ip.To16())
	network = (uint64(n[0]) << 56) | (uint64(n[1]) << 48) | (uint64(n[2]) << 40) | (uint64(n[3]) << 32) | (uint64(n[4]) << 24) | (uint64(n[5]) << 16) | (uint64(n[6]) << 8) | uint64(n[7])
	host = (uint64(n[8]) << 56) | (uint64(n[9]) << 48) | (uint64(n[10]) << 40) | (uint64(n[11]) << 32) | (uint64(n[12]) << 24) | (uint64(n[13]) << 16) | (uint64(n[14]) << 8) | uint64(n[15])
	return network, host
}

func uint64sToBig(left, right uint64) *big.Int {
	l := new(big.Int).SetUint64(left)
	r := new(big.Int).SetUint64(right)
	o := new(big.Int)
	o.Lsh(l, 64)
	o.Or(o, r)
	return o
}

func bigToUint64s(i *big.Int) (left, right uint64) {
	m := new(big.Int).SetUint64(maxIPv6)
	left = new(big.Int).Rsh(i, 64).Uint64()
	right = new(big.Int).And(i, m).Uint64()
	return left, right
}

// IPv6ToBig converts an IPv6 address to an equivalent `math.big.Int`
func IPv6ToBig(ip net.IP) *big.Int {
	network, host := IPv6ToUInts(ip)
	return uint64sToBig(network, host)
}

// BigToIPv6 converts a `math.big.Int` into the equivalent IPv6 address
func BigToIPv6(i *big.Int) net.IP {
	n, h := bigToUint64s(i)
	return UintsToIPv6(n, h)
}

// UintsToIPv6 converts a pair of integers to an IPv6 address
func UintsToIPv6(network, host uint64) net.IP {
	var ip = make(net.IP, net.IPv6len)
	ip[0] = byte(network >> 56)
	ip[1] = byte(network >> 48)
	ip[2] = byte(network >> 40)
	ip[3] = byte(network >> 32)
	ip[4] = byte(network >> 24)
	ip[5] = byte(network >> 16)
	ip[6] = byte(network >> 8)
	ip[7] = byte(network)
	ip[8] = byte(host >> 56)
	ip[9] = byte(host >> 48)
	ip[10] = byte(host >> 40)
	ip[11] = byte(host >> 32)
	ip[12] = byte(host >> 24)
	ip[13] = byte(host >> 16)
	ip[14] = byte(host >> 8)
	ip[15] = byte(host)
	return ip
}

// IPv4NetStartEnd determines the start and end IP address (represented in integers) based on an
// IPv4 address and an IPv4 subnet mask in CIDR notation (eg. "192.168.1.1/24", `mask` is 24).
func IPv4NetStartEnd(ip net.IP, mask int) (start, end uint32) {
	switch {
	case mask < 0:
		return minIPv4, maxIPv4
	case mask > 32:
		return maxIPv4, maxIPv4
	}

	i := IPv4ToUInt(ip)
	sm := uint32(maxIPv4 << uint32(32-mask))
	em := ^sm

	start = i & sm
	end = i | em
	return start, end
}

// IPv6NetStartEnd determines the start
func IPv6NetStartEnd(ip net.IP, mask int) (sNet, sHost, eNet, eHost uint64) {

	n, h := IPv6ToUInts(ip)

	var subj, cMask uint64
	switch {
	case mask < 0:
		return minIPv6, maxIPv6, minIPv6, maxIPv6
	case mask > 128:
		return maxIPv6, maxIPv6, maxIPv6, maxIPv6
	case mask <= 64:
		subj = n
		cMask = 64 - uint64(mask)
	case mask > 64 && mask <= 128:
		subj = h
		cMask = uint64(mask) - 64
	}

	sm := uint64(maxIPv6 << cMask)
	em := ^sm

	start := subj & sm
	end := subj | em

	switch {
	case mask >= 0 && mask <= 64:
		sNet = start
		sHost = minIPv6
		eNet = end
		eHost = maxIPv6
	case mask > 64 && mask <= 128:
		sNet = n
		sHost = start
		eNet = n
		eHost = end
	}
	return sNet, sHost, eNet, eHost
}

// IPv6NetStartEndBig is the same as `IPv6NetStartEnd` except that it returns a pair of
// `*math.big.Int` values instead of a 4-tuple of uint64
func IPv6NetStartEndBig(ip net.IP, mask int) (start, end *big.Int) {
	sn, sh, en, eh := IPv6NetStartEnd(ip, mask)
	start = uint64sToBig(sn, sh)
	end = uint64sToBig(en, eh)
	return start, end
}

// DetectIPVersion determines if the IP address is IPv4 or IPv6
func DetectIPVersion(ip net.IP) int {
	switch {
	case ip.To4() != nil && ip.To16() != nil:
		return 4
	case ip.To4() == nil && ip.To16() != nil:
		return 6
	}
	return 0
}
