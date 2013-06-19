package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/lindroth/socknet/lib"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	delim = '\n'
)

func setHeader(flags *string) http.Header {
	header := http.Header{}

	for _, flag := range strings.Split(*flags, ",") {
		value := strings.Split(flag, "=")
		if len(value) == 2 {
			fmt.Printf("flags %s = %s\n", value[0], value[1])
			header.Add(value[0], value[1])
		}
	}

	return header
}

func main() {
	headerFlags := flag.String("header", "", "X-MyHeader=123,X-MyOtherHeader=654")
	flag.Parse()
	header := setHeader(headerFlags)

	args := flag.Args()

	if len(args) < 2 {
		fmt.Printf("usage: socknet origin url \n")
		os.Exit(1)
	}

	origin := args[0]
	url := args[1]

	fmt.Printf("origin %s, url %s \n", origin, url)
	client := &lib.Socknet{}

	input, output, err := client.Connect(origin, url, header)

	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(os.Stdin)

	go func() {
		for {
			println("enter string:")
			if line, err := r.ReadString(delim); err == nil {
				input <- line
			} else {
				fmt.Println(err)
				os.Exit(1)
			}

		}
	}()

	for msg := range output {
		fmt.Printf("received: %s \n", msg)
	}

	fmt.Printf("Closing .. \n")

}
