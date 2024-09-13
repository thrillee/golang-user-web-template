package users

import (
	"github.com/thrillee/triq/internals/modules"
	"github.com/thrillee/triq/internals/schemas"
)

var nameApp = "User"

var appConfig = modules.ModuleConfig{Name: nameApp}

func mountRoutes(appFunc modules.FiberAppFunc) {
	appConfig.Routes = append(appConfig.Routes, appFunc)
}

func migrateModel(model schemas.Model) {
	appConfig.Models = append(appConfig.Models, &model)
}

func init() {
	migrateModel(User{})

	mountRoutes(NewUserREST)

	// log.Println(fmt.Sprintf("Mounting %s", nameApp))
}

func GetModuleConfig() *modules.ModuleConfig {
	return &appConfig
}
