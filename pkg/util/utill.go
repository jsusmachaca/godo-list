package util

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

var ErrInvalidDataType = errors.New("invalid data type")

func ReadJson(filestring string, model any) error {
	file, err := os.OpenFile(filestring, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(model)
	if err != nil {
		return err
	}

	return nil
}

func WriteJson(filestring string, model any) error {
	file, err := os.OpenFile(filestring, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(model)

	return nil
}

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
