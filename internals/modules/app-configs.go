package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thrillee/triq/internals/schemas"
)

type Module interface {
	GetModuleConfig() *ModuleConfig
}

type FiberAppFunc func(*fiber.App)

type ModuleConfig struct {
	Models []*schemas.Model
	Routes []FiberAppFunc
	Name   string
}
