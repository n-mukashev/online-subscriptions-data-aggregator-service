package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindJSONWithValidation(ctx *gin.Context, obj any) bool {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			out := make(map[string]string)
			val := reflect.ValueOf(obj).Elem()
			typ := val.Type()
			for _, fe := range verr {
				field := fe.StructField()
				jsonName := field

				if f, ok := typ.FieldByName(field); ok {
					tag := f.Tag.Get("json")
					if tag != "" {
						jsonName = strings.Split(tag, ",")[0] // берём до запятой
					}
				}

				switch fe.Tag() {
				case "required":
					out[jsonName] = "field is required"
				default:
					out[jsonName] = fmt.Sprintf("failed validation: %s", fe.Tag())
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
			return false
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return false
	}
	return true
}
