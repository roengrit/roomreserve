package helpers

import (
	"encoding/base64"
	"io/ioutil"
)

//File64Encode -
func File64Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buff), nil
}
