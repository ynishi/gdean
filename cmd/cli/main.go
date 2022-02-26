package main

import (
	"fmt"
	"strconv"

	"github.com/ynishi/gdean/service"
)

func main() {
	emv := service.MaxEmv(0.6,
		[]int32{-100, 200, 300},
		[]int32{-100, 200, 300},
	)
	fmt.Println("got max emv!:" + strconv.Itoa(int(emv)))
	r := service.ReportMaxEmvResults()
	fmt.Printf("{%v}", r)
}
