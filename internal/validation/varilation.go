package validation

import (
	"encoding/json"
	"errors"
	"io"
)

func RequestValidator(body io.ReadCloser, model any) error {
	err := json.NewDecoder(body).Decode(model)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			return ErrInvalidDataType
		}
		return err
	}

	return nil
}
