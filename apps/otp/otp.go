package otp

import (
	"github.com/thrillee/triq/internals/modules"
	"github.com/thrillee/triq/internals/schemas"
)

var nameApp = "OTP"

var appConfig = modules.ModuleConfig{Name: nameApp}

func mountRoutes(appFunc modules.FiberAppFunc) {
	appConfig.Routes = append(appConfig.Routes, appFunc)
}

func migrateModel(model schemas.Model) {
	appConfig.Models = append(appConfig.Models, &model)
}

func init() {
	migrateModel(OTP{})

	// mountRoutes(NewTodoREST)
}

func GetModuleConfig() *modules.ModuleConfig {
	return &appConfig
}
