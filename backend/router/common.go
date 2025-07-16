package router

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&BaseRouter{},
		&WecomRouter{},
		&LLMAppRouter{},
		&KFRouter{},
		&StaffRouter{},
		&PolicyRouter{},
	}
}
