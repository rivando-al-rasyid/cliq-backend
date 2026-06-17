package controller

import (
	"github.com/rivando-al-rasyid/cliq/internals/service"
)

type CliqController struct {
	CliqService *service.CliqService
}

func NewCliqController(CliqService *service.CliqService) *CliqController {
	return &CliqController{CliqService: CliqService}
}
