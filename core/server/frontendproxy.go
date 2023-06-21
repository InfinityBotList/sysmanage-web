package server

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/infinitybotlist/sysmanage-web/core/state"
	"golang.org/x/exp/slices"
)

func downView(w http.ResponseWriter, r *http.Request, msg string) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(msg))
}

func proxy(w http.ResponseWriter, r *http.Request) {
	var allowedHeaders = []string{
		"Content-Type",
		"Content-Encoding",
		"Content-Security-Policy",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Credentials",
		"Access-Control-Max-Age",
		"Access-Control-Expose-Headers",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Location",
	}

	allowedHeaders = append(allowedHeaders, state.ServerMeta.FrontendServer.ExtraHeadersToAllow...)

	url := state.ServerMeta.FrontendServer.Host + r.URL.Path

	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	// Special case optimization for OPTIONS requests, no need to send/read the body
	if r.Method == "OPTIONS" {
		// Fetch request, no body should be sent
		cli := &http.Client{
			Timeout: 2 * time.Minute,
		}

		req, err := http.NewRequest(r.Method, url, nil)

		if err != nil {
			downView(w, r, "Error creating request to backend")
			return
		}

		req.Header = r.Header

		resp, err := cli.Do(req)

		if err != nil {
			downView(w, r, "Error sending request to backend")
			return
		}

		for k, v := range resp.Header {
			if strings.HasPrefix("X-", k) || slices.Contains(allowedHeaders, k) {
				w.Header()[k] = v
			}
		}

		w.WriteHeader(resp.StatusCode)
		return
	}

	cli := &http.Client{
		Timeout: 2 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(r.Method, url, r.Body)

	if err != nil {
		downView(w, r, "Error creating request to backend")
		return
	}

	req.Header = r.Header

	resp, err := cli.Do(req)

	if err != nil {
		downView(w, r, "Error sending request to backend")
		return
	}

	defer resp.Body.Close()

	for k, v := range resp.Header {
		if strings.HasPrefix(k, "X-") || slices.Contains(allowedHeaders, k) {
			w.Header()[k] = v
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		downView(w, r, "Error reading response body from backend")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(bodyBytes)
}

func startFrontendServer() {
	runCmdArgs := strings.Split(state.ServerMeta.FrontendServer.RunCommand, " ")
	cmd := exec.Command(runCmdArgs[0], runCmdArgs[1:]...)

	var dir string

	if state.ServerMeta.FrontendServer.DirAbsolute {
		dir = state.ServerMeta.FrontendServer.Dir
	} else {
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		if !strings.HasSuffix(cwd, "/") {
			cwd += "/"
		}

		dir = cwd + state.ServerMeta.FrontendServer.Dir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, state.ServerMeta.FrontendServer.ExtraEnv...)

	go func() {
		err := cmd.Run()

		if err != nil {
			panic(err)
		}
	}()
}
