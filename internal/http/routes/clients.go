package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/middleware"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/utils"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

func Clients(r fiber.Router,
	authorizeUc auth.AuthorizeUseCase,
	findClientUc client.FindClientUseCase,
	listClientsUc client.ListClientsUseCase,
	updateClientUc client.UpdateClientUseCase,
	deleteClientUc client.DeleteClientUseCase,
) {
	r.Get("/:id", middleware.Authorize(authorizeUc, role.ResourceClient, role.ActionRead), findClient(findClientUc))
	r.Get("/", middleware.Authorize(authorizeUc, role.ResourceClient, role.ActionRead), listClients(listClientsUc))
	r.Patch("/:id", middleware.Authorize(authorizeUc, role.ResourceClient, role.ActionWrite), updateClient(updateClientUc))
	r.Delete("/:id", middleware.Authorize(authorizeUc, role.ResourceClient, role.ActionDelete), deleteClient(deleteClientUc))
}

// @Summary      Get client
// @Description  Get client data by the given ID
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        id             path      string                true  "Client ID (UUID)"
// @Success      200            {object}  dto.Client
// @Failure      400            {object}  utils.APIError
// @Failure      401            {object}  utils.APIError
// @Failure      403            {object}  utils.APIError
// @Failure      404            {object}  utils.APIError
// @Failure      500            {object}  utils.APIError
// @Security     BearerAuth
// @Router       /clients/{id} [get]
func findClient(uc client.FindClientUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cId, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.WriteError(c, err)
		}

		cl, err := uc.Execute(c.UserContext(), cId)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(cl)
	}
}

// @Summary      Listar clients
// @Description  Return a list with paginated clients
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        page           query     int     false "Page number"
// @Param        pageSize       query     int     false "Clients per page, default 10"
// @Success      200            {object}  dto.PaginatedClients
// @Failure      400            {object}  utils.APIError
// @Failure      401            {object}  utils.APIError
// @Failure      403            {object}  utils.APIError
// @Failure      500            {object}  utils.APIError
// @Security     BearerAuth
// @Router       /clients [get]
func listClients(uc client.ListClientsUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params dto.ListClientsParams

		if err := c.QueryParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		if params.PageSize <= 0 {
			params.PageSize = 10
		}

		cc, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(cc)
	}
}

// @Summary      Update client
// @Description  Update client data by the given ID
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        id             path      string                true  "Client ID (UUID)"
// @Param        body           body      dto.UpdateClientParams true  "Fields to update, if nil, it will be ignored"
// @Success      200            {object}  dto.Client
// @Failure      400            {object}  utils.APIError
// @Failure      401            {object}  utils.APIError
// @Failure      403            {object}  utils.APIError
// @Failure      404            {object}  utils.APIError
// @Failure      500            {object}  utils.APIError
// @Security     BearerAuth
// @Router       /clients/{id} [patch]
func updateClient(uc client.UpdateClientUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params dto.UpdateClientParams

		if err := c.BodyParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		cId, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.WriteError(c, err)
		}

		params.ClientID = cId

		cl, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(cl)
	}
}

// @Summary      Remove client
// @Description  Remove the cliente by the given ID
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        id             path      string  true  "Client ID (UUID)"
// @Success      204            {object}  nil     "No Content"
// @Failure      401            {object}  utils.APIError
// @Failure      403            {object}  utils.APIError
// @Failure      404            {object}  utils.APIError
// @Failure      500            {object}  utils.APIError
// @Security     BearerAuth
// @Router       /clients/{id} [delete]
func deleteClient(uc client.DeleteClientUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cId, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return utils.WriteError(c, err)
		}

		if err := uc.Execute(c.UserContext(), dto.DeleteClientParams{
			ClientID: cId,
		}); err != nil {
			return utils.WriteError(c, err)
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
