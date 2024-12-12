package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processUploadRequest(c *gin.Context) (models.Scope, uploadRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "media.delivery.http.processUploadRequest: unauthorized")
		return models.Scope{}, uploadRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	form, err := c.MultipartForm()
	if err != nil {
		h.l.Errorf(ctx, "media.delivery.http.processUploadRequest: %v", err)
		return models.Scope{}, uploadRequest{}, errWrongBody
	}

	req := uploadRequest{
		Files: form.File["files"],
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "media.delivery.http.processUploadRequest: %v", err)
		return models.Scope{}, uploadRequest{}, err
	}

	return sc, req, nil
}
