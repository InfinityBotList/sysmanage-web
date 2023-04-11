package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sysmanage-web/types"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

//go:embed all:frontend/build
var frontend embed.FS

var (
	config *types.Config
	rdb    *redis.Client
	ctx    = context.Background()
	v      *validator.Validate

	//go:embed data/servicegen/service.tmpl
	serviceTemplate string

	//go:embed data/servicegen/target.tmpl
	targetTemplate string
)

func ensureDpAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/__external__/") || config.DPDisable {
			// External route or under DPDisable, skip auth
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("X-DP-UserID") != "" {
			// Check if user is allowed
			for _, user := range config.AllowedUsers {
				if user == r.Header.Get("X-DP-UserID") {
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
					h := hmac.New(sha512.New, []byte(config.DPSecret))
					h.Write([]byte(ts))
					h.Write([]byte(r.Header.Get("X-DP-UserID")))
					hexed := hex.EncodeToString(h.Sum(nil))

					if r.Header.Get("X-DP-Signature") != hexed {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Unauthorized. Signature from deployproxy mismatch"))
						return
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

					// User is allowed
					next.ServeHTTP(w, r)
					return
				}
			}

			// User is not allowed
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized, user not in allowlist"))
			return
		} else {
			// User is not authenticated
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. Not running under deployproxy?"))
			return
		}
	})
}

func routeStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			// Serve frontend
			serverRoot, err := fs.Sub(frontend, "frontend/build")
			if err != nil {
				log.Fatal(err)
			}

			// Get file extension
			fileExt := ""
			if strings.Contains(r.URL.Path, ".") {
				fileExt = r.URL.Path[strings.LastIndex(r.URL.Path, "."):]
			}

			if fileExt == "" && r.URL.Path != "/" {
				r.URL.Path += ".html"
			}

			http.FileServer(http.FS(serverRoot)).ServeHTTP(w, r)
		} else {
			// Serve API
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	// Load config.yaml into Config struct
	file, err := os.Open("config.yaml")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)

	if err != nil {
		panic(err)
	}

	// Connect to redis
	rOptions, err := redis.ParseURL(config.RedisURL)

	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(rOptions)

	v = validator.New()

	// Create wildcard route
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(
		middleware.Recoverer,
		middleware.Logger,
		middleware.CleanPath,
		middleware.RealIP,
		ensureDpAuth,
		routeStatic,
		middleware.Timeout(30*time.Second),
	)

	loadApi(r)

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		w.Write([]byte("API endpoint not found..."))
	})

	// Create server
	s := &http.Server{
		Addr:    ":30010",
		Handler: r,
	}

	// Start server
	fmt.Println("Starting server on port 30010")
	err = s.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
