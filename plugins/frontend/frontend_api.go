package frontend

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/sysmanage-web/plugins/acl"
)

func loadFrontendApi(r chi.Router) {
	r.Post("/getRegisteredLinks", func(w http.ResponseWriter, r *http.Request) {
		var reg []Link

		if acl.Enabled() {
			var aclResultCache = make(map[string]bool)

			for _, link := range RegisteredLinks {
				// Cached
				if status, ok := aclResultCache[link.Plugin]; ok {
					if status {
						// Cached as true, append
						reg = append(reg, link)
					}

					// Cached as false, continue
					continue
				}

				userId := r.Header.Get("X-DP-UserID")

				if userId == "" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("SAFETY VIOLATION: user id is unset"))
					return
				}

				req, err := http.NewRequest("GET", "/api/"+link.Plugin+"/@frontend", nil)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Failed to check acls for " + link.Plugin))
					return
				}

				chk := acl.CheckACL(userId, req)

				if chk == nil {
					// Allowed
					reg = append(reg, link)
					aclResultCache[link.Plugin] = true
				} else {
					// Disallowed
					if len(chk.PerPlugin) > 0 || len(chk.PerRoute) > 0 {
						aclResultCache[link.Plugin] = false
					}
				}
			}
		} else {
			reg = RegisteredLinks
		}

		bytes, err := json.Marshal(reg)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not marshal registered links"))
			return
		}

		w.Write(bytes)
	})
}
