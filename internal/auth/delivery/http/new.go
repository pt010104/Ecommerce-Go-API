package http

import 

type Handler interface {
	Login()
	Signup()
}

type handler struct {
	authUC authUsecase.Usecase
}

func New() Handler {
	return &handler{
		authUC: authUC
	}
}