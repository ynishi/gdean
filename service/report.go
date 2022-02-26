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

type ReportRepository interface {
	Fetch() *pb.Report
	Put(*pb.Result)
}

var Repo ReportRepository

type Sqlite3ReportRepository struct {
	filename Sqlite3ConnInfo
}

type Sqlite3ConnInfo string

func DefaultSqlite3ConnInfo() *Sqlite3ConnInfo {
	var v Sqlite3ConnInfo = "report.db"
	return &v
}

func NewSqlite3ReportRepository(filename *Sqlite3ConnInfo) *Sqlite3ReportRepository {
	r := Sqlite3ReportRepository{filename: *filename}
	db := r.getDB()
	db.AutoMigrate(&ReportModel{})
	return &r
}

func (r *Sqlite3ReportRepository) getDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(string(r.filename)), &gorm.Config{})
	return db
}

func (r Sqlite3ReportRepository) Fetch() *pb.Report {
	var rep []ReportModel
	db := r.getDB()
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

func (r Sqlite3ReportRepository) Put(res *pb.Result) {
	db := r.getDB()
	db.Create(&ReportModel{MaxEmv: int(res.GetMaxEmv()), MaxEmvDesc: "desc", Seconds: res.GetCreateTime().Seconds, Nanos: res.GetCreateTime().Nanos})
}
