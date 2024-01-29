package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Zainal21/my-ewallet/app/appctx"
	"github.com/Zainal21/my-ewallet/app/consts"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/pkg/config"
	"github.com/Zainal21/my-ewallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type SignatureMiddleware struct {
	cnf *config.Config
}

func NewSignatureMiddleware(cnf *config.Config) *SignatureMiddleware {
	return &SignatureMiddleware{
		cnf: cnf,
	}
}

func (t *SignatureMiddleware) verifySignature(xCtx *fiber.Ctx, conf *config.Config) appctx.Response {
	signature := xCtx.Get("X-Signature")
	secretKey := "asololejos"
	payload := string(xCtx.Body())

	if signature == "" {
		return helpers.CreateErrorResponse(fiber.StatusBadRequest, consts.RequestNotValidMessage, nil)
	}

	if xCtx.Method() == "POST" || xCtx.Method() == "PUT" {
		payload := formDataToJSON(payload)
		jwtSignature := helpers.GenerateSignature(secretKey, payload)
		if signature != jwtSignature {
			logger.Info(fmt.Sprintf("Signature verify is not match, Server Signature : %v and Client Signature : %v", jwtSignature, signature))
			return helpers.CreateErrorResponse(fiber.StatusBadRequest, consts.SignatureNotValidMessage, nil)
		}

	}

	return *appctx.NewResponse().WithCode(fiber.StatusOK)
}

func (t *SignatureMiddleware) SignatureVerify(xCtx *fiber.Ctx, conf *config.Config) appctx.Response {
	return t.verifySignature(xCtx, conf)
}

func formDataToJSON(formString string) string {
	keyValues := strings.Split(formString, "&")

	var jsonBuffer strings.Builder
	jsonBuffer.WriteString("{")
	for i, keyValue := range keyValues {
		parts := strings.SplitN(keyValue, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value, err := url.QueryUnescape(parts[1])
			if err != nil {
				logger.Error(fmt.Sprintf("Error decoding value: %v", err))
				return ""
			}

			value = strings.ReplaceAll(value, "\n", "\\n")
			jsonBuffer.WriteString(fmt.Sprintf(`"%s":"%s"`, key, value))

			if i < len(keyValues)-1 {
				jsonBuffer.WriteString(",")
			}
		}
	}
	jsonBuffer.WriteString("}")

	jsonString := jsonBuffer.String()
	return jsonString
}
