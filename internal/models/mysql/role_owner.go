package mysql

import (
	"context"
	"github.com/quanxiang-cloud/goalie/internal/models"
	"gorm.io/gorm"
)

type roleOwnerRepo struct {
}

// NewRoleOwnerRepo new role owner repo
func NewRoleOwnerRepo() models.RoleOwnerRepo {
	return &roleOwnerRepo{}
}

func (r *roleOwnerRepo) TableName() string {
	return "role_owner"
}

// Search 全量获取角色下关联
func (r *roleOwnerRepo) Search(db *gorm.DB, roleID string, ownerType models.OwnerType, page, limit int) ([]*models.RoleOwner, int64, error) {
	ql := db.Table(r.TableName()).
		Where("role_id = ?", roleID)

	if ownerType != models.OwnerTypeNone {
		ql = ql.Where("type = ? ", ownerType)
	}

	var total int64
	roleOwners := make([]*models.RoleOwner, 0)

	_ = ql.Count(&total).
		Error

	if limit != 0 {
		ql = ql.Limit(limit).Offset((page - 1) * limit)
	}
	err := ql.
		Order("created_at ASC,id DESC").
		Find(&roleOwners).
		Error
	return roleOwners, total, err
}

// Delete 删除
func (r *roleOwnerRepo) Delete(db *gorm.DB, ids ...string) error {
	return db.
		Where("id in ?", ids).
		Delete(r).
		Error
}

// DeleteByOwnerID 根据拥有者id删除
func (r *roleOwnerRepo) DeleteByOwnerID(db *gorm.DB, _t models.OwnerType, ownerID string) error {
	return db.
		Where("type = ? and owner_id = ?", _t, ownerID).
		Delete(r).
		Error
}

// Create 添加角色用户
func (r *roleOwnerRepo) Create(db *gorm.DB, entity *models.RoleOwner) error {
	return db.Table(r.TableName()).
		Create(entity).
		Error
}

func (r *roleOwnerRepo) CreateInBatches(db *gorm.DB, entity ...*models.RoleOwner) error {
	return db.Table(r.TableName()).
		CreateInBatches(entity, len(entity)).
		Error
}

// TransferRoleSuper 转让超级管理员
func (r *roleOwnerRepo) TransferRoleSuper(db *gorm.DB, entity *models.RoleOwner) error {
	return db.Table(r.TableName()).
		Where("id = ?", entity.ID).
		Updates(map[string]interface{}{
			"owner_id":   entity.OwnerID,
			"created_at": entity.CreatedAt,
		}).
		Error
}

func (r *roleOwnerRepo) Get(db *gorm.DB, roleID, ownerID string) (*models.RoleOwner, error) {
	roleOwner := new(models.RoleOwner)
	err := db.Table(r.TableName()).
		Where("role_id = ? and owner_id = ?", roleID, ownerID).
		Find(roleOwner).
		Error
	if err != nil {
		return nil, err
	}
	if roleOwner.ID == "" {
		return nil, nil
	}
	return roleOwner, nil
}

// ListOwnerRole 获取用户角色
func (r *roleOwnerRepo) ListOwnerRole(ctx context.Context, db *gorm.DB, ownerID string, departmentID ...string) ([]*models.RoleOwner, int64, error) {
	ql := db.Table(r.TableName()).
		Where("owner_id = ? and type = ?", ownerID, models.Personnel)

	if len(departmentID) != 0 {
		ql.Or("owner_id in ?  and type = ?", departmentID, models.Department)
	}

	var total int64
	roleOwners := make([]*models.RoleOwner, 0)

	err := ql.Count(&total).
		Order("created_at ASC,id DESC").
		Find(&roleOwners).
		Error
	return roleOwners, total, err
}

// DeleteByRoleID 根据角色id删除
func (r *roleOwnerRepo) DeleteByRoleID(ctx context.Context, tx *gorm.DB, roleID ...string) error {
	return tx.WithContext(ctx).Where("role_id in (?)", roleID).Delete(r).Error
}

// DeleteByRoleIDAndOwnerID 根据角色id用户id删除
func (r *roleOwnerRepo) DeleteByRoleIDAndOwnerID(ctx context.Context, tx *gorm.DB, roleID string, ownerIDs ...string) error {
	return tx.WithContext(ctx).Where("role_id=? and owner_id in (?)", roleID, ownerIDs).Delete(r).Error
}

// DeleteByOwnerIDs 根据用户id删除
func (r *roleOwnerRepo) DeleteByOwnerIDs(ctx context.Context, tx *gorm.DB, ownerID ...string) error {
	return tx.WithContext(ctx).Where("owner_id in (?)", ownerID).Delete(r).Error
}
