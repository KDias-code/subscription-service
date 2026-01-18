package app

func (s Server) router() {
	subRouter := s.app.Group("/v1/subscriptions")

	subRouter.Get("/sum", s.handlers.SubscriptionsSum)
	subRouter.Get("/:id", s.handlers.FindByID)
	subRouter.Post("", s.handlers.Create)
	subRouter.Put("", s.handlers.Update)
	subRouter.Delete("/:id", s.handlers.Delete)
}
