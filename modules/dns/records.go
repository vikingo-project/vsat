package vsdns

import (
	"sort"
	"strings"
)

type Record struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Records struct {
	A     map[string][]string
	AAAA  map[string][]string
	CNAME map[string][]string
	MX    map[string][]string
	NS    map[string][]string
	TXT   map[string][]string
}

func getNames(recs map[string][]string) []string {
	names := make([]string, 0, len(recs))
	for name := range recs {
		names = append(names, name)
	}
	// sort by pattern length; from long to short
	sort.Slice(names, func(i, j int) bool {
		return len(names[i]) > len(names[j])
	})
	return names
}

func (rr *Records) LookupA(domain string) []string {
	domain = strings.ToLower(domain)
	names := getNames(rr.A)
	for _, name := range names {
		if match(name, domain) {
			return rr.A[name]
		}
	}
	return []string{}
}

func (rr *Records) LookupAAAA(domain string) []string {
	domain = strings.ToLower(domain)
	names := getNames(rr.AAAA)
	for _, name := range names {
		if match(name, domain) {
			return rr.AAAA[name]
		}
	}
	return []string{}
}
func (rr *Records) LookupCNAME(domain string) string {
	domain = strings.ToLower(domain)
	names := getNames(rr.CNAME)
	for _, name := range names {
		if match(name, domain) {
			return rr.CNAME[name][0]
		}
	}
	return ""
}

func (rr *Records) LookupNS(domain string) []string {
	domain = strings.ToLower(domain)
	names := getNames(rr.NS)
	for _, name := range names {
		if match(name, domain) {
			return rr.NS[name]
		}
	}
	return []string{}
}

func (rr *Records) LookupMX(domain string) []string {
	domain = strings.ToLower(domain)
	names := getNames(rr.MX)
	for _, name := range names {
		if match(name, domain) {
			return rr.MX[name]
		}
	}
	return []string{}
}

func (rr *Records) LookupTXT(domain string) []string {
	domain = strings.ToLower(domain)
	names := getNames(rr.TXT)
	for _, name := range names {
		if match(name, domain) {
			return rr.TXT[name]
		}
	}
	return []string{}
}
