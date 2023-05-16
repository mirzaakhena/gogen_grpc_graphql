package restapi

func (r *controller) RegisterRouter() {
	resource := r.Router.Group("/api/v1", r.authentication())
	resource.POST("/runmessagesend", r.authorization(), r.runMessageSendHandler())
	//!

}
