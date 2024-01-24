package zenValidator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/id"
)

type (
	Validator interface {
		Validate(interface{}) error
	}

	implementation struct {
		instance   *validator.Validate
		translator *ut.UniversalTranslator
	}
)

func New(opts ...ValidatorOption) (Validator, error) {
	options := defaultValidatorOptions()
	for _, opt := range opts {
		opt.Apply(&options)
	}

	enTranslator, _ := options.translator.GetTranslator("en")
	idTranslator, _ := options.translator.GetTranslator("id")

	validate := validator.New()
	if err := en.RegisterDefaultTranslations(validate, enTranslator); err != nil {
		return nil, err
	}
	if err := id.RegisterDefaultTranslations(validate, idTranslator); err != nil {
		return nil, err
	}

	instance := &implementation{instance: validate, translator: options.translator}
	return instance, nil
}

func (i *implementation) Validate(object interface{}) error {
	if err := i.instance.Struct(object); err != nil {
		return err
	}
	return nil
}
