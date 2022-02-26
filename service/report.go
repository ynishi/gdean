package service

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/ynishi/gdean/pb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// domain model
type ReportModel struct {
	gorm.Model
	MaxEmv     int
	MaxEmvDesc string
	Seconds    int64
	Nanos      int32
}

// report infra
var reportRepository = reporter()

func reporter() func(*pb.Result, bool) *pb.Report {
	db, _ := gorm.Open(sqlite.Open("report.db"), &gorm.Config{})
	db.AutoMigrate(&ReportModel{})

	return func(res *pb.Result, isFetch bool) *pb.Report {
		if isFetch {
			var rep []ReportModel
			db.Find(&rep)
			count := len(rep)
			res := make([]*pb.Result, count)
			for i, v := range rep {
				res[i] = &pb.Result{MaxEmv: int32(v.MaxEmv), CreateTime: &timestamp.Timestamp{
					Seconds: v.Seconds,
					Nanos:   v.Nanos,
				}}
			}
			return &pb.Report{
				NumberOfCalc: int32(count),
				Result:       res,
			}
		}
		db.Create(&ReportModel{MaxEmv: int(res.GetMaxEmv()), MaxEmvDesc: "desc", Seconds: res.GetCreateTime().Seconds, Nanos: res.GetCreateTime().Nanos})
		return nil
	}
}
