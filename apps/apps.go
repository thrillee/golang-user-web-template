package apps

import (
	"github.com/sirupsen/logrus"
	"github.com/thrillee/triq/apps/otp"
	"github.com/thrillee/triq/apps/users"
	"github.com/thrillee/triq/internals/modules"
	"github.com/thrillee/triq/internals/schemas"
	"github.com/thrillee/triq/internals/servers"
)

type installedApps struct {
	apps []*modules.ModuleConfig
}

var config installedApps = installedApps{}

func installApp(c *modules.ModuleConfig) {
	config.apps = append(config.apps, c)
}

func mountRoutes(routes []modules.FiberAppFunc) {
	for _, r := range routes {
		r(servers.GetAppServer())
	}
}

func mountMigratable(models []*schemas.Model) {
	for _, m := range models {
		schemas.AddToMigratables(m)
	}
}

func init() {
	installApp(users.GetModuleConfig())
	installApp(otp.GetModuleConfig())
}

func MountApps() {
	appCount := config.apps

	logrus.WithFields(logrus.Fields{
		"Total-apps": len(appCount),
	}).Println("Mount Apps")

	for _, v := range config.apps {
		// log.Printf("App %v\n", v.Name)

		mountRoutes(v.Routes)
		mountMigratable(v.Models)
	}
}
