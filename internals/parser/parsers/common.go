package parsers

import (
	"crypto/sha256"
	"encoding/binary"
	"regexp"

	"github.com/codeshelldev/goredacted/internals/config"
	"github.com/codeshelldev/goredacted/internals/parser"
)

func SelectReplacement(input string, targets []string) string {
	if len(targets) == 0 {
		return ""
	}

	hash := sha256.Sum256([]byte(input))
	index := binary.BigEndian.Uint64(hash[:8]) % uint64(len(targets))

	return targets[index]
}

type MatchDecision struct {
	Process bool
	Matched string
}

func ParserRules(name string, entries []config.RegexParserEntry) []config.RegexParserEntry {
	var rules []config.RegexParserEntry

	for _, entry := range entries {
		if entry.Parser == name {
			rules = append(rules, entry)
		}
	}

	return rules
}

func ShouldProcess(match string, p parser.Parser) MatchDecision {
	included := ParserRules(p.Name, config.CONFIG.Include)
	excluded := ParserRules(p.Name, config.CONFIG.Exclude)

	if len(included) == 0 && len(excluded) == 0 {
		return MatchDecision{ Process: true }
	}

	matchesAny := func(rules []config.RegexParserEntry) string {
		for _, rule := range rules {
			re, err := regexp.Compile(rule.Regex)
			if err != nil {
				continue
			}

			matched := re.FindString(match)
			if matched != "" {
				return matched
			}
		}

		return ""
	}

	matched := matchesAny(excluded)
	if matched != "" {
		return MatchDecision{ Process: false, Matched: matched }
	}

	if len(included) > 0 {
		matched := matchesAny(included)
		return MatchDecision{
			Process: matched != "",
			Matched: matched,
		}
	}

	return MatchDecision{ Process: true }
}

func CustomParserPatterns(name string) []string {
	rules := ParserRules(name, config.CONFIG.Include)

	patterns := make([]string, 0, len(rules))
	for _, rule := range rules {
		_, err := regexp.Compile(rule.Regex)
		if err != nil {
			continue
		}

		patterns = append(patterns, rule.Regex)
	}

	return patterns
}