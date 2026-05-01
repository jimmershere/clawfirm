// Command mcp-gateway is the ClawFirm in-house MCP gateway. It mediates
// every tool call from every engine (Dify, n8n, LangGraph, OpenHands, the
// OpenClaw shell) to upstream MCP servers, applying identity, allowlists,
// audit, and the MCP 2025-03-26 token-passthrough prohibition.
//
// See ../README.md for the full design rationale.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/allowlist"
	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/audit"
	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/identity"
	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/policy"
)

func main() {
	listen := flag.String("listen", env("MCP_GW_LISTEN", ":8443"), "TLS listen address")
	cfgPath := flag.String("config", env("MCP_GW_CONFIG", "/etc/mcp-gateway/allowlist.yaml"), "allowlist config")
	auditSink := flag.String("audit-sink", env("MCP_GW_AUDIT_SINK", "http://clawsecure:3188/api/events"), "ClawSecure events endpoint")
	oidcIssuer := flag.String("oidc-issuer", env("MCP_GW_OIDC_ISSUER", ""), "OIDC issuer URL (Authentik)")
	flag.Parse()

	allow, err := allowlist.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load allowlist: %v", err)
	}
	id, err := identity.NewVerifier(*oidcIssuer)
	if err != nil {
		log.Fatalf("init identity: %v", err)
	}
	auditor := audit.NewClawSecureSink(*auditSink)
	pol := policy.New(allow, id, auditor)

	mux := http.NewServeMux()
	mux.Handle("/v1/tools/", pol.Handle())
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "ok")
	}))

	srv := &http.Server{
		Addr:              *listen,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("clawfirm-mcp-gateway listening on %s", *listen)
		// TLS material is mounted at /etc/mcp-gateway/tls/ in production.
		// MVP scaffold uses HTTP for local dev; ListenAndServeTLS in prod.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
