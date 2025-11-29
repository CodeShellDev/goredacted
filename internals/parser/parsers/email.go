package parsers

import (
	"regexp"
	"strings"

	"github.com/codeshelldev/goredacted/internals/parser"
)

var emailParser = parser.Parser{
	Name: "email",
	Parse: func(this parser.Parser, content string) string {
		re, err := regexp.Compile(`([a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `{|}~-]+)@(?:\b([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+)\b)`)

		if err != nil {
			return content
		}

		return re.ReplaceAllStringFunc(content, func(match string) string {
			if Ingore(match, this) {
				return match
			}

			submatches := re.FindStringSubmatch(match)

			local := submatches[1]
			domain := submatches[2]

			parts := strings.Split(local, ".")

			localMasks := []string{}

			for i, part := range parts {
				if part == "" {
					continue
				}

				var mask string

				useFirstName := (i % 2 == 0)

				if useFirstName {
					mask = Match(local, FIRST_NAMES)
				} else {
					mask = Match(local, LAST_NAMES)
				}

				localMasks = append(localMasks, mask)
			}

			return strings.Join(localMasks, ".") + "@" + domain
		})
	},
}

func init() {
	parser.Register(emailParser)
}