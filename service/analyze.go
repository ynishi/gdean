package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog/log"
	"github.com/thanos-io/thanos/pkg/runutil"
	pb "github.com/ynishi/gdean/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// server impl
type AnalyzeServer struct {
	pb.AnalyzeServiceServer
	Repo AnalyzeRepository
}

func DefaultAnalyzeServiceServer() *AnalyzeServer {
	return &AnalyzeServer{}
}

// constructors for server
func DefaultAnalyzeServerWithRepo(ctx context.Context, repo AnalyzeRepository) *AnalyzeServer {
	server := DefaultAnalyzeServiceServer()
	if err := repo.Init(); err != nil {
		return nil
	}
	server.Repo = repo
	return server
}

// impl for GRPC interface
func (s *AnalyzeServer) GetMeta(ctx context.Context, in *pb.GetMetaRequest) (res *pb.GetMetaResponse, err error) {
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
	return &pb.CreateMetaResponse{Response: &pb.CreateMetaResponse_Meta{Meta: meta}}, nil
}

func (s *AnalyzeServer) PutMeta(ctx context.Context, in *pb.PutMetaRequest) (*pb.PutMetaResponse, error) {
	meta, err := s.Repo.Put(in.Id, in.MetaBody)
	if err != nil {
		return nil, err
	}
	return &pb.PutMetaResponse{Response: &pb.PutMetaResponse_UpdateTime{UpdateTime: meta.UpdateTime}}, nil
}

func (s *AnalyzeServer) DeleteMeta(ctx context.Context, in *pb.DeleteMetaRequest) (*pb.DeleteMetaResponse, error) {
	meta, err := s.Repo.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMetaResponse{Response: &pb.DeleteMetaResponse_DeleteTime{DeleteTime: meta.UpdateTime}}, nil
}

func (s *AnalyzeServer) GetMetaList(ctx context.Context, in *pb.GetMetaListRequest) (*pb.GetMetaListResponse, error) {
	ids, err := s.Repo.FetchFrom(in.StartId)
	if err != nil {
		return nil, err
	}
	return &pb.GetMetaListResponse{Response: &pb.GetMetaListResponse_Ids{Ids: &pb.GetMetaListIds{Ids: ids}}}, nil
}

// TODO
func (s *AnalyzeServer) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	return nil, nil
}

func (s *AnalyzeServer) MaxEmv(ctx context.Context, in *pb.MaxEmvRequest) (*pb.MaxEmvResponse, error) {
	log.Debug().Float32("p1", float32(in.GetTowPData().GetP1())).Msg("Recieved")
	maxEmv, err := calcMaxEmv(in.GetTowPData().GetP1(), in.GetTowPData().GetDataP1(), in.GetTowPData().GetDataP2())
	if err != nil {
		return nil, err
	}
	result := pb.Result{MaxEmv: maxEmv, CreateTime: timestamppb.Now()}
	return &pb.MaxEmvResponse{MaxEmv: result.MaxEmv, CreateTime: result.CreateTime}, nil
}

// internal domain(model) compatible with rdb
type Meta struct {
	gorm.Model
	Name        string
	Desc        string
	IsAvailable bool   `gorm:"column:is_available"`
	ParamDef    string `gorm:"column:param_def"`
}

// TODO: consider to refactor Repository, simplifly management DB inject or remove Repository interface(this interface can used when change to not rdb)
// repository
type AnalyzeRepository interface {
	Init() error
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
	Cache    *lru.Cache
}

func (s *DefaultAnalyzeRepository) Init() (err error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return err
	}
	dbx, err := db.DB()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			runutil.CloseWithErrCapture(&err, dbx, "close conn")
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		return db.AutoMigrate(&Meta{})
	})
	if err != nil {
		return err
	}
	cache, err := lru.New(128)
	if err != nil {
		return err
	}
	s.Cache = cache
	return nil
}

func (s *DefaultAnalyzeRepository) Fetch(id uint32) (*pb.Meta, error) {
	if s.Cache != nil {
		if v, ok := s.Cache.Get(id); ok {
			return v.(*pb.Meta), nil
		}
	}
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	var meta Meta
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		return db.First(&meta, id).Error
	})
	if err != nil {
		return nil, err
	}
	pbMeta, err := toPBMeta(&meta)
	if err != nil {
		return nil, err
	}
	if s.Cache != nil {
		s.Cache.Add(id, pbMeta)
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
	if err := db.Create(meta).Error; err != nil {
		return nil, err
	}
	pbMeta, err := toPBMeta(meta)
	if err != nil {
		return nil, err
	}
	if s.Cache != nil {
		s.Cache.Add(pbMeta.Id, pbMeta)
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
	var meta Meta
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		return db.First(&meta, id).Error
	})
	if err != nil {
		return nil, err
	}
	meta.Name = metaBody.Name
	meta.Desc = metaBody.Desc
	meta.IsAvailable = metaBody.IsAvailable
	meta.ParamDef = string(paramDef)
	meta.UpdatedAt = time.Now()
	ctxSave, cancelSave := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelSave()
	err = runutil.Retry(5*time.Second, ctxSave.Done(), func() error {
		return db.Save(&meta).Error
	})
	if err != nil {
		return nil, err
	}
	pbMeta, err := toPBMeta(&meta)
	if err != nil {
		return nil, err
	}
	if s.Cache != nil {
		s.Cache.Add(id, pbMeta)
	}
	return pbMeta, nil
}

func (s *DefaultAnalyzeRepository) Delete(id uint32) (*pb.Meta, error) {
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	var meta Meta
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		return db.First(&meta, id).Error
	})
	if err != nil {
		return nil, err
	}
	ctxDel, cancelDel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelDel()
	err = runutil.Retry(5*time.Second, ctxDel.Done(), func() error {
		return db.Delete(&meta).Error
	})
	if err != nil {
		return nil, err
	}
	pbMeta, err := toPBMeta(&meta)
	if err != nil {
		return nil, err
	}
	if s.Cache != nil {
		s.Cache.Remove(id)
	}
	return pbMeta, nil
}

func (s *DefaultAnalyzeRepository) FetchFrom(id uint32) ([]uint32, error) {
	lm := 10
	db, err := s.DBGetter.GetDB()
	if err != nil {
		return nil, err
	}
	metas := make([]Meta, lm)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var res *gorm.DB
	err = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		res = db.Limit(lm).Find(&metas)
		return res.Error
	})
	if err != nil {
		return nil, err
	}
	ids := make([]uint32, res.RowsAffected)
	for i, meta := range metas {
		ids[i] = uint32(meta.ID)
	}
	return ids, nil
}

// Repository implementation for rdb
// Sqlite Repoitory
// embedded Default(expext to injecte sqlite conn)
type SqliteAnalyzeRepository struct {
	*DefaultAnalyzeRepository
}

type SqliteAnalyzeConnInfo string

// has conn info struct work with sqlite repository
type SqlteDBGetter struct {
	Info SqliteAnalyzeConnInfo
}

// impl for GetDB abstruct func
func (s *SqlteDBGetter) GetDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(string(s.Info)), &gorm.Config{})
}

func DefaultSqliteAnalyzeConnInfo() *SqliteAnalyzeConnInfo {
	var v SqliteAnalyzeConnInfo = "analyze.db"
	return &v
}

// constructor for sqlite
func NewSqliteAnalyzeRepository(info *SqliteAnalyzeConnInfo) *SqliteAnalyzeRepository {
	return &SqliteAnalyzeRepository{DefaultAnalyzeRepository: &DefaultAnalyzeRepository{DBGetter: &SqlteDBGetter{Info: *info}}}
}

// Mysql Repostitory
type MysqlAnalyzeRepository struct {
	*DefaultAnalyzeRepository
}

// mysql conn info for Repository
type MysqlConnInfo struct {
	User     string
	Password string
	Host     string
	Port     uint
	DbName   string
}
type MysqlDBGetter struct {
	Info *MysqlConnInfo
}

// impl for GetDB abstruct func
func (s *MysqlDBGetter) GetDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", s.Info.User, s.Info.Password, s.Info.Host, s.Info.Port, s.Info.DbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func DefaultMysqlAnalyzeConnInfo() *MysqlConnInfo {
	var host = "localhost"
	if cfg.MysqlHost != "" {
		host = cfg.MysqlHost
	}
	var user = "analyze"
	if cfg.MysqlUser != "" {
		user = cfg.MysqlUser
	}
	var password = "password"
	if cfg.MysqlPassword != "" {
		password = cfg.MysqlPassword
	}
	var dbname = "analyze"
	if cfg.MysqlDbName != "" {
		dbname = cfg.MysqlDbName
	}
	var port uint = 3306
	if cfg.MysqlPort > 0 {
		port = cfg.MysqlPort
	}

	return &MysqlConnInfo{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		DbName:   dbname,
	}
}

// constructor for mysql
func NewMysqlAnalyzeRepository(info *MysqlConnInfo) *MysqlAnalyzeRepository {
	return &MysqlAnalyzeRepository{DefaultAnalyzeRepository: &DefaultAnalyzeRepository{DBGetter: &MysqlDBGetter{Info: info}}}
}

var ErrNilInput = errors.New("nil in input")

// helper converters
func toPBMeta(meta *Meta) (*pb.Meta, error) {
	if meta == nil {
		return nil, ErrNilInput
	}
	var paramDef map[string]string
	if err := json.Unmarshal([]byte(meta.ParamDef), &paramDef); err != nil {
		return nil, err
	}
	return &pb.Meta{
		Id: uint32(meta.ID),
		MetaBody: &pb.MetaBody{
			Name:        meta.Name,
			Desc:        meta.Desc,
			ParamDef:    paramDef,
			IsAvailable: meta.IsAvailable,
		},
		CreateTime: timestamppb.New(meta.CreatedAt),
		UpdateTime: timestamppb.New(meta.UpdatedAt),
	}, nil
}

func fromPBMetaBody(metaBody *pb.MetaBody) (*Meta, error) {
	if metaBody == nil {
		return nil, ErrNilInput
	}
	paramDef, err := json.Marshal(metaBody.ParamDef)
	if err != nil {
		return nil, err
	}
	return &Meta{Name: metaBody.Name, Desc: metaBody.Desc, IsAvailable: metaBody.IsAvailable, ParamDef: string(paramDef)}, nil
}
