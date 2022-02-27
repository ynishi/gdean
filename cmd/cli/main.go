package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ynishi/gdean/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	emv := service.MaxEmv(0.6,
		[]int32{-100, 200, 300},
		[]int32{-100, 200, 300},
	)
	log.Printf("got max emv!:%d", emv)
	r := service.ReportMaxEmvResults()
	log.Printf("{%v}", r)
}
