package common

import "github.com/gofiber/fiber/v2"

type URLHandler struct {
	Path    string
	Handler func(*fiber.Ctx) error
}

type URLRegister struct {
	urls []URLHandler
}

func (u URLRegister) GetURLs() []URLHandler {
	return u.urls
}

func (u URLRegister) RegisterURL(path string, handler func(*fiber.Ctx) error) {
	u.urls = append(u.urls, URLHandler{
		Path:    path,
		Handler: handler,
	})
}

func CreateURLRegister() *URLRegister {
	return &URLRegister{
		urls: []URLHandler{},
	}
}
