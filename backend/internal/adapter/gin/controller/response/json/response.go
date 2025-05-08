package response

import (
	iresponse "LCA/internal/adapter/gin/controller/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Object struct {
	SuccessBool bool   `json:"success"`        // if request success or not
	Message     string `json:"message"`        // message from request
	Data        *any   `json:"data,omitempty"` // data from request
}

func NewJSONResponse(logger *zap.Logger) iresponse.IResponse {
	return &JSONResponse{logger: logger}
}

// JSONHandler is a function to handle response.
func (jr Object) JSONHandler(c *gin.Context, httpStatus int) {
	c.JSON(httpStatus, gin.H{
		"success": jr.SuccessBool,
		"message": jr.Message,
		"data":    jr.Data,
	})
}

type JSONResponse struct {
	logger *zap.Logger
}

// Success is a function to response success with message
func (JSONResponse) Success(c *gin.Context, message string) {
	Object{
		SuccessBool: true,
		Message:     message,
		Data:        nil,
	}.JSONHandler(c, http.StatusOK)
}

// SuccessWithData is a function to response success with message and data
func (JSONResponse) SuccessWithData(c *gin.Context, message string, data any) {
	Object{
		SuccessBool: true,
		Message:     message,
		Data:        &data,
	}.JSONHandler(c, http.StatusOK)
}

// Fail is a function to response error with message
func (JSONResponse) Fail(c *gin.Context, message string) {
	Object{
		SuccessBool: false,
		Message:     message,
		Data:        nil,
	}.JSONHandler(c, http.StatusServiceUnavailable)
}

// FailWithError is a function to response error with message and error
func (jr JSONResponse) FailWithError(c *gin.Context, message string, err error) {
	jr.logger.Error(err.Error(),
		zap.String("message", message),
	)
	Object{
		SuccessBool: false,
		Message:     message,
		Data:        nil,
	}.JSONHandler(c, http.StatusServiceUnavailable)
}

// ValidatorFail is a function to response error with message
func (JSONResponse) ValidatorFail(c *gin.Context, message string) {
	Object{
		SuccessBool: false,
		Message:     message,
		Data:        nil,
	}.JSONHandler(c, http.StatusBadRequest)
}

// AuthFail is a function to response error with message
func (JSONResponse) AuthFail(c *gin.Context, message string) {
	Object{
		SuccessBool: false,
		Message:     message,
		Data:        nil,
	}.JSONHandler(c, http.StatusUnauthorized)
}
