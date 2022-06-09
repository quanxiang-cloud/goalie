package mysql

import (
	"github.com/quanxiang-cloud/goalie/internal/models"
	"gorm.io/gorm"
)

type funcRepo struct {
}

// NewFuncRepo new func repo
func NewFuncRepo() models.FuncRepo {
	return &funcRepo{}
}

func (f *funcRepo) TableName() string {
	return "func"
}

// List 获取功能集列表
func (f *funcRepo) List(db *gorm.DB) ([]*models.Func, error) {
	funcs := make([]*models.Func, 0)
	err := db.Table(f.TableName()).
		Order("created_at ASC").
		Find(&funcs).
		Error
	return funcs, err
}

// In 获取指定功能集
func (f *funcRepo) In(db *gorm.DB, ids ...string) ([]*models.Func, error) {
	funcs := make([]*models.Func, 0)
	err := db.Table(f.TableName()).
		Where("id in ?", ids).
		Find(&funcs).
		Error
	return funcs, err
}
