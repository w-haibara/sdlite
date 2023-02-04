package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

func ListenAndServe(addr string) error {
	http.Handle("/service/", new(serviceHandler))

	return http.ListenAndServe(addr, nil)
}

type serviceHandler struct {
	m sync.Map
}

type hosts [2]string

func (hs hosts) isEmpty() bool {
	for _, v := range hs {
		if v != "" {
			return false
		}
	}

	return true
}

func (hs hosts) toSlice() []string {
	res := make([]string, 0, len(hs))
	for _, v := range hs {
		res = append(res, v)
	}

	return res
}

func (hs hosts) current() string {
	return hs[0]
}

func (hs hosts) prev() string {
	return hs[1]
}

func (h *serviceHandler) registerServiceHost(name, host string) {
	hs := h.fetchServiceHosts(name)

	hs[1] = hs[0]
	hs[0] = host

	h.m.Store(name, hs)
}

func (h *serviceHandler) fetchServiceHosts(name string) hosts {
	v, ok := h.m.Load(name)
	if !ok {
		return hosts{}
	}

	hs := v.(hosts)
	return hs
}

func (h *serviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/service/")
	switch r.Method {
	case http.MethodGet:
		h.serveHTTPGet(name, w, r)
	case http.MethodPost:
		h.serveHTTPPost(name, w, r)
	}
}

func (h *serviceHandler) serveHTTPGet(name string, w http.ResponseWriter, r *http.Request) {
	hs := h.fetchServiceHosts(name)
	if hs.isEmpty() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	type result struct {
		Name        string
		Hosts       []string
		CurrentHost string
		PrevHost    string
	}
	res := result{
		Name:        name,
		Hosts:       hs.toSlice(),
		CurrentHost: hs.current(),
		PrevHost:    hs.prev(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("[serveHTTPGet] failed to encode result:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *serviceHandler) serveHTTPPost(name string, w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	type request struct {
		Host string
	}

	req := new(request)
	if err := dec.Decode(req); err != nil {
		log.Println("[serveHTTPPost] failed to decode request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.registerServiceHost(name, req.Host)
}
