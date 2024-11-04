package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
	"net/http"
)

type statusType bool

const (
	SUCCESS statusType = true
	FAILED  statusType = false
)

type ResponseFormat struct {
	Status  statusType `json:"status"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}

func response(context *gin.Context, data any, message string, code int) {
	var response ResponseFormat

	err, isError := data.(apperrors.AppError)
	if isError {
		response = ResponseFormat{
			Status: FAILED,
			Error:  err.Error(),
		}

	} else {
		response = ResponseFormat{
			Data:    data,
			Status:  SUCCESS,
			Message: message,
		}
	}

	context.JSON(code, response)
}

func JsonOk(context *gin.Context, data any) {
	response(context, data, "Ok", http.StatusOK)
}

func JsonDelete(context *gin.Context, data any, resourceType string) {
	response(context, data, fmt.Sprintf("%v Resource Deleted", resourceType), http.StatusOK)
}

func JsonCreated(context *gin.Context, data any, resourceType string) {
	response(context, data, fmt.Sprintf("%v created successfully", resourceType), http.StatusCreated)
}

func JsonModified(context *gin.Context, data any, resourceType string) {
	response(context, data, fmt.Sprintf("%v modified successfully", resourceType), http.StatusAccepted)
}

func JsonError(context *gin.Context, err apperrors.AppError) {
	response(context, err, err.Error(), err.GetCode())
}
