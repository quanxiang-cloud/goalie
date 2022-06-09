package goalie

import (
	"context"

	error2 "github.com/quanxiang-cloud/cabin/error"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/goalie/internal/models"
	"github.com/quanxiang-cloud/goalie/internal/models/mysql"
	"github.com/quanxiang-cloud/goalie/pkg/code"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"github.com/quanxiang-cloud/goalie/pkg/org"

	"gorm.io/gorm"
)

// Goalie 权限管理
type Goalie interface {
	// TransferRoleSuper 转让超级管理员
	TransferRoleSuper(ctx context.Context, req *TransferRoleSuperReq) (*TransferRoleSuperResp, error)
}

type goalie struct {
	db    *gorm.DB
	redis *redis.ClusterClient

	user org.User

	roleRepo      models.RoleRepo
	roleOwnerRepo models.RoleOwnerRepo
	funcRepo      models.FuncRepo
	roleFuncRepo  models.RoleFuncRepo
}

// NewGoalie 创建一个goalie
func NewGoalie(conf *config.Config, db *gorm.DB, opts ...Option) (Goalie, error) {

	g := &goalie{
		db: db,

		roleRepo:      mysql.NewRoleRepo(),
		roleOwnerRepo: mysql.NewRoleOwnerRepo(),
		funcRepo:      mysql.NewFuncRepo(),
		roleFuncRepo:  mysql.NewRoleFuncRepo(),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g, nil
}

// Option option
type Option func(g *goalie)

// WithUser org with user
func WithUser(user org.User) Option {
	return func(g *goalie) {
		g.user = user
	}
}

// TransferRoleSuperReq 转让超级管理员[参数]
type TransferRoleSuperReq struct {
	UserID     string
	Transferee string
}

// TransferRoleSuperResp 转让超级管理员[返回值]
type TransferRoleSuperResp struct {
}

// TransferRoleSuper 转让超级管理员
func (g *goalie) TransferRoleSuper(ctx context.Context, req *TransferRoleSuperReq) (*TransferRoleSuperResp, error) {
	tx := g.db.Begin()

	// 获取角色
	role, err := g.roleRepo.GetWithTag(tx, models.Super)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if role == nil {
		tx.Rollback()
		return nil, error2.New(code.UnknownRole)
	}

	// 必须是超级管理员才能转让超级管理员
	roleOwner, err := g.roleOwnerRepo.Get(tx, role.ID, req.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if roleOwner == nil {
		tx.Rollback()
		return nil, error2.New(code.InvalidRoleOwner)
	}

	err = g.roleOwnerRepo.TransferRoleSuper(tx, &models.RoleOwner{
		ID:        roleOwner.ID,
		OwnerID:   req.Transferee,
		CreatedAt: time2.NowUnix(),
	})

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil
}
