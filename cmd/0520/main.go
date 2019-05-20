package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/speps/go-hashids"
)

var (
	h bool

	t string
	s string
)

func init() {
	flag.BoolVar(&h, "h", false, "show document")

	flag.StringVar(&t, "t", "decrypt", "operation: encrypt, decrypt")
	flag.StringVar(&s, "s", "", "the string or integer you want to process")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `hashids version: hashids/1.0.0
Usage: hashids [-h] [-t type] [-s string]
Example: hashids -t encrypt -s 2001

Options:
`)
	flag.PrintDefaults()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
	} else {
		flag.Parse()

		if len(s) == 0 {
			fmt.Println("No string.")
		}

		hd := hashids.NewData()
		hd.Salt = "bed0481cea6b4855a60b8d0133e0cc85"
		hd.MinLength = 6
		h, _ := hashids.NewWithData(hd)

		if t == "encrypt" {
			i, _ := strconv.Atoi(s)
			e, err := h.Encode([]int{i})
			checkErr(err)
			fmt.Println("Encrypted:", e)
		} else {
			e, err := h.DecodeWithError(s)
			checkErr(err)
			fmt.Println("Decrypted:", e[0])
		}
	}
}
