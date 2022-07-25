package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func bindAndValidateWatch(watch *Watch, c *gin.Context) (map[string]string, error) {
	err := c.ShouldBind(watch)
	return validate(err), err
}

func bindAndValidateURL(url *URL, c *gin.Context) (map[string]string, error) {
	err := c.ShouldBind(url)
	return validate(err), err
}

func bindAndValidateQuery(query *Query, c *gin.Context) (map[string]string, error) {
	err := c.ShouldBind(query)
	return validate(err), err
}

func bindAndValidateFilter(filter *Filter, c *gin.Context) (map[string]string, error) {
	err := c.ShouldBind(filter)
	return validate(err), err
}

func bindAndValidateQueryUpdate(query *QueryUpdate, c *gin.Context) (map[string]string, error) {
	err := c.ShouldBind(query)
	return validate(err), err
}

func prettyError(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	default:
		return "No prettyError for " + fieldError.Tag()
	}
}

func validate(err error) map[string]string {
	out := make(map[string]string)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				out[fe.Field()] = prettyError(fe)
			}
		}
	}
	return out
}
