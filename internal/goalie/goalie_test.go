package goalie

//import (
//	"context"
//	"testing"
//
//	"github.com/quanxiang-cloud/goalie/internal/models"
//	"github.com/quanxiang-cloud/goalie/pkg/org"
//	"github.com/quanxiang-cloud/goalie/pkg/config"
//
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"github.com/stretchr/testify/suite"
//)
//
//type GoalieSuite struct {
//	suite.Suite
//
//	goalie Goalie
//}
//
//func TestGoalie(t *testing.T) {
//	suite.Run(t, new(GoalieSuite))
//}
//
//func (suite *GoalieSuite) SetupTest() {
//	conf, err := config.NewConfig("../../configs/config.yml")
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), conf)
//
//	err = logger.New(&conf.Log)
//	assert.Nil(suite.T(), err)
//
//	suite.goalie, err = NewGoalie(conf,
//		WithUser(org.NewUserMock()),
//	)
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), suite.goalie)
//
//}
//
//func (suite *GoalieSuite) TestRole() {
//	ctx := logger.ReentryRequsetID(context.Background(), "test-role")
//
//	roles, err := suite.goalie.ListRole(ctx, &ListRoleReq{})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), roles)
//
//}
//
//func (suite *GoalieSuite) TestRoleOwner() {
//	var (
//		userID    = "-1"
//		ownerType = models.Personnel
//	)
//	ctx := logger.ReentryRequsetID(context.Background(), "test-role-owner")
//
//	roles, err := suite.goalie.ListRole(ctx, &ListRoleReq{})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), roles)
//
//	var role, super *ListRoleVO
//	for _, elem := range roles.Roles {
//		if role != nil && super != nil {
//			break
//		}
//		if elem.Tag == string(models.Super) {
//			super = elem
//		} else {
//			role = elem
//		}
//	}
//	require.NotNil(suite.T(), role)
//	require.NotNil(suite.T(), super)
//
//	owners, err := suite.goalie.ListRoleOwner(ctx, &ListRoleOwnerReq{
//		RoleID: role.ID,
//	})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), owners)
//
//	_, err = suite.goalie.UpdateRoleOwner(ctx, &UpdateRoleOwnerReq{
//		RoleID: role.ID,
//		Add: []struct {
//			Type    int
//			OwnerID string
//		}{
//			{
//				Type:    1,
//				OwnerID: userID,
//			},
//		},
//	})
//	assert.Nil(suite.T(), err)
//
//	userRoles, err := suite.goalie.ListUserRole(ctx, &ListUserRoleReq{
//		UserID: userID,
//	})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), userRoles)
//
//	assert.Equal(suite.T(), 1, len(userRoles.Roles))
//	userRole := userRoles.Roles[0]
//	_, err = suite.goalie.UpdateRoleOwner(ctx, &UpdateRoleOwnerReq{
//		RoleID: role.ID,
//		Delete: []string{userRole.ID},
//	})
//	assert.Nil(suite.T(), err)
//
//	_, err = suite.goalie.UpdateOwnerRole(ctx, &UpdateOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//		Add:     []string{role.ID},
//	})
//	assert.Nil(suite.T(), err)
//	_, err = suite.goalie.DeleteOwnerRole(ctx, &DeleteOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//	})
//	assert.Nil(suite.T(), err)
//
//	_, err = suite.goalie.UpdateOwnerRole(ctx, &UpdateOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//		Add:     []string{role.ID},
//	})
//	assert.Nil(suite.T(), err)
//	ownerRoles, err := suite.goalie.ListOwnerRole(ctx, &ListOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//	})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), ownerRoles)
//	assert.Equal(suite.T(), 1, len(ownerRoles.Roles))
//	_, err = suite.goalie.UpdateOwnerRole(ctx, &UpdateOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//		Delete:  []string{ownerRoles.Roles[0].ID},
//	})
//	assert.Nil(suite.T(), err)
//
//	superOwner, err := suite.goalie.ListRoleOwner(ctx, &ListRoleOwnerReq{
//		RoleID: super.ID,
//	})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), superOwner)
//	assert.NotEqual(suite.T(), 0, len(superOwner.Owners))
//
//	_, err = suite.goalie.TransferRoleSuper(ctx, &TransferRoleSuperReq{
//		UserID:     superOwner.Owners[0].OwnerID,
//		Transferee: userID,
//	})
//	assert.Nil(suite.T(), err)
//	_, err = suite.goalie.TransferRoleSuper(ctx, &TransferRoleSuperReq{
//		UserID:     userID,
//		Transferee: superOwner.Owners[0].OwnerID,
//	})
//	assert.Nil(suite.T(), err)
//}
//
//func (suite *GoalieSuite) TestRoleFunc() {
//
//	var (
//		userID    = "-1"
//		ownerType = models.Personnel
//
//		// funcID = "-1"
//	)
//
//	ctx := logger.ReentryRequsetID(context.Background(), "test-role-owner")
//	roles, err := suite.goalie.ListRole(ctx, &ListRoleReq{})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), roles)
//	assert.GreaterOrEqual(suite.T(), len(roles.Roles), 1)
//
//	var role *ListRoleVO
//	for _, elem := range roles.Roles {
//		if models.RoleTag(elem.Tag) != models.Super {
//			role = elem
//			break
//		}
//	}
//	assert.NotNil(suite.T(), role)
//
//	roleFuncs, err := suite.goalie.ListRoleFunc(ctx, &ListRoleFuncReq{
//		RoleID: role.ID,
//	})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), roleFuncs)
//
//	funcTag, err := suite.goalie.ListFuncTag(ctx, &ListFuncTagReq{})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), funcTag)
//
//	_, err = suite.goalie.UpdateOwnerRole(ctx, &UpdateOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//		Add:     []string{role.ID},
//	})
//	assert.Nil(suite.T(), err)
//
//	userFuncTag, err := suite.goalie.ListUserFuncTag(ctx, &ListUserFuncTagReq{
//		UserID: userID,
//	})
//	assert.Nil(suite.T(), err)
//	assert.NotNil(suite.T(), userFuncTag)
//
//	_, err = suite.goalie.DeleteOwnerRole(ctx, &DeleteOwnerRoleReq{
//		Type:    int(ownerType),
//		OwnerID: userID,
//	})
//	assert.Nil(suite.T(), err)
//}
