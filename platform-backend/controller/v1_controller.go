package controller

import "platform-backend/service"

type V1Controller struct {
	svc service.V1Service
}

func NewV1Controller(svc service.V1Service) *V1Controller {
	return &V1Controller{svc: svc}
}
