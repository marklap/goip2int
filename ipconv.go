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
	"net"
)

var (
	minIPv4 = uint32(0)
	maxIPv4 = uint32(1<<32 - 1) // copied from standard lib math package constants (math.MaxUint32)
	minIPv6 = uint64(0)
	maxIPv6 = uint64(1<<64 - 1) // copied from standard lib math package constants (math.MaxUint64)
)

func ipv4ToUInt(ip net.IP) uint32 {
	n := []byte(ip.To4())
	return (uint32(n[0]) << 24) | (uint32(n[1]) << 16) | (uint32(n[2]) << 8) | uint32(n[3])
}

func uintToIPv4(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func ipv6ToUInts(ip net.IP) (network, host uint64) {
	n := []byte(ip.To16())
	network = (uint64(n[0]) << 56) | (uint64(n[1]) << 48) | (uint64(n[2]) << 40) | (uint64(n[3]) << 32) | (uint64(n[4]) << 24) | (uint64(n[5]) << 16) | (uint64(n[6]) << 8) | uint64(n[7])
	host = (uint64(n[8]) << 56) | (uint64(n[9]) << 48) | (uint64(n[10]) << 40) | (uint64(n[11]) << 32) | (uint64(n[12]) << 24) | (uint64(n[13]) << 16) | (uint64(n[14]) << 8) | uint64(n[15])
	return network, host
}

func uintsToIPv6(network, host uint64) net.IP {
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

func ipv4NetStartEnd(ip net.IP, mask int) (start, end uint32) {
	switch {
	case mask < 0:
		return minIPv4, maxIPv4
	case mask > 32:
		return maxIPv4, maxIPv4
	}

	i := ipv4ToUInt(ip)
	sm := uint32(maxIPv4 << uint32(32-mask))
	em := ^sm

	start = i & sm
	end = i | em
	return start, end
}

func ipv6NetStartEnd(ip net.IP, mask int) (sNet, sHost, eNet, eHost uint64) {

	n, h := ipv6ToUInts(ip)

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

func detectIPVersion(ip net.IP) int {
	switch {
	case ip.To4() != nil && ip.To16() != nil:
		return 4
	case ip.To4() == nil && ip.To16() != nil:
		return 6
	}
	return 0
}
