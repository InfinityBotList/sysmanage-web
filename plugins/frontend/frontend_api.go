package frontend

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/plugins/acl"
)

// Gets the registered links on the site taking ACLS etc into account
func GetRegisteredLinks(r *http.Request) ([]Link, error) {
	var reg []Link

	if plugins.Enabled("acl") {
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
				return nil, errors.New("user id is unset")
			}

			req, err := http.NewRequest("GET", "/api/"+link.Plugin+"/@frontend", nil)

			if err != nil {
				return nil, errors.New("failed to check acls for " + link.Plugin)
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

	return reg, nil
}

func loadFrontendApi(r chi.Router) {
	r.Post("/getRegisteredLinks", func(w http.ResponseWriter, r *http.Request) {
		reg, err := GetRegisteredLinks(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("SAFETY VIOLATION: " + err.Error()))
			return
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
