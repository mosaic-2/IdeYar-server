package util

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"html/template"
	"math/rand"
	"strconv"
)

type UserIDCtxKey struct{}

func GenerateVerificationCode() string {
	const charset = `ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	//const charset = `0123456789`
	//for front test
	const codeLen = 6
	b := make([]byte, codeLen)
	for i := 0; i < codeLen; i++ {
		b[i] = charset[rand.Int()%len(charset)]
	}
	return string(b)
}

func LoadVerificationEmail(code string) (string, error) {

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

func LoadChangeMailEmail(code string) (string, error) {

	tmp, err := template.ParseFiles("./internal/servicers/util/templates/changeMail.gohtml")
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

func LoadForgetPasswordEmail(code string) (string, error) {

	tmp, err := template.ParseFiles("./internal/servicers/util/templates/forgetPass.gohtml")
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

func GenerateFileName() string {
	return uuid.New().String()
}

func GetUserIDFromCtx(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}

	userID, _ := strconv.ParseInt(md["user-id"][0], 10, 64)
	return int64(userID)
}
