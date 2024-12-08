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

	files := form.File["files"]
	if len(files) == 0 {
		return models.Scope{}, uploadRequest{}, errNoFiles
	}

	req := uploadRequest{
		Files: files,
	}

	return sc, req, nil
}
