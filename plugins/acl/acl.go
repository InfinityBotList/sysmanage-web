package acl

import (
	"net/http"
	"strings"
	"sync"

	"github.com/infinitybotlist/sysmanage-web/core/plugins/constants"
)

var (
	routes  []ACLRouteEntry
	plugins []ACLPluginEntry
)

// Adds an ACL entry for a route to the list of ACL entries
func AddRoute(e ACLRouteEntry) {
	routes = append(routes, e)
}

// Adds an ACL entry for a plugin to the list of ACL entries
func AddPlugin(e ACLPluginEntry) {
	plugins = append(plugins, e)
}

func CheckACL(userId string, r *http.Request) *ACLCheck {
	// Call every acl route function in parallel and get return values
	// If any of them return false, the user is forbidden
	var wg sync.WaitGroup
	wg.Add(len(routes))

	var flaggedRoutes struct {
		sync.Mutex
		value   bool
		entries []ACLRouteEntry
	}

	for _, route := range routes {
		go func(route ACLRouteEntry) {
			defer wg.Done()
			if !route.CheckFunc(&ACLRouteData{
				Request: r,
				UserID:  userId,
			}) {
				defer flaggedRoutes.Unlock()
				flaggedRoutes.Lock()
				flaggedRoutes.value = true
				flaggedRoutes.entries = append(flaggedRoutes.entries, route)
			}
		}(route)
	}

	wg.Wait()

	if !flaggedRoutes.value {
		return &ACLCheck{
			PerRoute: flaggedRoutes.entries,
		}
	}

	// Next check per-plugin
	if !strings.HasPrefix(r.URL.Path, "/api/") {
		// Not a plugin
		return nil
	}

	// Get plugin name
	if len(strings.Split(r.URL.Path, "/")) < 3 {
		// Not a plugin
		return nil
	}

	pluginName := strings.Split(r.URL.Path, "/")[2]

	// Call every acl plugin function in parallel and get return values
	// If any of them return false, the user is forbidden

	var flaggedPlugins struct {
		sync.Mutex
		value   bool
		entries []ACLPluginEntry
	}

	for _, plugin := range plugins {
		go func(plugin ACLPluginEntry) {
			defer wg.Done()
			if !plugin.CheckFunc(&ACLPluginData{
				Plugin:  pluginName,
				Request: r,
				UserID:  userId,
			}) {
				defer flaggedPlugins.Unlock()
				flaggedPlugins.Lock()
				flaggedPlugins.value = true
				flaggedPlugins.entries = append(flaggedPlugins.entries, plugin)
			}
		}(plugin)
	}

	wg.Wait()

	if !flaggedPlugins.value {
		return &ACLCheck{
			PerPlugin: flaggedPlugins.entries,
		}
	}

	return nil
}

func MuxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get(constants.UserIdHeader)

		if userId == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("SAFETY VIOLATION: user id is unset"))
			return
		}

		chk := CheckACL(userId, r)

		if chk != nil {
			w.WriteHeader(http.StatusForbidden)

			if len(chk.PerRoute) > 0 {
				var aclentries []string

				for _, entry := range chk.PerRoute {
					aclentries = append(aclentries, entry.Name)
				}

				w.Write([]byte("Forbidden by per-route ACL:" + strings.Join(aclentries, ", ")))
				return
			}

			if len(chk.PerPlugin) > 0 {
				var aclentries []string

				for _, entry := range chk.PerPlugin {
					aclentries = append(aclentries, entry.Name)
				}

				w.Write([]byte("Forbidden by per-plugin ACL:" + strings.Join(aclentries, ", ")))
				return
			}

			panic("SAFETY VIOLATION: acl check failed with no entries returned")
		}

		next.ServeHTTP(w, r)
	})
}
