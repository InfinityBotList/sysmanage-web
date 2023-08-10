package authdp

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/infinitybotlist/sysmanage-web/core/plugins/constants"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"golang.org/x/exp/slices"
)

func DpAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Auth-exempt routes should be excluded
		if slices.Contains(state.AuthExemptRoutes, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("X-DP-Host") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Host header not found."))
			return
		}

		if r.Header.Get("X-DP-Host") != url {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. Domain rebind detected. Expected " + url + " but got " + r.Header.Get("X-DP-Host")))
			return
		}

		if r.Header.Get("X-DP-UserID") == "" {
			// User is not authenticated
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. Not running under deployproxy?"))
			return
		}

		// User is allowed, set constants.UserIdHeader to user id for other plugins to use it
		r.Header.Set(constants.UserIdHeader, r.Header.Get("X-DP-UserID"))

		fmt.Println(r.Header.Get(constants.UserIdHeader), allowedUsers)

		// Check if user is allowed
		if len(allowedUsers) > 0 && !slices.Contains(allowedUsers, r.Header.Get(constants.UserIdHeader)) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. User not allowed to access this site."))
			return
		}

		// User is possibly allowed
		if r.Header.Get("X-DP-Signature") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Signature header not found."))
			return
		}

		// Check for X-DP-Timestamp
		if r.Header.Get("X-DP-Timestamp") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Timestamp header not found."))
			return
		}

		ts := r.Header.Get("X-DP-Timestamp")

		// Validate DP-Secret next
		if dpSecret != "" {
			h := hmac.New(sha512.New, []byte(dpSecret))
			h.Write([]byte(ts))
			h.Write([]byte(r.Header.Get("X-DP-UserID")))
			hexed := hex.EncodeToString(h.Sum(nil))

			if r.Header.Get("X-DP-Signature") != hexed {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized. Signature from deployproxy mismatch"))
				return
			}
		}

		// Check if timestamp is valid
		timestamp, err := strconv.ParseInt(ts, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Timestamp is not a valid integer."))
			return
		}

		if time.Now().Unix()-timestamp > 10 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Timestamp is too old."))
			return
		}

		next.ServeHTTP(w, r)
	})
}
