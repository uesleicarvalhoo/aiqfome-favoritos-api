package routes

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/context"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/utils"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
)

func Me(r fiber.Router,
	getClientFavoritesUc favorites.GetClientFavoritesUseCase,
	addProductToFavoritesUc favorites.AddProductToFavoritesUseCase,
	removeProductFromFavoritesUc favorites.RemoveProductFromFavoritesUseCase,
) {
	r.Get("/", getMe())
	r.Get("/favorites", getClientFavorites(getClientFavoritesUc))
	r.Post("/favorites", addProductToFavorites(addProductToFavoritesUc))
	r.Delete("/favorites/product/:id", removeProductFromFavorites(removeProductFromFavoritesUc))
}

// @Summary      Get client favorites
// @Description  Retrieve paginated list of favorite products for the authenticated client
// @Tags         Me/Favorites
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number, starts from 0"
// @Param        pageSize  query     int  false  "Items per page, default 10"
// @Success      200       {object}  dto.ClientFavorites
// @Failure      422       {object}  utils.APIError "Invalid params"
// @Failure      401       {object}  utils.APIError
// @Failure      404       {object}  utils.APIError
// @Failure      403       {object}  utils.APIError
// @Failure      500       {object}  utils.APIError
// @Security     BearerAuth
// @Router       /me/favorites [get]
func getClientFavorites(uc favorites.GetClientFavoritesUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params dto.GetClientFavoritesParams

		if err := c.QueryParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		cl, err := context.GetClient(c.UserContext())
		if err != nil {
			return utils.WriteError(c, err)
		}

		params.ClientID = cl.ID

		fv, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(fv)
	}
}

// @Summary      Add product to favorites
// @Description  Add a product to the authenticated client's favorites list
// @Tags         Me/Favorites
// @Accept       json
// @Produce      json
// @Param        favorite  body      dto.AddProductToFavoritesParams  true  "Product to add"
// @Success      200       {object}  dto.ProductFavorite             "Added favorite"
// @Failure      401       {object}  utils.APIError
// @Failure      404       {object}  utils.APIError
// @Failure      422       {object}  utils.APIError "Invalid params"
// @Failure      500       {object}  utils.APIError
// @Security     BearerAuth
// @Router       /me/favorites [post]
func addProductToFavorites(uc favorites.AddProductToFavoritesUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params dto.AddProductToFavoritesParams

		if err := c.BodyParser(&params); err != nil {
			return utils.WriteError(c, err)
		}

		cl, err := context.GetClient(c.UserContext())
		if err != nil {
			return utils.WriteError(c, err)
		}

		params.ClientID = cl.ID
		p, err := uc.Execute(c.UserContext(), params)
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(p)
	}
}

// @Summary      Remove product from favorites
// @Description  Remove a product from the authenticated client's favorites list
// @Tags         Me/Favorites
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  nil "Success"
// @Failure      401  {object}  utils.APIError
// @Failure      404  {object}  utils.APIError
// @Failure      422  {object}  utils.APIError "Invalid params"
// @Failure      500  {object}  utils.APIError
// @Security     BearerAuth
// @Router       /me/favorites/product/{id} [delete]
func removeProductFromFavorites(uc favorites.RemoveProductFromFavoritesUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return utils.WriteError(c, domainerror.Wrap(err, domainerror.InvalidParams, "id do produto inválido", map[string]any{
				"product_id": c.Params("id"),
			}))
		}

		cl, err := context.GetClient(c.UserContext())
		if err != nil {
			return utils.WriteError(c, err)
		}

		params := dto.RemoveProductFromFavoritesParams{
			ClientID:  cl.ID,
			ProductID: pID,
		}

		if err := uc.Execute(c.UserContext(), params); err != nil {
			return utils.WriteError(c, err)
		}

		return c.SendStatus(http.StatusOK)
	}
}

// @Summary      Get current client data
// @Description  Get current client data
// @Tags         Me
// @Accept       json
// @Produce      json
// @Success      200  {object}  nil "Success"
// @Failure      401  {object}  utils.APIError
// @Failure      404  {object}  utils.APIError
// @Security     BearerAuth
// @Router       /me [get]
func getMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl, err := context.GetClient(c.UserContext())
		if err != nil {
			return utils.WriteError(c, err)
		}

		return c.Status(http.StatusOK).JSON(cl)
	}
}
