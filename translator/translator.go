package zenTranslator

import (
	"sync"

	ut "github.com/go-playground/universal-translator"
)

var (
	translator *ut.UniversalTranslator
	once       sync.Once
)

func GetTranslator(opts ...Option) *ut.UniversalTranslator {
	once.Do(func() {
		options := defaultOption()
		for _, opt := range opts {
			opt.Apply(&options)
		}
		translator = ut.New(options.defaultLanguage, options.supportedLanguage...)
	})
	return translator
}
