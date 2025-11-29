package parsers

import (
	"crypto/sha256"
	"regexp"
	"slices"

	"github.com/codeshelldev/goredacted/internals/config"
	"github.com/codeshelldev/goredacted/internals/parser"
)

func Match(input string, targets []string) string {
	hash := sha256.Sum256([]byte(input))

	index := int(hash[0]) << 24 | int(hash[1]) << 16 | int(hash[2]) << 8 | int(hash[3])
	if index < 0 {
		index = -index
	}

	return targets[index % len(targets)]
}

func Ingore(match string, parser parser.Parser) bool {
	if len(config.CONFIG.EXCLUDE) <= 0 && len(config.CONFIG.INCLUDE) <= 0 {
		return false
	}

	cleanFunc := func(xxcludes []config.XXCLUDE) []config.XXCLUDE {
		i := 0

		for _, xxclude := range xxcludes {
			if xxclude.PARSER == parser.Name {
				xxcludes[i] = xxclude
				i++
			}
		}
		return xxcludes[:i]
	}

	included := cleanFunc(config.CONFIG.INCLUDE)
	excluded := cleanFunc(config.CONFIG.EXCLUDE)

	containsFunc := func(try config.XXCLUDE) bool {
		re, err := regexp.Compile(try.REGEX)

		if err != nil {
			return false
		}

		return re.MatchString(match)
	}

	isIncluded := slices.ContainsFunc(config.CONFIG.INCLUDE, containsFunc)
	isExcluded := slices.ContainsFunc(config.CONFIG.EXCLUDE, containsFunc)

	if isIncluded {
		return false
	}
	
	if isExcluded {
		return true
	}

	if len(included) > 0 && len(excluded) == 0 {
		return true
	}

	if len(excluded) > 0 && len(included) == 0 {
		return false
	}

	return true
}