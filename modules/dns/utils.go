package vsdns

import (
	"net"
	"regexp"
	"strings"
)

func parseRecords(userRecords []Record) Records {
	var records Records
	records.A = make(map[string][]string)
	records.AAAA = make(map[string][]string)
	records.CNAME = make(map[string][]string)
	records.MX = make(map[string][]string)
	records.NS = make(map[string][]string)
	records.TXT = make(map[string][]string)

	for _, record := range userRecords {
		switch record.Type {
		case "A":
			name := strings.ToLower(strings.TrimSpace(record.Name))
			content := strings.TrimSpace(record.Content)
			// check if value is IPv4
			if isIP(content) {
				records.A[name] = append(records.A[name], content)
			}
		case "AAAA":
			name := strings.ToLower(strings.TrimSpace(record.Name))
			content := strings.TrimSpace(record.Content)
			// check if value is IPv6
			if isIP(content) {
				records.AAAA[name] = append(records.AAAA[name], content)
			}
		case "CNAME":
			name := strings.ToLower(strings.TrimSpace(record.Name))
			content := strings.TrimSpace(record.Content)
			if isDomain(content) {
				records.CNAME[name] = append(records.CNAME[name], content)
			}
		case "MX":
			name := strings.ToLower(strings.TrimSpace(record.Name))
			content := strings.TrimSpace(record.Content)
			records.MX[name] = append(records.MX[name], content)
		case "NS":
			name := strings.ToLower(strings.TrimSpace(record.Name))
			content := strings.TrimSpace(record.Content)
			records.NS[name] = append(records.NS[name], content)
		case "TXT":
			name := strings.ToLower(strings.TrimSpace(record.Name))
			content := strings.TrimSpace(record.Content)
			records.TXT[name] = append(records.TXT[name], content)
		}
	}
	return records
}

// isDomain returns true if string is a valid domain
func isDomain(domain string) bool {
	if isIP(domain) {
		return false
	}
	match, _ := regexp.MatchString(`^([a-zA-Z0-9\*]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,9}$`, domain)
	return match
}

// isIP returns true if string IP address
func isIP(ip string) bool {
	return (net.ParseIP(ip) != nil)
}

// match returns true if wildcard value is mathed
func match(name string, value string) bool {
	var result strings.Builder
	for i, literal := range strings.Split(name, "*") {
		if i > 0 {
			result.WriteString(".*")
		}
		result.WriteString(regexp.QuoteMeta(literal))
	}
	matched, _ := regexp.MatchString(result.String(), value)
	return matched
}
