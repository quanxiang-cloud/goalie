package models

import (
	"context"
	"gorm.io/gorm"
)

// RoleTag 角色标签
type RoleTag string

const (
	// Super 超级管理员
	Super RoleTag = "super"
	// Warden 管理员
	Warden RoleTag = "warden"
	// Viewer 查看者
	Viewer RoleTag = "viewer"
)

// Role 角色
type Role struct {
	ID   string
	Name string
	// Tag
	// super 超级管理员
	// warden 管理员
	// viewer 查看者
	Tag string

	CreatedAt int64
	UpdatedAt int64

	CreatedBy string
	UpdatedBy string
	TenantID  string
}

// RoleRepo 角色[存储服务]
type RoleRepo interface {
	// Search 查询角色列表
	Search(db *gorm.DB, tenantID string) ([]*Role, int64, error)

	// Get 获取角色
	Get(db *gorm.DB, id string) (*Role, error)

	// In 获取指定角色
	In(db *gorm.DB, ids ...string) ([]*Role, int64, error)

	// GetWithTag tag获取角色
	GetWithTag(db *gorm.DB, tag RoleTag) (*Role, error)

	// Create 创建角色
	Create(ctx context.Context, tx *gorm.DB, role *Role) error

	// Update 更新角色
	Update(ctx context.Context, tx *gorm.DB, role *Role) error

	// Delete 删除角色
	Delete(ctx context.Context, tx *gorm.DB, id string) error

	// SearchByIDs 查询角色列表
	SearchByIDs(ctx context.Context, db *gorm.DB, ids []string) []Role

	// PageList 分页查询
	PageList(ctx context.Context, db *gorm.DB, tenantID, name string, page, limit int) ([]Role, int64)
}
