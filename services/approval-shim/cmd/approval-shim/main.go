// Command approval-shim adapts every agent framework's native approval
// primitive to ClawSecure's REST API.
//
// See ../README.md for the design rationale.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jimmershere/clawfirm/services/approval-shim/internal/adapters"
	"google.golang.org/grpc"
)

func main() {
	listen := flag.String("listen", env("SHIM_LISTEN", ":50051"), "gRPC listen address")
	csURL := flag.String("clawsecure-url", env("CLAWSECURE_URL", "http://clawsecure:3188"), "ClawSecure base URL")
	httpListen := flag.String("http-listen", env("SHIM_HTTP_LISTEN", ":50052"), "HTTP listen for n8n/Dify webhook adapters")
	flag.Parse()

	backend := adapters.NewClawSecure(*csURL)

	// HTTP webhook adapters (n8n + Dify)
	mux := http.NewServeMux()
	mux.Handle("/n8n/approve", adapters.N8NHandler(backend))
	mux.Handle("/dify/human-input", adapters.DifyHandler(backend))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	go func() {
		log.Printf("approval-shim http listening on %s", *httpListen)
		if err := http.ListenAndServe(*httpListen, mux); err != nil {
			log.Fatal(err)
		}
	}()

	// gRPC server (LangGraph + OpenHands)
	lis, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	// adapters.RegisterApprovalServer(srv, adapters.NewGRPCServer(backend))
	// (Proto registration lands in Weeks 5-6.)
	log.Printf("approval-shim grpc listening on %s; backend=%s", *listen, *csURL)
	_ = srv.Serve(lis)
}

func env(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}
