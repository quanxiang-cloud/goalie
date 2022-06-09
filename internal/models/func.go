package models

import "gorm.io/gorm"

// Func 功能集合
type Func struct {
	ID string
	// PID 父级ID 用于构建多叉树
	PID string
	// FuncTag 预留给前端
	FuncTag  string
	Name     string
	Describe string

	CreatedAt int64
	UpdatedAt int64
}

// FuncRepo 功能集[存储服务]
type FuncRepo interface {
	// Create 创建功能集
	// Create(c *redis.ClusterClient, f *Func) error

	// Get 获取功能集
	// Get(c *redis.ClusterClient, pid string) (*Func, error)

	// Delete 删除功能集合
	// Delete(c *redis.ClusterClient, ids ...string) error

	// List 获取功能集列表
	List(db *gorm.DB) ([]*Func, error)

	// In 获取指定功能集
	In(db *gorm.DB, ids ...string) ([]*Func, error)
}
