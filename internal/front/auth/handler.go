package auth

import (
	response "api-mini-shop/pkg/http/response"
	"api-mini-shop/pkg/utls"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	DBPool      *sqlx.DB
	AuthService *AuthService
}

func NewAuthHandler(db_pool *sqlx.DB) *AuthHandler {
	return &AuthHandler{
		DBPool:      db_pool,
		AuthService: NewAuthService(db_pool),
	}
}

// @Summary      Login
// @Description  Authenticates a user and returns a token
// @Tags         Admin/Auth
// @Accept       json
// @Produce      json
// @Param        user  body      auth.LoginRequest  true  "Credentials to use"
// @Success      200   {object}  auth.LoginResponse
// @Failure      400   {object}  utls.Error
// @Failure      401   {object}  utls.Error
// @Router       /front/auth [post]
func (au *AuthHandler) Login(c *fiber.Ctx) error {
	var login_request LoginRequest
	v := utls.NewValidator()

	if err := login_request.Bind(c, v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utls.Translate("login_failed", nil, c),
				-1000,
				err,
			),
		)
	}

	resp, err := au.AuthService.Login(login_request.UserName, login_request.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			response.NewResponseError(
				utls.Translate(err.MessageID, nil, c),
				-1000,
				fmt.Errorf(utls.Translate(err.Err.Error(), nil, c)),
			),
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.NewResponse(
			utls.Translate("login_success", nil, c),
			1000,
			resp,
		),
	)
}
