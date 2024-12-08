package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary Upload media files
// @Description Upload one or multiple media files
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param files formData file true "Media files to upload"
// @Success 200 {object} response.Resp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/media/upload [post]
func (h handler) Upload(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processUploadRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "media.delivery.http.Upload: %v", err)
		response.Error(c, err)
		return
	}

	err = h.uc.Upload(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "media.delivery.http.Upload: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}
