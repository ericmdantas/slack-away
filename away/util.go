package away

import (
	"fmt"
	"strings"
)

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func formatUser(user string) string {
	userFormatted := fmt.Sprintf("<@%s>", user)

	if userFormatted == "<@>" {
		userFormatted = ""
	}

	return userFormatted
}
