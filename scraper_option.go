package main

import "regexp"

type scraperOptions struct {
	ExcludedPatterns      []*regexp.Regexp
	ignoreWhitespaceInURL bool
}
