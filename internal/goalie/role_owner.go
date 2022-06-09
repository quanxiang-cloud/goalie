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
	"github.com/quanxiang-cloud/goalie/pkg/org"

	"net/http"

	"gorm.io/gorm"
)

//RoleOwner interface
type RoleOwner interface {
	// ListRoleOwner 获取角色的拥有者列表
	ListRoleOwner(ctx context.Context, header http.Header, req *ListRoleOwnerRequest) (*ListRoleOwnerResponse, error)
	// UpdateRoleOwner 修改角色拥有者
	UpdateRoleOwner(ctx context.Context, req *UpdateRoleOwnerRequest) (*UpdateRoleOwnerResponse, error)
	// ListUserRole 用户token获取用户角色列表
	ListUserRole(ctx context.Context, req *ListUserRoleRequest) (*ListUserRoleResponse, error)
	// DeleteOwnerRole 删除拥有者角色
	DeleteOwnerRole(ctx context.Context, req *DeleteOwnerRoleRequest) (*DeleteOwnerRoleResponse, error)
	// UpdateOwnerRole 修改拥有者角色
	UpdateOwnerRole(ctx context.Context, req *UpdateOwnerRoleRequest) (*UpdateOwnerRoleResponse, error)
	// ListOwnerRole 获取拥有者角色列表
	ListOwnerRole(ctx context.Context, req *ListOwnerRoleRequest) (*ListOwnerRoleResponse, error)
	// OthUserRoleList 其它获取用户角色列表
	OthUserRoleList(ctx context.Context, req *OthOneUserRoleListRequest) (*OthOneUserRoleListResponse, error)
	// OthRoleListByIDs 其它根据ids获取用户角色列表
	OthRoleListByIDs(ctx context.Context, req *OthRoleListByIDsRequest) (*OthRoleListByIDsResponse, error)
	// OthDelUserRoleList 其它删除用户角色列表
	OthDelUserRoleList(ctx context.Context, req *OthDelRequest) (*OthDelResponse, error)
}

const (
	systemAdminOwnerID = "1"
)

type roleOwner struct {
	db            *gorm.DB
	user          org.User
	roleRepo      models.RoleRepo
	roleOwnerRepo models.RoleOwnerRepo
}

// NewRoleOwner 创建逻辑层
func NewRoleOwner(conf *config.Config, db *gorm.DB) RoleOwner {
	return &roleOwner{
		db:            db,
		user:          org.NewUser(conf.InternalNet),
		roleOwnerRepo: mysql.NewRoleOwnerRepo(),
		roleRepo:      mysql.NewRoleRepo(),
	}
}

// ListRoleOwnerRequest 获取角色的拥有者[参数]
type ListRoleOwnerRequest struct {
	RoleID string `json:"roleID" form:"roleID"`
	Type   int    `json:"type" form:"type"`
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
}

// ListRoleOwnerResponse 获取角色的拥有者[返回值]
type ListRoleOwnerResponse struct {
	Owners []*ListRoleOwnerVO `json:"owners"`
	Total  int64              `json:"total"`
}

// ListRoleOwnerVO 角色拥有者列表VO
type ListRoleOwnerVO struct {
	ID          string             `json:"id"`
	Type        int                `json:"type"`
	OwnerID     string             `json:"ownerID"`
	OwnerName   string             `json:"ownerName"`
	Phone       string             `json:"phone,omitempty"`
	Email       string             `json:"email,omitempty"`
	Deps        []ListRoleOwnerDep `json:"deps,omitempty"`
	CreatedTime int64              `json:"createdTime,omitempty"`
	PID         string             `json:"pid,omitempty"`
}

//ListRoleOwnerDep req
type ListRoleOwnerDep struct {
	DepartmentID   string `json:"departmentID,omitempty"`
	DepartmentName string `json:"departmentName,omitempty"`
	PID            string `json:"pid"`
}

func (r *roleOwner) ListRoleOwner(ctx context.Context, header http.Header, req *ListRoleOwnerRequest) (*ListRoleOwnerResponse, error) {
	roleOwners, total, err := r.roleOwnerRepo.Search(r.db,
		req.RoleID,
		models.OwnerType(req.Type),
		req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	if len(roleOwners) == 0 {
		return nil, nil
	}
	// 获取用户和部门信息
	userIDs := make([]string, 0, len(roleOwners))
	departmentIDs := make([]string, 0, len(roleOwners))
	for _, elem := range roleOwners {
		if elem.Type == models.Personnel {
			userIDs = append(userIDs, elem.OwnerID)
		} else {
			departmentIDs = append(departmentIDs, elem.OwnerID)
		}
	}
	userMap := make(map[string]*org.OneUserResponse)
	departmentMap := make(map[string]*org.DepOneResponse)
	if len(userIDs) > 0 {
		reqUser := &org.GetUserByIDsRequest{
			IDs: userIDs,
		}
		userInfos, err := r.user.GetUserByIDs(ctx, reqUser)
		if err != nil {
			return nil, err
		}

		for k := range userInfos.Users {
			userMap[userInfos.Users[k].ID] = &userInfos.Users[k]
		}

	}
	if len(departmentIDs) > 0 {
		departments, err := r.user.GetDepByIDs(ctx, &org.GetDepByIDsRequest{IDs: departmentIDs})
		if err != nil {
			return nil, err
		}

		for k := range departments.Deps {
			departmentMap[departments.Deps[k].ID] = &departments.Deps[k]
		}
	}

	// 构建返回值
	roleOwnerList := make([]*ListRoleOwnerVO, 0, len(roleOwners))
	for _, elem := range roleOwners {
		vo := &ListRoleOwnerVO{
			ID:          elem.ID,
			Type:        int(elem.Type),
			OwnerID:     elem.OwnerID,
			CreatedTime: elem.CreatedAt,
		}
		if vo.Type == int(models.Personnel) {
			if v, ok := userMap[elem.OwnerID]; ok && v != nil {
				vo.OwnerName = userMap[elem.OwnerID].Name
				vo.Email = userMap[elem.OwnerID].Email
				if len(userMap[elem.OwnerID].Dep) > 0 {
					dep := make([]ListRoleOwnerDep, 0, len(userMap[elem.OwnerID].Dep))
					for k := range userMap[elem.OwnerID].Dep {
						for k1 := range userMap[elem.OwnerID].Dep[k] {
							dep = append(dep, ListRoleOwnerDep{
								DepartmentID:   userMap[elem.OwnerID].Dep[k][k1].ID,
								DepartmentName: userMap[elem.OwnerID].Dep[k][k1].Name,
								PID:            userMap[elem.OwnerID].Dep[k][k1].PID,
							})
						}

					}
					vo.Deps = append(vo.Deps, dep...)
				}
			}
		} else {
			if v, ok := departmentMap[vo.OwnerID]; ok && v != nil {
				vo.OwnerName = departmentMap[vo.OwnerID].Name
				vo.PID = departmentMap[vo.OwnerID].PID
			} else {
				total = total - 1
				continue
			}
		}
		roleOwnerList = append(roleOwnerList, vo)

	}

	return &ListRoleOwnerResponse{
		Owners: roleOwnerList,
		Total:  total,
	}, nil
}

// UpdateRoleOwnerRequest 修改角色拥有者[参数]
type UpdateRoleOwnerRequest struct {
	RoleID string   `json:"roleID"`
	Delete []string `json:"delete"`
	Add    []struct {
		Type    int    `json:"type"`
		OwnerID string `json:"ownerID"`
	} `json:"add"`
}

// UpdateRoleOwnerResponse 修改角色拥有者[返回值]
type UpdateRoleOwnerResponse struct{}

func (r *roleOwner) UpdateRoleOwner(ctx context.Context, req *UpdateRoleOwnerRequest) (*UpdateRoleOwnerResponse, error) {
	// 超级管理员不允许修改
	ro, err := r.roleRepo.Get(r.db, req.RoleID)
	if err != nil {
		return nil, err
	}
	if ro == nil {
		return nil, error2.New(code.RoleNotExist)
	}
	if ro.Tag == string(models.Super) {
		return nil, error2.New(code.ErrSuperUpdate)
	}
	tx := r.db.Begin()
	// 删除
	if len(req.Delete) > 0 {
		dels := make([]string, 0, len(req.Delete))
		for k := range req.Delete {
			if ro.Tag == string(models.Super) && req.Delete[k] == systemAdminOwnerID {
				continue
			}
			dels = append(dels, req.Delete[k])
		}
		err := r.roleOwnerRepo.DeleteByRoleIDAndOwnerID(ctx, tx, req.RoleID, dels...)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 添加
	roleOwners := make([]*models.RoleOwner, 0, len(req.Add))
	for _, elem := range req.Add {
		roleOwners = append(roleOwners, &models.RoleOwner{
			ID:        id2.ShortID(0),
			RoleID:    req.RoleID,
			OwnerID:   elem.OwnerID,
			Type:      models.OwnerType(elem.Type),
			CreatedAt: time2.NowUnix(),
		})

	}

	err = r.roleOwnerRepo.CreateInBatches(tx, roleOwners...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil
}

// ListUserRoleRequest 获取用户角色列表[参数]
type ListUserRoleRequest struct {
	UserID       string   `json:"userID"`
	DepartmentID []string `json:"departmentID"`
}

// ListUserRoleResponse 获取用户角色列表[返回值]
type ListUserRoleResponse struct {
	Roles []*ListRoleVO `json:"roles"`
	Total int64         `json:"total"`
}

// ListRoleVO 角色列表VO
type ListRoleVO struct {
	ID     string `json:"id"`
	RoleID string `json:"roleID"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
}

func (r *roleOwner) ListUserRole(ctx context.Context, req *ListUserRoleRequest) (*ListUserRoleResponse, error) {
	// 获取用户角色
	userRole, _, err := r.roleOwnerRepo.ListOwnerRole(ctx, r.db, req.UserID, req.DepartmentID...)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(userRole))
	ownerRoleMap := make(map[string]*models.RoleOwner, len(userRole))

	for _, elem := range userRole {
		roleIDs = append(roleIDs, elem.RoleID)
		ownerRoleMap[elem.RoleID] = elem
	}

	roles := r.roleRepo.SearchByIDs(ctx, r.db, roleIDs)
	if err != nil {
		return nil, err
	}

	listRoleVO := make([]*ListRoleVO, 0, len(roles))
	for _, elem := range roles {
		listRoleVO = append(listRoleVO, &ListRoleVO{
			ID:     ownerRoleMap[elem.ID].ID,
			RoleID: elem.ID,
			Tag:    elem.Tag,
			Name:   elem.Name,
		})
	}

	return &ListUserRoleResponse{
		Roles: listRoleVO,
		Total: int64(len(roles)),
	}, nil
}

// DeleteOwnerRoleRequest 删除拥有者关联[参数]
type DeleteOwnerRoleRequest struct {
	Type    int    `json:"type"`
	OwnerID string `json:"ownerID"`
}

// DeleteOwnerRoleResponse 删除拥有者关联[返回值]
type DeleteOwnerRoleResponse struct {
}

func (r *roleOwner) DeleteOwnerRole(ctx context.Context, req *DeleteOwnerRoleRequest) (*DeleteOwnerRoleResponse, error) {
	err := r.roleOwnerRepo.DeleteByOwnerID(r.db, models.OwnerType(req.Type), req.OwnerID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateOwnerRoleRequest 修改拥有者角色[参数]
type UpdateOwnerRoleRequest struct {
	Type    int      `json:"type"`
	OwnerID string   `json:"ownerID"`
	Delete  []string `json:"delete"`
	Add     []string `json:"add"`
}

// UpdateOwnerRoleResponse 修改拥有者角色[返回值]
type UpdateOwnerRoleResponse struct{}

func (r *roleOwner) UpdateOwnerRole(ctx context.Context, req *UpdateOwnerRoleRequest) (*UpdateOwnerRoleResponse, error) {
	// 超级管理员不允许修改
	role, err := r.roleRepo.GetWithTag(r.db, models.Super)
	if err != nil {

		return nil, err
	}

	for _, roleID := range req.Delete {
		if role.ID == roleID {
			return nil, error2.New(code.ErrSuperUpdate)
		}
	}
	tx := r.db.Begin()
	// 删除

	err = r.roleOwnerRepo.DeleteByOwnerID(tx, models.OwnerType(req.Type), req.OwnerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加
	roleOwners := make([]*models.RoleOwner, 0, len(req.Add))
	for _, roleID := range req.Add {
		if role.ID == roleID {
			tx.Rollback()
			return nil, error2.New(code.ErrSuperUpdate)
		}
		roleOwners = append(roleOwners, &models.RoleOwner{
			ID:        id2.ShortID(0),
			RoleID:    roleID,
			OwnerID:   req.OwnerID,
			Type:      models.OwnerType(req.Type),
			CreatedAt: time2.NowUnix(),
		})

	}

	err = r.roleOwnerRepo.CreateInBatches(tx, roleOwners...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil
}

// ListOwnerRoleRequest 获取拥有者角色列表
type ListOwnerRoleRequest struct {
	Type    int    `json:"type" form:"type"`
	OwnerID string `json:"ownerID" form:"ownerID"`
}

// ListOwnerRoleResponse 获取拥有者角色列表[返回值]
type ListOwnerRoleResponse struct {
	Roles []*ListRoleVO `json:"roles"`
	Total int64         `json:"total"`
}

func (r *roleOwner) ListOwnerRole(ctx context.Context, req *ListOwnerRoleRequest) (*ListOwnerRoleResponse, error) {
	ownerRoles, _, err := r.roleOwnerRepo.ListOwnerRole(ctx, r.db, req.OwnerID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(ownerRoles))
	ownerRoleMap := make(map[string]*models.RoleOwner, len(ownerRoles))
	for _, elem := range ownerRoles {
		roleIDs = append(roleIDs, elem.RoleID)
		ownerRoleMap[elem.RoleID] = elem
	}

	roles := r.roleRepo.SearchByIDs(ctx, r.db, roleIDs)
	if err != nil {
		return nil, err
	}

	listRoleVO := make([]*ListRoleVO, 0, len(roles))
	for _, elem := range roles {
		listRoleVO = append(listRoleVO, &ListRoleVO{
			ID:     ownerRoleMap[elem.ID].ID,
			RoleID: elem.ID,
			Tag:    elem.Tag,
			Name:   elem.Name,
		})
	}

	return &ListOwnerRoleResponse{
		Roles: listRoleVO,
		Total: int64(len(roles)),
	}, nil
}

//OthOneUserRoleListRequest 其它服务请求一个用户的角色集合
type OthOneUserRoleListRequest struct {
	UserID string   `json:"userID"`
	DepIDs []string `json:"depIDs"`
}

//OthOneUserRoleListResponse 其它服务查询一个人的角色返回信息
type OthOneUserRoleListResponse struct {
	RoleIDs []string `json:"roleIDs"`
}

func (r *roleOwner) OthUserRoleList(ctx context.Context, req *OthOneUserRoleListRequest) (*OthOneUserRoleListResponse, error) {
	// 获取用户角色
	userRole, _, err := r.roleOwnerRepo.ListOwnerRole(ctx, r.db, req.UserID, req.DepIDs...)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(userRole))

	for _, elem := range userRole {
		roleIDs = append(roleIDs, elem.RoleID)
	}

	return &OthOneUserRoleListResponse{
		RoleIDs: roleIDs,
	}, nil
}

//OthRoleListByIDsRequest 其它服务根据ids角色集合
type OthRoleListByIDsRequest struct {
	RoleIDs []string `json:"roleIDs"`
}

//OthRoleListByIDsResponse 其它服务根据IDs的角色返回信息
type OthRoleListByIDsResponse struct {
	Roles []ListRoleVO `json:"roles"`
}

//OthRoleListByIDs 其它根据ids获取用户角色列表
func (r *roleOwner) OthRoleListByIDs(ctx context.Context, req *OthRoleListByIDsRequest) (*OthRoleListByIDsResponse, error) {
	roles := r.roleRepo.SearchByIDs(ctx, r.db, req.RoleIDs)
	responses := make([]ListRoleVO, 0, len(roles))
	for k := range roles {
		ro := ListRoleVO{}
		ro.ID = roles[k].ID
		ro.Name = roles[k].Name
		ro.Tag = roles[k].Tag
		responses = append(responses, ro)
	}

	return &OthRoleListByIDsResponse{
		Roles: responses,
	}, nil
}

// OthDelRequest 其它服务删除权限
type OthDelRequest struct {
	IDs   []string `json:"ids"`
	DelBy string   `json:"delBy"`
}

// OthDelResponse 其它服务删除权限
type OthDelResponse struct {
}

//OthDelUserRoleList 不删除系统设定的管理员
func (r *roleOwner) OthDelUserRoleList(ctx context.Context, req *OthDelRequest) (*OthDelResponse, error) {
	ids := make([]string, 0, len(req.IDs))
	for k := range req.IDs {
		if req.IDs[k] == systemAdminOwnerID {
			continue
		}
		ids = append(ids, req.IDs[k])
	}
	err := r.roleOwnerRepo.DeleteByOwnerIDs(ctx, r.db, ids...)
	return &OthDelResponse{}, err
}
