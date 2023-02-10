package main

import (
	"fmt"

	"github.com/mpkondrashin/vone"
)

func main() {
	v1 := vone.NewVOne(url, token)
	reserve, err := v1.SandboxGetDailyReserve()
	fmt.Println(err)
	fmt.Println(reserve)
}
