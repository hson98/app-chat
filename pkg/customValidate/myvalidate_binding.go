package customvalidate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func RegisterValidate() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notspace", func(fl validator.FieldLevel) bool {
			//match:=strings.Split(fl.Param()," ")
			// convert field value to string
			value := fl.Field().String()
			whitespace := regexp.MustCompile(`\s`).MatchString(value)
			if whitespace {
				return false
			}
			return true
		})
	}
}
