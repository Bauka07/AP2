package utils

import (
	"fmt"
	"strings"
)

func GetKey(path string) (string, error){
	key := strings.TrimPrefix(path, "/data/")
	key =	strings.Trim(key, "/")
	if key == "" {
		return "",	fmt.Errorf("missing id")
	}

	return key, nil
}
