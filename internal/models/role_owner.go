package models

import (
	"context"
	"gorm.io/gorm"
)

// OwnerType 拥有者类型
type OwnerType int

const (
	// OwnerTypeNone none
	OwnerTypeNone OwnerType = 0
	// Personnel 人员
	Personnel OwnerType = 1
	// Department 部门
	Department OwnerType = 2
)

const (
	// MaxSuperOwner 最多超级管理员人数
	MaxSuperOwner = 1
)

// RoleOwner 角色关联
type RoleOwner struct {
	ID     string
	RoleID string
	// Type 1: 人员 2:部门
	Type    OwnerType
	OwnerID string

	CreatedAt int64
}

// RoleOwnerRepo 角色关联[存储服务]
type RoleOwnerRepo interface {
	// Search 全量获取角色下关联
	Search(db *gorm.DB, roleID string, ownerType OwnerType, page, lmit int) ([]*RoleOwner, int64, error)

	// Delete 删除
	Delete(db *gorm.DB, ids ...string) error

	// DeleteByOwnerID 根据拥有者id删除
	DeleteByOwnerID(db *gorm.DB, _t OwnerType, ownerID string) error

	// Create 添加角色用户
	Create(db *gorm.DB, entity *RoleOwner) error

	// CreateInBatches 批量添加角色用户
	CreateInBatches(db *gorm.DB, entity ...*RoleOwner) error

	// TransferRoleSuper 转让超级管理员
	TransferRoleSuper(db *gorm.DB, entity *RoleOwner) error

	// Get 获取指定角色关联
	Get(db *gorm.DB, roleID, ownerID string) (*RoleOwner, error)

	// ListOwnerRole 获取拥有者角色
	ListOwnerRole(ctx context.Context, db *gorm.DB, ownerID string, departmentID ...string) ([]*RoleOwner, int64, error)
	// DeleteByRoleID 根据角色id删除
	DeleteByRoleID(ctx context.Context, db *gorm.DB, roleIDs ...string) error
	// DeleteByRoleIDAndOwnerID 根据角色id用户id删除
	DeleteByRoleIDAndOwnerID(ctx context.Context, db *gorm.DB, roleID string, ownerIDs ...string) error
	// DeleteByOwnerIDs 根据用户ids删除
	DeleteByOwnerIDs(ctx context.Context, db *gorm.DB, ownerIDs ...string) error
}
