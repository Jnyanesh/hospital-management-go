package doctor

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/HemanthGowda/hospital-management-go/internal/auth"
	"github.com/HemanthGowda/hospital-management-go/internal/domain"
	"github.com/HemanthGowda/hospital-management-go/internal/httpx"
	"github.com/HemanthGowda/hospital-management-go/internal/user"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var request RegisterDoctorRequest
	principal, ok := auth.PrincipalFrom(c.Context())

	if !ok {
		return domain.Unauthorized("Authentication is required")
	}

	if !principal.HasAny(user.RoleAdmin, user.RoleDoctor) {
		return domain.Forbidden()
	}
	if err := httpx.Decode(c, &request); err != nil {
		return err
	}

	response, err := h.service.Register(c.Context(), request)
	if err != nil {
		return err
	}
	return httpx.WriteJSON(c, fiber.StatusCreated, response)
}

func (h *Handler) Get(c fiber.Ctx) error {
	principal, ok := auth.PrincipalFrom(c.Context())

	if !ok {
		return domain.Unauthorized("Authentication is required")
	}
	if !principal.HasAny(user.RoleAdmin, user.RoleDoctor, user.RolePatient) {
		return domain.Forbidden()
	}
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return domain.BadRequest("Doctor ID is required")
	}
	response, err := h.service.GetDoctorById(c.Context(), id)
	if err != nil {
		return domain.NotFound("Doctor not found")
	}
	return httpx.WriteJSON(c, fiber.StatusOK, response)
}
