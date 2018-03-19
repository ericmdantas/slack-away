package away

import (
	"fmt"
	"strings"
)

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func formatUser(user *string) {
	*user = fmt.Sprintf("<@%s>", *user)

	if *user == "<@>" {
		*user = ""
	}
}
