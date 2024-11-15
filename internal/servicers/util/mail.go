package util

import (
	"bytes"
	"html/template"
)

func verificationEmail(code string) (string, error) {

	tmp, err := template.ParseFiles("./internal/servicers/util/templates/verification.gohtml")
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)

	err = tmp.Execute(b, code)
	if err != nil {
		return "", err
	}

	message := string(b.Bytes())

	return message, nil
}
