package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/johanhenriksson/montai/mount"
	"github.com/johanhenriksson/montai/s3fs"
)

type app struct {
	*mount.Manager
}

func main() {
	log.Println("hello team")
	log.Println("montai v0.1.0")

	// settings
	root := "/montai"
	if setroot := os.Getenv("ROOT_DIR"); setroot != "" {
		root = setroot
	}
	log.Println("mount root:", root)

	// create montai manager
	manager := mount.NewManager(root)
	manager.AddProvider(s3fs.New())

	// web server
	m := &app{manager}
	http.HandleFunc("/mount", requestWrapper(m.handleMount))
	http.HandleFunc("/unmount", requestWrapper(m.handleUnmount))
	http.HandleFunc("/mounts", m.handleGetMounts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (m *app) handleMount(req *mount.Request) (interface{}, error) {
	return m.Mount(req)
}

func (m *app) handleGetMounts(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, m.Mounts)
}

func (m *app) handleUnmount(req *mount.Request) (interface{}, error) {
	return m.Unmount(req.Name)
}

type reqHandler func(*mount.Request) (interface{}, error)

func requestWrapper(handler reqHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mount.Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := handler(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, result)
	}
}

func jsonResponse(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)
}
