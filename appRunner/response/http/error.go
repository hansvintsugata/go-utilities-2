package zenHttpResponse

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	zenMandatory "github.com/hansvintsugata/go-utilities-2/mandatory"
	zenTranslator "github.com/hansvintsugata/go-utilities-2/translator"
)

type (
	errorResponse struct {
		Error *Error `json:"error,omitempty"`
	}
	Error struct {
		Message string       `json:"message"`
		Code    int          `json:"code"`
		Errors  []ErrorField `json:"errors,omitempty"`
	}
	ErrorField struct {
		Field  string `json:"field"`
		Reason string `json:"reason"`
	}
)

func (e *errorResponse) Val() interface{} {
	return e
}

func newErrorResponse(ctx context.Context, responseCode int, err error, errLocalization map[string]map[string]string) Response {
	var errStruct *Error = nil
	if err != nil {
		eMsg, e := getErrors(ctx, responseCode, err, errLocalization)
		errStruct = &Error{
			Message: eMsg,
			Code:    responseCode,
			Errors:  e,
		}
	}
	return &errorResponse{Error: errStruct}
}

func getErrors(ctx context.Context, responseCode int, err error, errLocalization map[string]map[string]string) (string, []ErrorField) {
	if responseCode != http.StatusBadRequest {
		return err.Error(), nil
	}

	x, ok := err.(validator.ValidationErrors)
	if ok {
		data := buildCustomErrors(ctx, x, errLocalization)
		return data[0].Reason, data
	}
	return err.Error(), nil
}

func buildCustomErrors(ctx context.Context, errs validator.ValidationErrors, errLocalization map[string]map[string]string) []ErrorField {
	mandatory := zenMandatory.FromContext(ctx)
	translator, _ := zenTranslator.GetTranslator().GetTranslator(mandatory.Language())
	errors := make([]ErrorField, 0)
	for _, err := range errs {
		if len(errLocalization) > 0 && len(errLocalization[mandatory.Language()]) > 0 {
			if x, exist := errLocalization[mandatory.Language()][fmt.Sprintf("%s.%s", err.StructNamespace(), err.Tag())]; exist {
				errors = append(errors, ErrorField{
					Field:  err.Field(),
					Reason: x,
				})
				continue
			} else if x, exist := errLocalization[mandatory.Language()][err.StructNamespace()]; exist {
				errors = append(errors, ErrorField{
					Field:  err.Field(),
					Reason: x,
				})
				continue
			}
		}

		errors = append(errors, ErrorField{
			Field:  err.Field(),
			Reason: err.Translate(translator),
		})
	}
	return errors
}
