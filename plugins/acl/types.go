package acl

import "net/http"

// Per-route ACLs

// A function that checks if a user is allowed to access a route
type ACLRouteFunc func(d *ACLRouteData) bool

// Defines a ACL entry on a route
type ACLRouteEntry struct {
	Name        string
	Description string
	CheckFunc   ACLRouteFunc
}

type ACLRouteData struct {
	Request *http.Request
	UserID  string
}

// Per-plugin ACLs
type ACLPluginFun func(d *ACLPluginData) bool

type ACLPluginEntry struct {
	Name        string
	Description string
	CheckFunc   ACLPluginFun
}

type ACLPluginData struct {
	Plugin  string
	Request *http.Request
	UserID  string
}

type ACLCheck struct {
	// Per-route ACLs
	PerRoute  []ACLRouteEntry
	PerPlugin []ACLPluginEntry
}
