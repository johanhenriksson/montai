package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/johanhenriksson/montai/mount"
	"github.com/johanhenriksson/montai/s3fs"
)

type server struct {
	*mount.Manager
}

func Serve() {
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
	m := &server{manager}
	http.HandleFunc("/mount", requestWrapper(m.handleMount))
	http.HandleFunc("/unmount", requestWrapper(m.handleUnmount))
	http.HandleFunc("/mounts", m.handleGetMounts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *server) handleMount(req *mount.Request) (interface{}, error) {
	return s.Mount(req)
}

func (s *server) handleGetMounts(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, s.Mounts)
}

func (s *server) handleUnmount(req *mount.Request) (interface{}, error) {
	return s.Unmount(req.Name)
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

func isContainer() bool {
	path := "/proc/1/cgroup"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "docker")
}
