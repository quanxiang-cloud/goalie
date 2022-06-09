package models

import (
	"context"
	"gorm.io/gorm"
)

// RoleFunc 角色功能集
type RoleFunc struct {
	ID     string
	RoleID string
	FuncID string

	CreatedAt int64
	UpdatedAt int64
}

// RoleFuncRepo 角色功能集[存储服务]
type RoleFuncRepo interface {
	// Create 角色添加功能
	Create(db *gorm.DB, roleFunc *RoleFunc) error

	// CreateInBatch 批量添加
	CreateInBatches(db *gorm.DB, roleFunc ...*RoleFunc) error

	// Delete 角色删除功能
	Delete(db *gorm.DB, ids ...string) error

	// Get 获取角色功能
	Get(db *gorm.DB, roleID, funcID string) (*RoleFunc, error)

	// List 获取角色功能集
	List(db *gorm.DB, roleID ...string) ([]*RoleFunc, error)

	// DeleteByRoleID 角色删除功能
	DeleteByRoleID(ctx context.Context, db *gorm.DB, ids ...string) error
}
