package service

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "github.com/ynishi/gdean/pb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockRepository struct {
	*DefaultAnalyzeRepository
}

type MockDBGetter struct {
	Mock sqlmock.Sqlmock
	Db   *gorm.DB
}

// constructor for mock
func NewMockRepository() (*MockRepository, sqlmock.Sqlmock) {
	getter := MockDBGetter{}
	getter.GetDB()
	return &MockRepository{DefaultAnalyzeRepository: &DefaultAnalyzeRepository{DBGetter: &getter}}, getter.Mock
}

func (s *MockDBGetter) GetDB() (*gorm.DB, error) {
	if s.Db != nil {
		return s.Db, nil
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	gdb, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}}, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s.Db = gdb
	s.Mock = mock
	return s.Db, nil
}

// helper based on https://github.com/data-dog/go-sqlmock/issues/262
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var anyTime = AnyTime{}

func TestShouldFetch(t *testing.T) {
	repo, mock := NewMockRepository()
	q := regexp.QuoteMeta("SELECT * FROM `meta` WHERE `meta`.`id` = ? AND `meta`.`deleted_at` IS NULL ORDER BY `meta`.`id` LIMIT 1")
	rows := sqlmock.NewRows([]string{"id", "name", "desc", "is_available", "param_def"}).AddRow(1, "name1", "desc of name1", true, "null")
	mock.ExpectQuery(q).WithArgs(1).WillReturnRows(rows)
	meta := pb.MetaBody{Name: "name1", Desc: "desc of name1", IsAvailable: true}
	m, err := repo.Fetch(1)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	if err == nil {
		assert.Equal(t, m.Id, uint32(1), "id should be 1")
		assert.Equal(t, m.MetaBody, &meta, "metabody should be set")
	}
}

func TestShouldCreate(t *testing.T) {
	repo, mock := NewMockRepository()
	meta := pb.MetaBody{Name: "name1", Desc: "desc of name1", IsAvailable: true}
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `meta`*").WithArgs(anyTime, anyTime, nil, "name1", "desc of name1", true, "null").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	m, err := repo.Create(&meta)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	if err == nil {
		assert.Equal(t, m.Id, uint32(1), "id should be 1")
		assert.Equal(t, m.MetaBody, &pb.MetaBody{Name: "name1", Desc: "desc of name1", IsAvailable: true, ParamDef: nil}, "metabody should be set")
	}
}

func TestShouldPut(t *testing.T) {
	repo, mock := NewMockRepository()
	q := regexp.QuoteMeta("SELECT * FROM `meta` WHERE `meta`.`id` = ? AND `meta`.`deleted_at` IS NULL ORDER BY `meta`.`id` LIMIT 1")
	rows := sqlmock.NewRows([]string{"id", "name", "desc", "is_available", "param_def"}).AddRow(1, "name1", "desc of name", true, "null")
	mock.ExpectQuery(q).WithArgs(1).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `meta`*").WithArgs(anyTime, anyTime, nil, "name2", "desc of name2", false, `{"x":"float"}`, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	meta := pb.MetaBody{Name: "name2", Desc: "desc of name2", IsAvailable: false, ParamDef: map[string]string{"x": "float"}}
	m, err := repo.Put(1, &meta)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	if err == nil {
		assert.Equal(t, m.Id, uint32(1), "id should be 1")
	}
}

func TestShouldDelete(t *testing.T) {
	repo, mock := NewMockRepository()
	q := regexp.QuoteMeta("SELECT * FROM `meta` WHERE `meta`.`id` = ? AND `meta`.`deleted_at` IS NULL ORDER BY `meta`.`id` LIMIT 1")
	rows := sqlmock.NewRows([]string{"id", "name", "desc", "is_available", "param_def"}).AddRow(1, "name1", "desc of name", true, "null")
	mock.ExpectQuery(q).WithArgs(1).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `meta` SET `deleted_at`=*").WithArgs(anyTime, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(q).WithArgs(1).WillReturnError(sql.ErrNoRows)
	_, err := repo.Delete(1)
	assert.Nil(t, err)
	fetched, err := repo.Fetch(1)
	assert.EqualError(t, sql.ErrNoRows, "sql: no rows in result set")
	assert.Nil(t, fetched)
}
