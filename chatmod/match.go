package chatmod

import (
	"regexp"
)

func Match(messagePayload string) bool {
	for _, word := range badWordsCache {
		matched, err := regexp.MatchString(word, messagePayload)
		if err != nil || matched {
			return false
		}
	}
	return true
}

func CensorChatMesage(messagePayload string) (string, error) {
	messagePayloadBytes := []byte(messagePayload)
	for _, word := range badWordsCache {
		wordBytes := []byte(word)

		pattern, err := regexp.Compile(string(wordBytes))
		if err != nil {
			return "", err
		}

		replaceBytes := func() []byte {
			out := []byte{}
			for range wordBytes {
				out = append(out, '*')
			}
			return out
		}()
		messagePayloadBytes = pattern.ReplaceAll(messagePayloadBytes, replaceBytes)
	}
	return string(messagePayloadBytes), nil
}
