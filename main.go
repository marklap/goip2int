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
package main

import (
	"./ipconv"
	"flag"
	"math"
	"math/big"
	"net"
	"os"
)

const (
	mInt2IP int = iota
	mIP2Int
	mCIDR2IPRange
	t32
	t64
	tBig
)

var (
	fBigInt bool
	fIPv6   bool
)

func init() {
	flag.BoolVar(&fBigInt, "bigint", false, `(IPv6 ONLY) return "big" integers`)
	flag.BoolVar(&fIPv6, "6", false, `Return IPv6 address if possible`)
	flag.Parse()
}

func determineMode(arg string) int {

	return 0
}

func inspectIP(arg string) (net.IP, int) {
	ip := net.ParseIP(arg)
	if ip == nil {
		return nil, 0
	}

	return ip, determineIPVersion(ip)
}

func arg2Int(arg string, reqType int) interface{} {
	var i, max32, max64 = new(big.Int), new(big.Int), new(big.Int)
	max32.SetUint64(uint64(math.MaxUint32))
	max64.SetUint64(math.MaxUint64)

	if i, success := i.SetString(arg, 10); !success {
		return nil
	}

	switch reqType {
	case t32:
		if i.Cmp(max32) <= 0 {
			return uint32(i.Uint64())
		} else {
			return uint32(0)
		}
	case t64:
		if i.Cmp(max64) <= 0 {
			return i.Uint64()
		} else {
			return uint64(0)
		}
	case tBig:
		return i
	}
}

func int2IP(arg string) string {

}

func ip2Int(arg string) string {
	var res string

	ip, ver := inspectIP(arg)
	switch ver {
	case 4:
		return fmt.Sprintf("%s", ipconv.IPv4ToUInt(ip))
	case 6:
		network, host := ipconv.IPv6ToUints(ip)
		if fBigInt {
			return ipconv.UintsToBig(network, host).String()
		} else {
			return fmt.Sprintf("%d::%d", network, host)
		}
	}
}

func cidr2IPRange(arg string) string {

}

func main() {

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var res string
	arg := flag.Arg(0)
	switch determineMode(arg) {
	case mInt2IP:
		fmt.Println(int2IP(arg))
	case mIP2Int:
		fmt.Println(ip2Int(arg))
	case mCIDR2IPRange:
		fmt.Println(cidr2IPRange(arg))
	}
}
