package helper

import (
	"fmt"
	"strings"
)

type ExtractVarResult struct {
	Resource string
	ID       string
	Appendix string
}

// checks s for ${}
func ExtractVar(s string) *ExtractVarResult {
	if !strings.Contains(s, "${") && !strings.Contains(s, "}") {
		return nil
	}

	leftPart := strings.Index(s, "${")
	rightPart := strings.Index(s, "}")

	res := s[leftPart+2 : rightPart]

	split := strings.Split(res, ":")
	resource := split[0]
	id := split[1]

	fullSplit := strings.Split(s, fmt.Sprintf("${%s}", res))
	appendix := ""

	if len(fullSplit) > 1 {
		appendix = fullSplit[1]
	}
	return &ExtractVarResult{
		Resource: resource,
		ID:       id,
		Appendix: appendix,
	}
}
