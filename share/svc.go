package share

import "github.com/alphabatem/nft-proxy/share/repository/sqlrepo"

//this file is services context manually instead of "github.com/babilu-online/common/context"

type serviceContext struct {
	cfg EnvConfig
	db  sqlrepo.DbContext
	v1  *gin.RouterGroup
}
type ServiceContext interface {
	GetCfg() EnvConfig
	GetV1() *gin.RouterGroup
	GetDb() sqlrepo.DbContext
}

func NewServiceContext(
	cfg EnvConfig,
	v1 *gin.RouterGroup,
	db sqlrepo.DbContext,
) *serviceContext {
	return &serviceContext{
		cfg: cfg,
		v1:  v1,
		db:  db,
	}
}

func (svc *serviceContext) GetCfg() EnvConfig {
	return svc.cfg
}
func (svc *serviceContext) GetV1() *gin.RouterGroup {
	return svc.v1
}
func (svc *serviceContext) GetDb() sqlrepo.DbContext {
	return svc.db
}
