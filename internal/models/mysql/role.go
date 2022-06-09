package mysql

import (
	"context"
	"github.com/quanxiang-cloud/goalie/internal/models"
	page2 "github.com/quanxiang-cloud/goalie/pkg/page"
	"gorm.io/gorm"
)

type roleRepo struct {
}

// NewRoleRepo new role repo
func NewRoleRepo() models.RoleRepo {
	return &roleRepo{}
}

func (r *roleRepo) TableName() string {
	return "role"
}

// Search 查询角色列表
func (r *roleRepo) Search(db *gorm.DB, tenantID string) ([]*models.Role, int64, error) {
	roles := make([]*models.Role, 0)
	var total int64
	db = db.Table(r.TableName())
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}

	err := db.Count(&total).
		Find(&roles).
		Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Get 获取角色
func (r *roleRepo) Get(db *gorm.DB, id string) (*models.Role, error) {
	role := new(models.Role)
	err := db.Table(r.TableName()).
		Where("id = ?", id).
		Find(role).
		Error
	if err != nil {
		return nil, err
	}
	if role.ID == "" {
		return nil, nil
	}
	return role, nil
}

// In 获取指定角色
func (r *roleRepo) In(db *gorm.DB, ids ...string) ([]*models.Role, int64, error) {
	roles := make([]*models.Role, 0, len(ids))
	var total int64

	err := db.Table(r.TableName()).
		Where("id in ?", ids).
		Count(&total).
		Find(&roles).
		Error
	return roles, total, err
}

// GetWithTag tag获取角色
func (r *roleRepo) GetWithTag(db *gorm.DB, tag models.RoleTag) (*models.Role, error) {
	role := new(models.Role)
	err := db.Table(r.TableName()).
		Where("tag = ?", tag).
		Find(role).
		Error
	if err != nil {
		return nil, err
	}
	if role.ID == "" {
		return nil, nil
	}
	return role, nil
}

// Delete 删除
func (r *roleRepo) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	return tx.WithContext(ctx).Where("id=?", id).Delete(r).Error
}

// Update 更新
func (r *roleRepo) Update(ctx context.Context, tx *gorm.DB, role *models.Role) error {
	return tx.WithContext(ctx).Table(r.TableName()).Updates(role).Error
}

// Create 创建
func (r *roleRepo) Create(ctx context.Context, tx *gorm.DB, role *models.Role) error {
	return tx.WithContext(ctx).Table(r.TableName()).Create(role).Error
}

// SearchByIDs 根据id批量查询
func (r *roleRepo) SearchByIDs(ctx context.Context, db *gorm.DB, ids []string) []models.Role {
	roles := make([]models.Role, 0)
	affected := db.WithContext(ctx).
		Table(r.TableName()).
		Where("id in (?)", ids).
		Find(&roles).RowsAffected
	if affected > 0 {
		return roles
	}
	return roles
}

// PageList 分页查询
func (r *roleRepo) PageList(ctx context.Context, db *gorm.DB, tenantID, name string, page, limit int) ([]models.Role, int64) {
	db = db.WithContext(ctx).Table(r.TableName())
	if tenantID == "" {
		db = db.Where("tenant_id=? or tenant_id is null", tenantID)
	} else {
		db = db.Where("tenant_id=?", tenantID)
	}

	if name != "" {
		db = db.Where("name=?", "%"+name+"%")
	}
	db = db.Order("updated_at desc")
	roles := make([]models.Role, 0)
	var num int64
	db.Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex)

	affected := db.Find(&roles).RowsAffected
	if affected > 0 {
		return roles, num
	}
	return nil, 0

}
