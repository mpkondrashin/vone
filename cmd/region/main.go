package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mpkondrashin/vone"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: region <token> (or \"-\" for stdin)")
		os.Exit(2)
	}
	token := os.Args[1]
	if token == "-" {
		var sb strings.Builder
		io.Copy(&sb, os.Stdin)
		token = strings.TrimSpace(sb.String())
	}
	d, err := vone.DetectVisionOneDomain(context.TODO(), token, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if d == "" {
		fmt.Println("Domain not detected (wrong token?)")
		os.Exit(1)
	}
	fmt.Printf("Domain: %s\n", d)
}
