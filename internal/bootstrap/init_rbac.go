package bootstrap

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type casbinOption struct {
	Prefix string
}

func NewCasbinOption(conf *viper.Viper) casbinOption {
	return casbinOption{
		Prefix: conf.GetString("casbin.prefix"),
	}
}

func NewCasbin(gorm *gorm.DB, option casbinOption) *casbin.SyncedEnforcer {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(gorm, option.Prefix, "permission")
	if err != nil {
		panic(err)
	}

	casbinModel, err := model.NewModelFromFile("internal/adapter/casbin/casbinText.go")
	if err != nil {
		panic(err)
	}

	syncedCachedEnforcer, err := casbin.NewSyncedEnforcer(casbinModel, adapter)
	if err != nil {
		panic(err)
	}

	err = syncedCachedEnforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
	return syncedCachedEnforcer
}
