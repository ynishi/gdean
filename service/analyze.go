package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/ynishi/gdean/pb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// server impl
type AnalyzeServer struct {
	Server pb.AnalyzeServiceServer
	Repo   AnalyzeRepository
}

// constructors for server
func AnalyzeServerWithSqliteRepo(ctx context.Context, server pb.AnalyzeServiceServer, info string) *AnalyzeServer {
	return DefaultAnalyzeServerWithRepo(ctx, server, NewSqliteAnalyzeRepository(info))
}

func DefaultAnalyzeServerWithRepo(ctx context.Context, server pb.AnalyzeServiceServer, repo AnalyzeRepository) *AnalyzeServer {
	return &AnalyzeServer{Server: server, Repo: repo}
}

// internal domain(model) compatible with rdb
type Meta struct {
	gorm.Model
	Name        string
	Desc        string
	IsAvailable bool   `gorm:"column:is_available"`
	ParamDef    string `gorm:"column:param_def"`
}

type AnalyzeRepository interface {
	Fetch(id uint32) (*pb.Meta, error)
	Create(*pb.MetaBody) (*pb.Meta, error)
	Put(uint32, *pb.MetaBody) (*pb.Meta, error)
	Delete(uint32) (*pb.Meta, error)
	FetchFrom(uint32) ([]uint32, error)
}

// to inject db conn to repository
type AnalyzeDBGetter interface {
	GetDB() (*gorm.DB, error)
}

// to default func definition injectable conn
type DefaultAnalyzeRepository struct {
	DBGetter AnalyzeDBGetter
}

// embedded Default(expext to injecte sqlite conn)
type SqliteAnalyzeRepository struct {
	*DefaultAnalyzeRepository
}

// has conn info struct work with sqlite repository
type SqlteDBGetter struct {
	Info string
}

func NewSqliteAnalyzeRepository(info string) *SqliteAnalyzeRepository {
	return &SqliteAnalyzeRepository{DefaultAnalyzeRepository: &DefaultAnalyzeRepository{DBGetter: &SqlteDBGetter{Info: info}}}
}

func (s *DefaultAnalyzeRepository) Init() error {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return err
	}
	db.AutoMigrate(&Meta{})
	return nil
}

func (s *SqlteDBGetter) GetDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(s.Info), &gorm.Config{})
}

func (s *DefaultAnalyzeRepository) Fetch(id uint32) (*pb.Meta, error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	var meta Meta
	db.First(&meta, id)
	pbMeta, err := toPBMeta(&meta)
	if err != nil {
		return nil, err
	}
	return pbMeta, nil
}

func (s *DefaultAnalyzeRepository) Create(metaBody *pb.MetaBody) (*pb.Meta, error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	meta, err := fromPBMetaBody(metaBody)
	if err != nil {
		return nil, err
	}
	res := db.Create(meta)
	if res.Error != nil {
		return nil, err
	}
	pbMeta, err := toPBMeta(meta)
	if err != nil {
		return nil, err
	}
	return pbMeta, nil
}

func (s *DefaultAnalyzeRepository) Put(id uint32, metaBody *pb.MetaBody) (*pb.Meta, error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	paramDef, err := json.Marshal(metaBody.ParamDef)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var meta Meta
	db.First(&meta, id)
	meta.Name = metaBody.Name
	meta.Desc = metaBody.Desc
	meta.IsAvailable = metaBody.IsAvailable
	meta.ParamDef = string(paramDef)
	meta.UpdatedAt = now
	res := db.Save(&meta)
	if res.Error != nil {
		return nil, err
	}
	pbMeta, err := toPBMeta(&meta)
	if err != nil {
		return nil, err
	}
	return pbMeta, nil
}

func (s *DefaultAnalyzeRepository) Delete(id uint32) (*pb.Meta, error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	var meta Meta
	db.First(&meta, id)
	res := db.Delete(&meta)
	if res.Error != nil {
		return nil, res.Error
	}
	pbMeta, err := toPBMeta(&meta)
	if err != nil {
		return nil, err
	}
	return pbMeta, nil
}

func (s *DefaultAnalyzeRepository) FetchFrom(id uint32) ([]uint32, error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	metas := make([]Meta, 10)
	db.Limit(10).Find(&metas)
	var acc []uint32
	for _, meta := range metas {
		acc = append(acc, uint32(meta.ID))
	}
	return acc, nil
}

func (s *AnalyzeServer) GetMeta(ctx context.Context, in *pb.GetMetaRequest) (*pb.GetMetaResponse, error) {
	meta, err := s.Repo.Fetch(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetMetaResponse{Response: &pb.GetMetaResponse_Meta{Meta: meta}}, nil
}

func (s *AnalyzeServer) CreateMeta(ctx context.Context, in *pb.CreateMetaRequest) (*pb.CreateMetaResponse, error) {
	meta, err := s.Repo.Create(in.MetaBody)
	if err != nil {
		return nil, err
	}
	return &pb.CreateMetaResponse{Response: &pb.CreateMetaResponse_Meta{meta}}, nil
}

func (s *AnalyzeServer) PutMeta(ctx context.Context, in *pb.PutMetaRequest) (*pb.PutMetaResponse, error) {
	meta, err := s.Repo.Put(in.Id, in.MetaBody)
	if err != nil {
		return nil, err
	}
	return &pb.PutMetaResponse{Response: &pb.PutMetaResponse_UpdateTime{meta.UpdateTime}}, nil
}

func (s *AnalyzeServer) DeleteMeta(ctx context.Context, in *pb.DeleteMetaRequest) (*pb.DeleteMetaResponse, error) {
	meta, err := s.Repo.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMetaResponse{Response: &pb.DeleteMetaResponse_DeleteTime{meta.UpdateTime}}, nil
}

// TODO
func (s *AnalyzeServer) GetMetaList(ctx context.Context, in *pb.GetMetaListRequest) (*pb.GetMetaListResponse, error) {
	ids, err := s.Repo.FetchFrom(in.StartId)
	if err != nil {
		return nil, err
	}
	return &pb.GetMetaListResponse{Response: &pb.GetMetaListResponse_Ids{&pb.GetMetaListIds{Ids: ids}}}, nil
}

// TODO
func (s *AnalyzeServer) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	return nil, nil
}

// helper converters
func toPBMeta(meta *Meta) (*pb.Meta, error) {
	var paramDef map[string]string
	err := json.Unmarshal([]byte(meta.ParamDef), &paramDef)
	if err != nil {
		return nil, err
	}
	metaBody := pb.MetaBody{Name: meta.Name, Desc: meta.Desc, ParamDef: paramDef, IsAvailable: meta.IsAvailable}
	ct, err := ptypes.TimestampProto(meta.CreatedAt)
	if err != nil {
		return nil, err
	}
	ut, err := ptypes.TimestampProto(meta.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &pb.Meta{Id: uint32(meta.ID), MetaBody: &metaBody, CreateTime: ct, UpdateTime: ut}, nil
}

func fromPBMetaBody(metaBody *pb.MetaBody) (*Meta, error) {
	paramDef, err := json.Marshal(metaBody.ParamDef)
	if err != nil {
		return nil, err
	}
	return &Meta{Name: metaBody.Name, Desc: metaBody.Desc, IsAvailable: metaBody.IsAvailable, ParamDef: string(paramDef)}, nil
}
