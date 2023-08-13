package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nikitamirzani323/wl_agen_backend_api/controllers"
	"github.com/nikitamirzani323/wl_agen_backend_api/middleware"
)

func Init() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		// c.Set("Content-Security-Policy", "frame-ancestors 'none'")
		// c.Set("X-XSS-Protection", "1; mode=block")
		// c.Set("X-Content-Type-Options", "nosniff")
		// c.Set("X-Download-Options", "noopen")
		// c.Set("Strict-Transport-Security", "max-age=5184000")
		// c.Set("X-Frame-Options", "SAMEORIGIN")
		// c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Get("/ipaddress", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      "data",
			"BASEURL":     c.BaseURL(),
			"HOSTNAME":    c.Hostname(),
			"IP":          c.IP(),
			"IPS":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomain":   c.Subdomains(),
		})
	})
	app.Get("/dashboard", monitor.New())

	app.Post("/api/login", controllers.CheckLogin)
	app.Post("/api/valid", middleware.JWTProtected(), controllers.Home)
	app.Post("/api/alladmin", middleware.JWTProtected(), controllers.Adminhome)
	app.Post("/api/saveadmin", middleware.JWTProtected(), controllers.AdminSave)
	app.Post("/api/alladminrule", middleware.JWTProtected(), controllers.Adminrulehome)
	app.Post("/api/saveadminrule", middleware.JWTProtected(), controllers.AdminruleSave)

	app.Post("/api/curr", middleware.JWTProtected(), controllers.Currhome)
	app.Post("/api/currsave", middleware.JWTProtected(), controllers.CurrSave)
	app.Post("/api/provider", middleware.JWTProtected(), controllers.Providerhome)
	app.Post("/api/providersave", middleware.JWTProtected(), controllers.ProviderSave)
	app.Post("/api/categame", middleware.JWTProtected(), controllers.CateGamehome)
	app.Post("/api/categamesave", middleware.JWTProtected(), controllers.CateGameSave)
	app.Post("/api/gamesave", middleware.JWTProtected(), controllers.GameSave)
	app.Post("/api/catebank", middleware.JWTProtected(), controllers.CateBankhome)
	app.Post("/api/catebanksave", middleware.JWTProtected(), controllers.CateBankSave)
	app.Post("/api/banktypesave", middleware.JWTProtected(), controllers.BankTypeSave)
	app.Post("/api/master", middleware.JWTProtected(), controllers.Masterhome)
	app.Post("/api/masteragenadmin", middleware.JWTProtected(), controllers.Masteragenadmin)
	app.Post("/api/mastersave", middleware.JWTProtected(), controllers.MasterSave)
	app.Post("/api/masteradminsave", middleware.JWTProtected(), controllers.MasteradminSave)
	app.Post("/api/masteragensave", middleware.JWTProtected(), controllers.MasteragenSave)
	app.Post("/api/masteragenadminsave", middleware.JWTProtected(), controllers.MasteragenadminSave)
	app.Post("/api/domain", middleware.JWTProtected(), controllers.Domainhome)
	app.Post("/api/domainsave", middleware.JWTProtected(), controllers.DomainSave)

	return app
}
