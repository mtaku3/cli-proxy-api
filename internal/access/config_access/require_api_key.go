package configaccess

import (
	"os"
	"strconv"
	"strings"
)

// requireAPIKey reports whether REQUIRE_API_KEY forbids serving the proxy API
// without authentication when no api-keys are configured.
func requireAPIKey() bool {
	enabled, err := strconv.ParseBool(strings.TrimSpace(os.Getenv("REQUIRE_API_KEY")))
	return err == nil && enabled
}
