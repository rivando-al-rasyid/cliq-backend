package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rivando-al-rasyid/cliq/internals/dto"
	"github.com/rivando-al-rasyid/cliq/internals/pkg"
	"github.com/rivando-al-rasyid/cliq/internals/service"
)

type CliqController struct {
	CliqService *service.CliqService
}

func NewCliqController(cliqService *service.CliqService) *CliqController {
	return &CliqController{CliqService: cliqService}
}

// CreateSlug godoc
// @Summary      Create a new slug
// @Description  Generates a shortened URL slug for an authenticated user.
// @Tags         cliq
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body           body      dto.Link  true  "Slug Creation Payload"
// @Success      201            {object}  dto.Response "Slug created successfully"
// @Failure      400            {object}  dto.Response "Invalid request payload"
// @Failure      401            {object}  dto.Response "Unauthorized / Invalid token"
// @Failure      500            {object}  dto.Response "Internal server error"
// @Router       /link/create [post]
func (c *CliqController) CreateSlug(ctx *gin.Context) {
	claimsRaw, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			dto.NewError("Unauthorized", errors.New("missing claims")),
		)
		return
	}

	claims, ok := claimsRaw.(*pkg.Claims)
	if !ok {
		ctx.JSON(
			http.StatusUnauthorized,
			dto.NewError("Unauthorized", errors.New("invalid claims")),
		)
		return
	}

	var body dto.Link
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("[CliqController.CreateSlug] bind error: %v\n", err)

		ctx.JSON(
			http.StatusBadRequest,
			dto.NewError("Invalid request payload", err),
		)
		return
	}

	slug, err := c.CliqService.CreateSlug(ctx.Request.Context(), claims.ID, body)
	if err != nil {
		log.Printf("[CliqController.CreateSlug] service error: %v\n", err)

		ctx.JSON(
			http.StatusInternalServerError,
			dto.NewError("Create slug failed", err),
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		dto.NewSuccess("Slug created successfully", gin.H{
			"slug": slug,
		}),
	)
}
