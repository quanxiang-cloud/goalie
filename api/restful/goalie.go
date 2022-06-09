package restful

import (
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/goalie/pkg/header2"
	"gorm.io/gorm"
	"net/http"

	"github.com/quanxiang-cloud/goalie/internal/goalie"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"github.com/quanxiang-cloud/goalie/pkg/org"

	"github.com/gin-gonic/gin"
)

// Goalie goalie
type Goalie struct {
	goalie goalie.Goalie
}

// NewGoalie new goalie
func NewGoalie(conf *config.Config, db *gorm.DB) (*Goalie, error) {
	g, err := goalie.NewGoalie(conf, db,
		goalie.WithUser(org.NewUser(conf.InternalNet)),
	)
	if err != nil {
		return nil, err
	}
	return &Goalie{
		goalie: g,
	}, nil
}

// TransferRoleSuper 转让超级管理员
func (g *Goalie) TransferRoleSuper(c *gin.Context) {
	req := &goalie.TransferRoleSuperReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	profile := header2.GetProfile(c)
	req.UserID = profile.UserID

	resp.Format(g.goalie.TransferRoleSuper(ginheader.MutateContext(c), req)).Context(c)
}
