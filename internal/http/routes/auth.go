package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	authDTO "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	_ "github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/utils"
)

func Auth(r fiber.Router, signInUc auth.SignInUseCase, signUpUc auth.SignUpUseCase, refreshUc auth.RefreshTokenUseCase) {
	r.Post("/sign-in", signIn(signInUc))
	r.Post("/sign-up", signUp(signUpUc))
	r.Post("/token/refresh", tokenRefresh(refreshUc))
}

// signIn godoc
// @Summary      Sign in a client
// @Description  Authenticate client using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.SignInParams  true  "SignIn credentials"
// @Success      200          {object}  dto.AuthTokens    "Access and Refresh Tokens"
// @Failure      401          {object}  utils.APIError
// @Failure      404          {object}  utils.APIError
// @Failure      422          {object}  utils.APIError "Invalid params"
// @Failure      500          {object}  utils.APIError
// @Router       /auth/sign-in [post]
func signIn(uc auth.SignInUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params authDTO.SignInParams

		if err := c.BodyParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		t, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(t)
	}
}

// signUp godoc
// @Summary      Sign up a new client
// @Description  Register a new client with name, email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        registration  body      dto.SignUpParams  true  "SignUp registration data"
// @Success      201           {object}  dto.Client        "Created client"
// @Failure      400           {object}  utils.APIError
// @Failure      422           {object}  utils.APIError "Invalid params"
// @Failure      409           {object}  utils.APIError  "Conflict (email exists)"
// @Failure      500           {object}  utils.APIError
// @Router       /auth/sign-up [post]
func signUp(uc auth.SignUpUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params authDTO.SignUpParams

		if err := c.BodyParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		cl, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusCreated).JSON(cl)
	}
}

// tokenRefresh godoc
// @Summary      Refresh tokens
// @Description  Generate new access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        tokenRefresh  body      dto.RefreshTokenParams  true  "Refresh token data"
// @Success      200           {object}  dto.AuthTokens          "New tokens"
// @Failure      400           {object}  utils.APIError
// @Failure      401           {object}  utils.APIError
// @Failure      422           {object}  utils.APIError "Invalid params"
// @Failure      500           {object}  utils.APIError
// @Router       /auth/token/refresh [post]
func tokenRefresh(uc auth.RefreshTokenUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params authDTO.RefreshTokenParams

		if err := c.BodyParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		t, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(t)
	}
}
