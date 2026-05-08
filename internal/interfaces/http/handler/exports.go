package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandleError is the exported variant of handleError for use by handler sub-packages.
func HandleError(c *gin.Context, err error) {
	handleError(c, err)
}

// GetUserIDFromContext is the exported variant of getUserIDFromContext for use by handler sub-packages.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	return getUserIDFromContext(c)
}
