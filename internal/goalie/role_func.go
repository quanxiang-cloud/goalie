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

// RoleFunc  逻辑处理接口
type RoleFunc interface {
	//CreateInBatches 给角色赋予具体的功能
	CreateInBatches(ctx context.Context, req *UpdateRoleFuncRequest) (*UpdateRoleFuncResponse, error)
	ListRoleFunc(ctx context.Context, req *ListRoleFuncRequest) (*ListRoleFuncResponse, error)
}

type roleFunc struct {
	db           *gorm.DB
	roleRepo     models.RoleRepo
	roleFuncRepo models.RoleFuncRepo
	funcRepo     models.FuncRepo
}

// UpdateRoleFuncRequest 修改角色功能集[参数]
type UpdateRoleFuncRequest struct {
	RoleID string   `json:"roleID"`
	FuncID []string `json:"funcID"`
}

// UpdateRoleFuncResponse 修改角色功能集[返回]
type UpdateRoleFuncResponse struct {
}

func (r *roleFunc) CreateInBatches(ctx context.Context, req *UpdateRoleFuncRequest) (*UpdateRoleFuncResponse, error) {
	old, err := r.roleRepo.Get(r.db, req.RoleID)
	if old == nil {
		return nil, error2.New(code.RoleNotExist)
	}
	// 超级管理员 不允许修改
	if old.Tag == string(models.Super) {
		return nil, error2.New(code.ErrSuperUpdate)
	}

	tx := r.db.Begin()
	// 删除需要删除的func
	err = r.roleFuncRepo.DeleteByRoleID(ctx, tx, req.RoleID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// 添加新的func
	adds := make([]*models.RoleFunc, 0, len(req.FuncID))
	for _, funcID := range req.FuncID {
		adds = append(adds, &models.RoleFunc{
			ID:        id2.ShortID(0),
			RoleID:    req.RoleID,
			FuncID:    funcID,
			CreatedAt: time2.NowUnix(),
			UpdatedAt: time2.NowUnix(),
		})
	}
	err = r.roleFuncRepo.CreateInBatches(tx, adds...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	//TODO 清除用户的权限缓存
	return nil, nil
}

// ListRoleFuncRequest 获取角色功能集[参数]
type ListRoleFuncRequest struct {
	RoleID string `json:"roleID" form:"roleID"`
}

// ListRoleFuncResponse 获取角色功能集[返回值]
type ListRoleFuncResponse struct {
	Func         []ListRoleFuncVO `json:"func"`
	LastSaveTime int64            `json:"lastSaveTime"`
}

// ListRoleFuncVO 权限功能集VO
type ListRoleFuncVO struct {
	ID      string           `json:"id"`
	PID     string           `json:"pid"`
	Name    string           `json:"name"`
	FuncTag string           `json:"funcTag"`
	Has     bool             `json:"has"`
	Child   []ListRoleFuncVO `json:"child"`
}

func (r *roleFunc) ListRoleFunc(ctx context.Context, req *ListRoleFuncRequest) (*ListRoleFuncResponse, error) {
	// 获取全部func
	funcList, err := r.funcRepo.List(r.db)
	if err != nil {
		return nil, err
	}
	if len(funcList) == 0 {
		return nil, error2.New(code.ErrNoFuncs)
	}
	//获取角色
	roleOne, err := r.roleRepo.Get(r.db, req.RoleID)
	if err != nil {
		return nil, err
	}
	tmp := make([]ListRoleFuncVO, 0, len(funcList))
	var lastSaveTime int64
	if roleOne.Tag == string(models.Super) {
		for _, f := range funcList {
			vo := ListRoleFuncVO{
				ID:      f.ID,
				Name:    f.Name,
				FuncTag: f.FuncTag,
				PID:     f.PID,
			}
			if f.UpdatedAt > lastSaveTime {
				lastSaveTime = f.CreatedAt
			}

			vo.Has = true

			tmp = append(tmp, vo)
		}
	} else {
		// 获取角色的func
		roleFuncs, err := r.roleFuncRepo.List(r.db, req.RoleID)
		if err != nil {
			return nil, err
		}

		// 转换为map
		roleFuncMap := make(map[string]struct{}, len(roleFuncs))
		for _, rf := range roleFuncs {
			roleFuncMap[rf.FuncID] = struct{}{}
			if rf.UpdatedAt > lastSaveTime {
				lastSaveTime = rf.CreatedAt
			}
		}
		for _, f := range funcList {
			vo := ListRoleFuncVO{
				ID:      f.ID,
				Name:    f.Name,
				FuncTag: f.FuncTag,
				PID:     f.PID,
			}
			if _, ok := roleFuncMap[f.ID]; ok {
				vo.Has = true
			}
			tmp = append(tmp, vo)
		}

	}

	trees := r.makeTrees(tmp)

	return &ListRoleFuncResponse{
		Func:         trees,
		LastSaveTime: lastSaveTime,
	}, nil
}

/*
makeRoot 提取pid=""的节点
*/
func (r *roleFunc) makeTrees(list []ListRoleFuncVO) []ListRoleFuncVO {
	var outs = make([]ListRoleFuncVO, 0)
	var mps = make(map[string][]ListRoleFuncVO)
	for k, v := range list {
		if v.PID == "" {
			outs = append(outs, list[k])
		} else {
			mps[v.PID] = append(mps[v.PID], v)
		}

	}
	if len(outs) > 0 {
		for k := range outs {
			r.makeTree(&outs[k], mps)
		}

	}

	return outs
}

/*
makeTree 发现每个节点的子数据
*/
func (r *roleFunc) makeTree(dep *ListRoleFuncVO, mps map[string][]ListRoleFuncVO) {
	for k := range mps {
		if k == dep.ID {
			dep.Child = append(dep.Child, mps[k]...)
		}
	}
	if len(dep.Child) > 0 {
		for i := 0; i < len(dep.Child); i++ {
			r.makeTree(&dep.Child[i], mps)
		}
	}
}

// NewRoleFunc 创建逻辑层
func NewRoleFunc(conf *config.Config, db *gorm.DB) RoleFunc {
	return &roleFunc{
		db:           db,
		roleRepo:     mysql.NewRoleRepo(),
		roleFuncRepo: mysql.NewRoleFuncRepo(),
		funcRepo:     mysql.NewFuncRepo(),
	}
}
