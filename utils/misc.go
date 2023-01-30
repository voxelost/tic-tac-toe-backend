package utils

import (
	"main/chatmod"
	"strings"
)

func PreprocessChatPayload(payload string) (preprocessedPayload string, shouldForward bool) {
	defer recover()

	censoredPayload, err := chatmod.CensorChatMesage(payload)
	if err != nil {
		return "", false
	}

	censoredPayload = strings.TrimSpace(censoredPayload)

	if len(censoredPayload) > 200 {
		censoredPayload = censoredPayload[:200]
	} else if len(censoredPayload) < 1 {
		return "", false
	}

	return censoredPayload, true
}
