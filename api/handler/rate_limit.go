package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kavehjamshidi/arvan-challenge/service/rate_limit/contract"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"net/http"
)

type RateLimitHandler struct {
	rateLimitService contract.RateLimitService
}

func NewRateLimitHandler(
	rateLimitService contract.RateLimitService,
) RateLimitHandler {
	return RateLimitHandler{
		rateLimitService: rateLimitService,
	}
}

func (u RateLimitHandler) HandleRateLimit(c *fiber.Ctx) error {
	userID := c.Get("user-id")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(response[string]{
			Error:   "user-id header is required",
			Message: MsgValidationError,
		})
	}

	err := u.rateLimitService.CheckRateLimit(c.Context(), userID)
	if err != nil {
		if errors.Is(err, utils.ErrTooManyRequests) {
			return c.Status(http.StatusTooManyRequests).JSON(response[string]{
				Error:   err.Error(),
				Message: MsgTooManyRequests,
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(response[string]{
			Error:   err.Error(),
			Message: MsgFailed,
		})
	}

	return c.Next()
}
