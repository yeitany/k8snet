package utils

import "github.com/goccy/go-graphviz"

var allowedFormat []graphviz.Format = []graphviz.Format{
	graphviz.XDOT,
	graphviz.SVG,
	graphviz.PNG,
	graphviz.JPG,
}

var allowedLayout []graphviz.Layout = []graphviz.Layout{
	graphviz.CIRCO,
	graphviz.DOT,
	graphviz.FDP,
	graphviz.NEATO,
	graphviz.OSAGE,
	graphviz.PATCHWORK,
	graphviz.SFDP,
	graphviz.TWOPI,
}

type UrlParmeters struct {
	BlacklistNamespaces []string        `json:"blacklist_namespaces"`
	WhitelistNamespaces []string        `json:"whitelist_namespaces"`
	Targets             []string        `json:"targets"`
	Format              graphviz.Format `json:"format"`
	Layout              graphviz.Layout `json:"layout"`
}

func (u *UrlParmeters) isLayoutValid() bool {
	for i := range allowedLayout {
		if u.Layout == allowedLayout[i] {
			return true
		}
	}
	return false
}

func (u *UrlParmeters) isFormatValid() bool {
	for i := range allowedFormat {
		if u.Format == allowedFormat[i] {
			return true
		}
	}
	return false
}

func (u *UrlParmeters) IsValid() bool {
	return u.isFormatValid() && u.isLayoutValid()
}

type CtxKey string
