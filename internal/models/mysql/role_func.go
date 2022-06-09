package mysql

import (
	"context"
	"github.com/quanxiang-cloud/goalie/internal/models"
	"gorm.io/gorm"
)

type roleFuncRepo struct {
}

// NewRoleFuncRepo new role func repo
func NewRoleFuncRepo() models.RoleFuncRepo {
	return &roleFuncRepo{}
}

func (r *roleFuncRepo) TableName() string {
	return "role_func"
}

// Create 角色添加功能
func (r *roleFuncRepo) Create(db *gorm.DB, roleFunc *models.RoleFunc) error {
	return db.Table(r.TableName()).
		Create(roleFunc).
		Error
}

// CreateInBatch 批量添加
func (r *roleFuncRepo) CreateInBatches(db *gorm.DB, roleFunc ...*models.RoleFunc) error {
	return db.Table(r.TableName()).
		CreateInBatches(roleFunc, len(roleFunc)).
		Error
}

// Delete 角色删除功能
func (r *roleFuncRepo) Delete(db *gorm.DB, ids ...string) error {
	return db.
		Where("id in ?", ids).Delete(r).
		Error
}

// Get 获取角色功能
func (r *roleFuncRepo) Get(db *gorm.DB, roleID, funcID string) (*models.RoleFunc, error) {
	return nil, nil
}

// List 获取角色功能集
func (r *roleFuncRepo) List(db *gorm.DB, roleID ...string) ([]*models.RoleFunc, error) {
	roleFunc := make([]*models.RoleFunc, 0)
	err := db.Table(r.TableName()).
		Where("role_id in ?", roleID).
		Order("created_at DESC").
		Find(&roleFunc).
		Error
	return roleFunc, err
}

// DeleteByRoleID 角色删除功能
func (r *roleFuncRepo) DeleteByRoleID(ctx context.Context, tx *gorm.DB, ids ...string) error {
	return tx.WithContext(ctx).Where("role_id in (?)", ids).Delete(r).Error

}
