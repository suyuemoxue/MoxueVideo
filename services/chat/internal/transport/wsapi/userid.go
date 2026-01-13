package wsapi

import (
	"net/http"
	"strconv"
	"strings"
)

func ExtractUserID(r *http.Request) (uint64, bool) {
	raw := strings.TrimSpace(r.URL.Query().Get("uid"))
	if raw == "" {
		raw = strings.TrimSpace(r.Header.Get("X-User-ID"))
	}
	if raw == "" {
		raw = strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
	}
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}
