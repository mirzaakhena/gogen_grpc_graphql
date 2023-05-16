package graphqlserver

func (r *controller) RegisterRouter() {
	r.fields["reverseMessage"] = r.sendMessageHandler()
}
