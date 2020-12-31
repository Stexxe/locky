//go:generate go run gen/generator.go

package resources

import (
	"errors"
	"fmt"
)

func Get(name string) (string, error) {
	b, err := GetBin(name)

	if err != nil {
		return "", err
	}

	return string(*b), nil
}

func GetBin(name string) (*[]byte, error) {
	b, ok := resources[name]

	if !ok {
		return &[]byte{}, errors.New(fmt.Sprintf("Embedded file named %s not found", name))
	}

	return b, nil
}


