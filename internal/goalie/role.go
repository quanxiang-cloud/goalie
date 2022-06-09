package goalie

import (
	"context"
	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/goalie/internal/models"
	"github.com/quanxiang-cloud/goalie/internal/models/mysql"
	"github.com/quanxiang-cloud/goalie/pkg/code"
	"github.com/quanxiang-cloud/goalie/pkg/config"

	"gorm.io/gorm"
)

// Role  Service处理接口
type Role interface {

	// CreateRole 创建角色
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error)
	// UpdateRole 更新角色
	UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*UpdateRoleResponse, error)
	// PageListRole 分页批量获取角色
	PageListRole(ctx context.Context, req *PageListRoleRequest) (*PageListRoleResponse, error)

	// GetRole 获取单个角色
	GetRole(ctx context.Context, req *GetRoleRequest) (*GetRoleResponse, error)

	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error)
}

type role struct {
	db *gorm.DB

	roleRepo      models.RoleRepo
	roleOwnerRepo models.RoleOwnerRepo
	roleFuncRepo  models.RoleFuncRepo
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name      string `json:"name" binding:"required,max=30"`
	Tag       string `json:"tag"`
	CreatedBy string
	TenantID  string
}

// CreateRoleResponse 创建角色响应
type CreateRoleResponse struct {
	ID string `json:"id"`
}

func (r *role) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	response := new(CreateRoleResponse)
	var nowAt = time2.NowUnix()
	var role = new(models.Role)
	role.ID = id2.ShortID(0)
	role.Name = req.Name
	if req.Tag == "" {
		role.Tag = string(models.Viewer)
	} else {
		role.Tag = req.Tag
	}

	role.CreatedBy = req.CreatedBy
	role.UpdatedBy = req.CreatedBy
	role.CreatedAt = nowAt
	role.UpdatedAt = nowAt
	role.TenantID = req.TenantID
	tx := r.db.Begin()
	if err := r.roleRepo.Create(ctx, tx, role); err != nil {
		tx.Rollback()
		return nil, error2.New(code.ErrNameUsed)
	}
	tx.Commit()
	response.ID = role.ID
	return response, nil
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	ID        string `json:"id"  binding:"required"`
	Name      string `json:"name" binding:"max=30"`
	Tag       string `json:"tag"`
	UpdatedBy string
}

// UpdateRoleResponse 更新角色请求
type UpdateRoleResponse struct {
}

func (r *role) UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*UpdateRoleResponse, error) {

	var nowAt = time2.NowUnix()

	old, err := r.roleRepo.Get(r.db, req.ID)
	if err != nil || old == nil {
		return nil, error2.New(code.RoleNotExist)
	}

	old.Name = req.Name
	old.Tag = req.Tag
	if old.Tag == "" && req.Tag == "" {
		old.Tag = string(models.Viewer)
	}
	old.UpdatedBy = req.UpdatedBy
	old.UpdatedAt = nowAt
	tx := r.db.Begin()
	if err := r.roleRepo.Update(ctx, tx, old); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &UpdateRoleResponse{}, nil
}

// GetRoleRequest 获取角色[参数]
type GetRoleRequest struct {
	ID string `json:"id" form:"id"`
}

// GetRoleResponse 获取角色返回值[参数]
type GetRoleResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	CreatedAt int64  `json:"createdAt"`
}

// GetRole 获取角色信息
func (r *role) GetRole(ctx context.Context, req *GetRoleRequest) (*GetRoleResponse, error) {
	old, err := r.roleRepo.Get(r.db, req.ID)
	if err != nil {
		return nil, err
	}
	if old == nil {
		return nil, error2.New(code.RoleNotExist)
	}
	return &GetRoleResponse{
		ID:        old.ID,
		Name:      old.Name,
		Tag:       old.Tag,
		CreatedAt: old.CreatedAt,
	}, nil
}

// PageListRoleRequest 批量查询请求
type PageListRoleRequest struct {
	Name     string `json:"name" form:"name"`
	Page     int    `json:"page" form:"page" binding:"required,min=0"`
	Limit    int    `json:"limit" form:"limit" binding:"required"`
	TenantID string
}

// PageListRoleResponse 批量查询响应
type PageListRoleResponse struct {
	TotalCount int64             `json:"totalCount"`
	Data       []GetRoleResponse `json:"data"`
}

func (r *role) PageListRole(ctx context.Context, req *PageListRoleRequest) (*PageListRoleResponse, error) {

	search, total := r.roleRepo.PageList(ctx, r.db, req.TenantID, req.Name, req.Page, req.Limit)
	res := &PageListRoleResponse{}
	if len(search) > 0 {
		data := make([]GetRoleResponse, 0, len(search))
		for k := range search {
			data = append(data, GetRoleResponse{
				ID:        search[k].ID,
				Tag:       search[k].Tag,
				Name:      search[k].Name,
				CreatedAt: search[k].CreatedAt,
			})
		}
		res.TotalCount = total
		res.Data = data
	}

	return res, nil
}

// DeleteRoleRequest 删除角色请求
type DeleteRoleRequest struct {
	ID string `json:"id" form:"id" uri:"id" binding:"required"`
}

// DeleteRoleResponse 删除角色响应
type DeleteRoleResponse struct {
}

func (r *role) DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	role, err := r.roleRepo.Get(r.db, req.ID)
	if err != nil || role == nil {
		return nil, error2.New(code.RoleNotExist)
	}
	tx := r.db.Begin()
	if err := r.roleRepo.Delete(ctx, tx, req.ID); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := r.roleFuncRepo.DeleteByRoleID(ctx, tx, req.ID); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := r.roleOwnerRepo.DeleteByRoleID(ctx, tx, req.ID); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &DeleteRoleResponse{}, nil
}

// NewRole 创建Service
func NewRole(conf *config.Config, db *gorm.DB) Role {
	return &role{
		db:            db,
		roleRepo:      mysql.NewRoleRepo(),
		roleFuncRepo:  mysql.NewRoleFuncRepo(),
		roleOwnerRepo: mysql.NewRoleOwnerRepo(),
	}
}
