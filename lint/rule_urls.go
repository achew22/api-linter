package lint

import "strings"

// A list of mapping functions, each of which returns the rule URL for
// the given rule name, and if not found, return an empty string.
//
// At Google, we inject additional rule URL mappings into this list.
// Example: google_rule_url_mappings.go
// package lint
// func init() {
//   ruleURLMappings = append(ruleURLMappings, internalRuleURLMapping)
// }
//
// func internalRuleURLMapping(ruleName string) string {
//   ...
// }
var ruleURLMappings = []func(string) string{
	coreRuleURL,
	clientLibrariesRuleUrl,
	cloudRuleUrl,
}

func coreRuleURL(ruleName string) string {
	return groupUrl(ruleName, "core")
}

func clientLibrariesRuleUrl(ruleName string) string {
	return groupUrl(ruleName, "client-libraries")
}

func cloudRuleUrl(ruleName string) string {
	return groupUrl(ruleName, "cloud")
}

func groupUrl(ruleName, groupName string) string {
	base := "https://linter.aip.dev/"
	nameParts := strings.Split(ruleName, "::") // e.g., client-libraries::0122::camel-case-uris -> ["client-libraries", "0122", "camel-case-uris"]
	if len(nameParts) == 0 || nameParts[0] != groupName {
		return ""
	}
	path := strings.TrimPrefix(strings.Join(nameParts[1:], "/"), "0")
	return base + path
}

func getRuleURL(ruleName string, nameURLMappings []func(string) string) string {
	for i := len(nameURLMappings) - 1; i >= 0; i-- {
		if url := nameURLMappings[i](ruleName); url != "" {
			return url
		}
	}
	return ""
}
