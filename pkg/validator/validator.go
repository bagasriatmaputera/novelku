package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	v          = validator.New()
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]+$`) //regex based our requirment
)

func init() {
	// Mengambil key asli json -> misal Entity Email yang di return ya json key "email"
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})

	// Costum validasi
	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRe.MatchString(fl.Field().String())
	})
}

func Validate(s any) map[string]string {
	// dari lib validate untuk periksa seluruh field di sebuah struct
	err := v.Struct(s)

	if err == nil {
		return nil
	}

	// apabila bukan error validasi maka beri error default aja
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return map[string]string{"error": err.Error()}
	}

	// wadah untuk error validasi
	out := make(map[string]string, len(verrs))

	for _, fe := range verrs {
		out[fe.Field()] = Message(fe)
	}

	return out
}

// meesage fungsi
func Message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field wajib di isi"
	case "email":
		return "format email tidak valid"
	case "username":
		return "hanya boleh huruf, angka, dan underscore"
	case "min":
		return "minimal" + fe.Param() + "karakter"
	case "max":
		return "maksimal" + fe.Param() + "karakter"
	default:
		return "tidak valid"
	}
}
