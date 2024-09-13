package servers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type HttpServerProps struct {
	Host string
	Port int
}

type HttpServer struct {
	server    *fiber.App
	AddrProps *HttpServerProps
	mounted   bool
}

var httpServer *HttpServer = &HttpServer{
	server:  new(),
	mounted: true,
}

func new() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	return app
}

func GetAppServer() *fiber.App {
	return httpServer.server
}

func (p *HttpServer) setHttpAddr(addr *HttpServerProps) {
	p.AddrProps = addr
}

func (p *HttpServer) listen() {
	addr := fmt.Sprintf("%s:%d", p.AddrProps.Host, p.AddrProps.Port)

	log.Fatal(p.server.Listen(addr))
}

func ListenAndServe(addr *HttpServerProps) {
	httpServer.setHttpAddr(addr)
	httpServer.listen()
}
