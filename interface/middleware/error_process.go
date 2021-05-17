package middleware

import (
	"errors"
	"fmt"
	error2 "github.com/damondu/greddit/domain/error"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"reflect"
)

var (
	ErrorInternalError = errors.New("Woops! Something went wrong :(")
)

func ValidationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required;", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s;", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s;", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("Invalid email format;")
	case "len":
		return fmt.Sprintf("%s must be %s characters long;", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid;", e.Field())
}

// Errors This method collects all errors and submits them to Rollbar
func Errors() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Find out what type of error it is
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{"code": error2.CommonError, "errorMsg": e.Error()})
					}
				case gin.ErrorTypeBind:
					var fieldErrs []validator.FieldError
					fieldErrs = e.Err.(validator.ValidationErrors)
					var errorMsg = ""
					for _, err := range fieldErrs {
						errorMsg += ValidationErrorToText(err)
					}

					// Make sure we maintain the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.JSON(status, gin.H{"code": error2.CommonError, "errorMsg": errorMsg})
				case gin.ErrorTypePrivate:
					if reflect.TypeOf(e.Err) == reflect.TypeOf(&error2.ApplicationError{}) {
						appError := e.Err.(*error2.ApplicationError)
						c.JSON(http.StatusForbidden, gin.H{"code": appError.Code, "errorMsg": appError.Msg})
					}
				default:
					// Log all other errors
					if !c.Writer.Written() {
						c.JSON(http.StatusInternalServerError, gin.H{"code": error2.CommonError, "errorMsg": ErrorInternalError.Error()})
						log.Panic(c.Err())
					}
				}
			}
		}
	}
}
