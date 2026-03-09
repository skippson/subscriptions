package httphandlers

import "github.com/gofiber/fiber/v2"

func (h *ApiHandlers) MapApiRoutes(router fiber.Router, mw Middleware) {
	router.Use(mw.SetRequestID())

	sub := router.Group("/subscriptions")

	sub.Get("/", h.ListSubscriptions())
	sub.Post("/", h.AddNewSubscription())
	sub.Get("/total_cost", h.GetTotalCost())
	sub.Get("/:id", h.GetSubscriptionByID())
	sub.Patch("/:id", h.UpdateSubscription())
	sub.Delete("/:id", h.DeleteSubscription())
}
