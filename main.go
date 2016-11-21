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
	"flag"
	"fmt"
	"os"

	"./ipconv"
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
	flagBigInt  bool
	flagIPv6    bool
	flagVersion bool
)

func init() {
	flag.BoolVar(&flagIPv6, "6", false, "input argument is an IPv6 address")
	flag.BoolVar(&flagBigInt, "big", false, "output a big int instead of a pair of 64 bit unsigned ints")
	flag.BoolVar(&flagVersion, "version", false, "print the version and exit")
	flag.Parse()
}

func mode(arg string) int {
	return -1
}

func int2IP(arg string) string {
	return ""
}

func ip2Int(arg string) string {
	return ""
}

func cidr2IPRange(arg string) string {
	return ""
}

func main() {
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if flagVersion {
		fmt.Println(ipconv.Version)
		os.Exit(0)
	}

	arg := flag.Arg(0)
	switch mode(arg) {
	case -1:
		fmt.Println("Under Construction")
	case mInt2IP:
		fmt.Println(int2IP(arg))
	case mIP2Int:
		fmt.Println(ip2Int(arg))
	case mCIDR2IPRange:
		fmt.Println(cidr2IPRange(arg))
	}
}
