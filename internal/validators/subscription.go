package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/models"
)

func RegisterValidators(v *validator.Validate) {
	v.RegisterValidation("monthyear", func(fl validator.FieldLevel) bool {
		m, ok := fl.Field().Interface().(models.MonthYear)
		if !ok {
			return false
		}
		t := m.ToTime()
		month := t.Month()
		year := t.Year()
		return month >= 1 && month <= 12 && year > 0
	})

	v.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		id, ok := fl.Field().Interface().(uuid.UUID)
		return ok && id.String() != ""
	})

	v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		str, ok := fl.Field().Interface().(string)
		return ok && strings.TrimSpace(str) != ""
	})
}
