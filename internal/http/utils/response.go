package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
)

func WriteResponse(c *fiber.Ctx, body any, status int) error {
	if body == nil {
		return c.SendStatus(status)
	}

	return c.Status(status).JSON(body)
}

func WriteRawResponse(c *fiber.Ctx, body []byte) error {
	return c.Status(http.StatusOK).Send(body)
}

func WriteError(c *fiber.Ctx, err error) error {
	apiErr, sc := toAPIError(err)

	span := trace.SpanFromContext(c.UserContext())
	trace.AddSpanTags(span, map[string]string{
		"error_code": apiErr.Code,
	})

	logger.ErrorF(c.UserContext(), fmt.Sprintf("error on route %s", c.Route().Path), logger.Fields{
		"message": apiErr.Message,
		"code":    apiErr.Code,
	})

	return WriteResponse(c, apiErr, sc)
}

func toAPIError(err error) (APIError, int) {
	if domainErr, ok := err.(*domainerror.Error); ok {
		dt := domainErr.Details
		if dt == nil {
			dt = make(map[string]any)
		}

		if _, ok := dt["error"]; !ok && domainErr.Cause != nil {
			dt["error"] = domainErr.Cause.Error()
		}

		return APIError{
			Code:      domainErr.Code.String(),
			Message:   domainErr.Message,
			Timestamp: time.Now().Format(time.RFC3339),
			Details:   domainErr.Details,
		}, domainerror.StatusCode(domainErr.Code)
	}

	return APIError{
		Code:    domainerror.Default.String(),
		Message: err.Error(),
	}, http.StatusInternalServerError
}
