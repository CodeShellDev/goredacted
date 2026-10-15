package parsers

import (
	"crypto/hmac"
	"crypto/sha256"
	"net"
	"regexp"

	"github.com/codeshelldev/goredacted/internals/parser"
	"github.com/codeshelldev/goredacted/internals/seed"
)

var ipParser = parser.Parser{
	Name: "ip",
	Parse: func(this parser.Parser, content string) string {
		re, err := regexp.Compile(`\b((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\b|(\b(?:[A-Fa-f0-9]{1,4}:){7}[A-Fa-f0-9]{1,4}\b)`)

		if err != nil {
			return content
		}

		return re.ReplaceAllStringFunc(content, func(candidate string) string {
			decision := ShouldProcess(candidate, this)
			if !decision.Process {
				return candidate
			}

			return AnonymizeIP(candidate)
		})
	},
}

func init() {
	parser.Register(ipParser)
}

func AnonymizeIP(ipStr string) string {
    ip := net.ParseIP(ipStr).To4()

    if ip == nil {
        return ""
    }

    mac := hmac.New(sha256.New, []byte(seed.SEED))
    mac.Write([]byte(ipStr))

    hash := mac.Sum(nil)

    anonIP := net.IPv4(hash[0], hash[1], hash[2], hash[3])

    return anonIP.String()
}