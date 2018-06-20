package main

import (
	"net/url"
	"regexp"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var validSchemes = map[string]struct{}{
	"":      {},
	"http":  {},
	"https": {},
}

var atomToAttributes = map[atom.Atom][]string{
	atom.A:      {"href"},
	atom.Frame:  {"src"},
	atom.Iframe: {"src"},
	atom.Img:    {"src"},
	atom.Link:   {"href"},
	atom.Script: {"src"},
	atom.Source: {"src", "srcset"},
	atom.Track:  {"src"},
}

type scraper struct {
	excludedPatterns []*regexp.Regexp
	whitespaceRegexp *regexp.Regexp
}

func newScraper(rs []*regexp.Regexp, b bool) scraper {
	r := (*regexp.Regexp)(nil)

	if b {
		r = compileRegexp("(\t|\n|\r)")
	}

	return scraper{rs, r}
}

func (sc scraper) Scrape(n *html.Node, base *url.URL) map[string]error {
	us := map[string]error{}

	for _, n := range scrape.FindAllNested(n, func(n *html.Node) bool {
		_, ok := atomToAttributes[n.DataAtom]
		return ok
	}) {
		for _, a := range atomToAttributes[n.DataAtom] {
			s := scrape.Attr(n, a)

			if s == "" || sc.isURLExcluded(s) {
				continue
			}

			u, err := url.Parse(s)

			if err != nil {
				us[s] = err
				continue
			}

			if _, ok := validSchemes[u.Scheme]; !ok {
				continue
			}

			s = base.ResolveReference(u).String()

			if sc.whitespaceRegexp != nil {
				s = sc.whitespaceRegexp.ReplaceAllString(s, "")
			}

			us[s] = nil
		}
	}

	return us
}

func (sc scraper) isURLExcluded(u string) bool {
	for _, r := range sc.excludedPatterns {
		if r.MatchString(u) {
			return true
		}
	}

	return false
}

func compileRegexp(s string) *regexp.Regexp {
	r, err := regexp.Compile(s)

	if err != nil {
		panic(err)
	}

	return r
}
