package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 只处理已经产生的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v\n", err.Error())

			// 根据错误类型返回不同的状态码
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
			}
		}
	}
}
