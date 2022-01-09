package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/golang/glog"
)

func main() {
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/debug/pprof/", pprof.Index)
	http.HandleFunc("/debug/pprof/profile", pprof.Profile)
	http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	http.HandleFunc("/debug/pprof/trace", pprof.Trace)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	os.Setenv("VERSION", "0.0.1")
	version := os.Getenv("VERSION")
	fmt.Println(version)
	w.Header().Set("VERSION", version)

	//X-REAL-IP
	//X-FORWARD-FOR
	clientIP := getCurrentIP(r)
	httpCode := http.StatusOK

	fmt.Printf("client ip is %s", clientIP)
	fmt.Printf("http status is %d", httpCode)

}

func getCurrentIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")

	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "working")
}
