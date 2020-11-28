package app

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// checkAnswerForSelect 检查为合法的 [选择题] 答案
func checkAnswerForSelect(answer string) (ok bool) {

	ok, err := regexp.MatchString("^[A-Za-z]$", answer)
	if err != nil {
		log.Error().Err(err).Msg("检查")
	}

	return
}

// checkAnswerForFill 检查为合法的 [简答题] 答案
func checkAnswerForFill(answer string) bool { return strings.HasPrefix(answer, "#") }

// HashForSHA1 SHA1 散列
func HashForSHA1(d string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(d))
	return hex.EncodeToString(h.Sum(nil))
}

func DoWordSplit(s string) (words []string, err error) {
	words, err = AC.Cli.C.GetWordSegmentation(s)

	if err != nil {
		return
	}

	for k, v := range words {
		words[k] = strings.ReplaceAll(v, "\u0000", "")
	}

	return words, nil
}
