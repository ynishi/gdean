package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ynishi/gdean/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	c, _ := service.InitConfig()
	service.DialInfo = fmt.Sprintf("%s:%d", c.Host, c.Port)
	emv := service.MaxEmv(0.6,
		[]int32{-100, 200, 300},
		[]int32{-100, 200, 300},
	)
	log.Printf("got max emv!:%d", emv)
	r := service.ReportMaxEmvResults()
	log.Printf("{%v}", r)
}
