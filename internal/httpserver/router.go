package httpserver

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/HemanthGowda/hospital-management-go/internal/auth"
	"github.com/HemanthGowda/hospital-management-go/internal/doctor"
	"github.com/HemanthGowda/hospital-management-go/internal/domain"
	"github.com/HemanthGowda/hospital-management-go/internal/httpx"
	"github.com/HemanthGowda/hospital-management-go/internal/patient"
	"github.com/HemanthGowda/hospital-management-go/internal/user"
	"github.com/HemanthGowda/hospital-management-go/openapi"
)

type Dependencies struct {
	Ping     func(context.Context) error
	Tokens   *auth.JWTService
	Users    *user.Handler
	Auth     *auth.Handler
	Patients *patient.Handler
	Doctors  *doctor.Handler
}

func NewRouter(deps Dependencies) (*fiber.App, error) {
	spec, err := openapi.JSON()
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		ReadTimeout: 5 * time.Second,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return httpx.WriteError(c, err)
		},
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(_ fiber.Ctx, recovered any) {
			slog.Error("panic", "error", recovered, "stack", string(debug.Stack()))
		},
	}))
	app.Use(auth.Middleware(deps.Tokens))

	app.Get("/actuator/health", healthHandler(deps.Ping))
	app.Get("/actuator/info", infoHandler)
	app.Get("/v3/api-docs", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Send(spec)
	})
	app.Get("/swagger-ui.html", swaggerUI)
	app.Post("/api/v1/auth/login", deps.Auth.Login)
	app.Post("/api/v1/users", deps.Users.Register)
	app.Put("/api/v1/patients/:id", deps.Patients.Update)
	app.Post("/api/v1/doctors", deps.Doctors.Register)
	app.Use(unmatchedRoute)

	return app, nil
}

func healthHandler(ping func(context.Context) error) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		status := fiber.StatusOK
		body := map[string]string{"status": "UP"}
		if err := ping(ctx); err != nil {
			status = fiber.StatusServiceUnavailable
			body["status"] = "DOWN"
		}
		return c.Status(status).JSON(body)
	}
}

func infoHandler(c fiber.Ctx) error {
	return httpx.WriteJSON(c, fiber.StatusOK, map[string]string{
		"name":    "hospital-management-go",
		"version": "0.1.0",
	})
}

func swaggerUI(c fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <title>Hospital Management API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css"/>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
<script>
window.ui = SwaggerUIBundle({
  url: "/v3/api-docs",
  dom_id: "#swagger-ui",
  persistAuthorization: true
});
</script>
</body>
</html>`)
}

func unmatchedRoute(c fiber.Ctx) error {
	if _, ok := auth.PrincipalFrom(c.Context()); !ok {
		return domain.Unauthorized("Authentication is required")
	}
	return domain.NotFound("Resource not found")
}
