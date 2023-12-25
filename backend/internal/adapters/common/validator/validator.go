package validator

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid"
)

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("ISO8601", iso8601)
	v.RegisterValidation("imagex_name", imagexName)
	v.RegisterValidation("imagex_url", imagexUrl)
	return v
}

func imagexUrl(fl validator.FieldLevel) bool {
	regString := `https?://.*(jpe?g|png)$`
	reg := regexp.MustCompile(regString)
	return reg.MatchString(fl.Field().String())
}

func imagexName(fl validator.FieldLevel) bool {
	regString := `(jpe?g|png)$`
	reg := regexp.MustCompile(regString)
	return reg.MatchString(fl.Field().String())
}

func iso8601(fl validator.FieldLevel) bool {
	regString := `^(\d{4})-(\d{2})-(\d{2})(T(\d{2}):(\d{2}):(\d{2})(\.\d{0,3})?(Z|[+-](\d{2}):(\d{2})))?$`
	reg := regexp.MustCompile(regString)
	return reg.MatchString(fl.Field().String())
}

func NewULID(t time.Time) string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
