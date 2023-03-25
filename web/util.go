package web

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	. "github.com/broodjeaap/go-watch/models"
)

func bindAndValidateWatch(watch *Watch, c *gin.Context) (map[string]string, error) {
	err := c.ShouldBind(watch)
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

func buildFilterTree(filters []Filter, connections []FilterConnection) {
	filterMap := make(map[FilterID]*Filter, len(filters))
	for i := range filters {
		filter := &filters[i]
		filterMap[filter.ID] = filter
	}
	for i := range connections {
		connection := &connections[i]
		outputFilter := filterMap[connection.OutputID]
		inputFilter := filterMap[connection.InputID]

		outputFilter.Children = append(outputFilter.Children, inputFilter)
		//log.Println("Adding", inputFilter.Name, "as child to", outputFilter.Name)
		inputFilter.Parents = append(inputFilter.Parents, outputFilter)
		//log.Println("Adding", outputFilter.Name, "as parent to", inputFilter.Name)
	}
	/*
		for _, filter := range filters {
			log.Println("Children of", filter.Name)
			for _, child := range filter.Children {
				log.Println("   ", child.Name)
			}

			log.Println("Parents of", filter.Name)
			for _, parent := range filter.Parents {
				log.Println("   ", parent.Name)
			}
		}
	*/
}
