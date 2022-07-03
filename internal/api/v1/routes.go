package v1

func (s *Server) initRoutes() {

	apiV1 := s.Group("/api/v1")
	{
		apiV1.GET("/list", s.handler.makeGetDeliveryHandler(s.ss.deliveryService))
	}
}
