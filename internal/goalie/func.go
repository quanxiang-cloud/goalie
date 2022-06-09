package goalie

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/goalie/internal/models"
	"github.com/quanxiang-cloud/goalie/internal/models/mysql"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"gorm.io/gorm"
)

//Func interface
type Func interface {

	//List 当前系统功能
	List(ctx context.Context, req *GetListFuncRequest) (*GetListFuncResponse, error)

	// IN 获取指定功能集
	IN(ctx context.Context, req *InFuncRequest) (*InFuncResponse, error)
	// ListUserFuncTag 获取用户全量tag
	ListUserFuncTag(ctx context.Context, req *ListUserFuncTagRequest) (*ListUserFuncTagResponse, error)
	// ListFuncTag 获取全量tag
	ListFuncTag(ctx context.Context, req *ListFuncTagRequest) (*ListFuncTagResponse, error)
}

type funcs struct {
	db *gorm.DB

	funcRepo      models.FuncRepo
	roleOwnerRepo models.RoleOwnerRepo
	roleFuncRepo  models.RoleFuncRepo
	roleRepo      models.RoleRepo
}

//GetListFuncRequest req
type GetListFuncRequest struct {
}

//GetListFuncResponse response
type GetListFuncResponse struct {
	List []OneFuncResp `json:"list"`
}

//OneFuncResp func resp
type OneFuncResp struct {
	ID          string `json:"id"`
	PID         string `json:"pid"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	CreateTime  int64  `json:"createTime"`
	UpdateTime  int64  `json:"updateTime"`
}

func (f *funcs) List(ctx context.Context, req *GetListFuncRequest) (*GetListFuncResponse, error) {
	list, err := f.funcRepo.List(f.db)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		res := new(GetListFuncResponse)
		for k := range list {
			res.List = append(res.List, OneFuncResp{
				ID:          list[k].ID,
				Tag:         list[k].FuncTag,
				Name:        list[k].Name,
				Description: list[k].Describe,
				CreateTime:  list[k].CreatedAt,
				UpdateTime:  list[k].UpdatedAt,
			})
		}
		return res, nil
	}
	return nil, nil
}

//InFuncRequest req
type InFuncRequest struct {
}

//InFuncResponse response
type InFuncResponse struct {
}

func (f *funcs) IN(ctx context.Context, req *InFuncRequest) (*InFuncResponse, error) {
	fmt.Println("implement me")
	return nil, nil
}

// ListUserFuncTagRequest 获取用户tag[参数]
type ListUserFuncTagRequest struct {
	UserID       string   `json:"userID"`
	DepartmentID []string `json:"departmentID"`
}

// ListUserFuncTagResponse 获取用户tag[返回值]
type ListUserFuncTagResponse struct {
	Tag   []string `json:"tag"`
	Total int64    `json:"total"`
}

// ListUserFuncTag 获取用户tag
func (f *funcs) ListUserFuncTag(ctx context.Context, req *ListUserFuncTagRequest) (*ListUserFuncTagResponse, error) {
	// 获取用户角色
	userRoles, _, err := f.roleOwnerRepo.ListOwnerRole(ctx, f.db, req.UserID, req.DepartmentID...)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(userRoles))
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}
	//判断是否有超管角色
	roles, _, err := f.roleRepo.In(f.db, roleIDs...)
	if err != nil {
		return nil, err
	}
	var hasSuper = false
	for k := range roles {
		if roles[k].Tag == string(models.Super) {
			hasSuper = true
			break
		}
	}
	var fcs = make([]*models.Func, 0)
	if hasSuper {
		fs, err := f.funcRepo.List(f.db)
		if err != nil {
			return nil, err
		}
		fcs = append(fcs, fs...)
	} else {
		roleFuncs, err := f.roleFuncRepo.List(f.db, roleIDs...)
		if err != nil {
			return nil, err
		}
		funcIDs := make([]string, 0, len(roleFuncs))
		for _, roleFunc := range roleFuncs {
			funcIDs = append(funcIDs, roleFunc.FuncID)
		}

		fs, err := f.funcRepo.In(f.db, funcIDs...)
		if err != nil {
			return nil, err
		}
		fcs = append(fcs, fs...)

	}

	tags := make([]string, 0, len(fcs))
	for _, fc := range fcs {
		tags = append(tags, fc.FuncTag)
	}
	return &ListUserFuncTagResponse{
		Tag: tags,
	}, nil
}

// ListFuncTagRequest 获取全量tag[参数]
type ListFuncTagRequest struct{}

// ListFuncTagResponse 获取全量tag[返回]
type ListFuncTagResponse struct {
	Tag []string `json:"tag"`
}

// ListFuncTag 获取全量tag
func (f *funcs) ListFuncTag(ctx context.Context, req *ListFuncTagRequest) (*ListFuncTagResponse, error) {
	fcs, err := f.funcRepo.List(f.db)
	if err != nil {
		return nil, err
	}
	tags := make([]string, 0, len(fcs))
	for _, fc := range fcs {
		tags = append(tags, fc.FuncTag)
	}
	return &ListFuncTagResponse{
		Tag: tags,
	}, nil
}

// NewFunc 创建逻辑层
func NewFunc(conf *config.Config, db *gorm.DB) Func {
	return &funcs{
		db:            db,
		funcRepo:      mysql.NewFuncRepo(),
		roleFuncRepo:  mysql.NewRoleFuncRepo(),
		roleOwnerRepo: mysql.NewRoleOwnerRepo(),
		roleRepo:      mysql.NewRoleRepo(),
	}
}
