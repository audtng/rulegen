package main

	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"


func (s *HTTPServer) handlePostFile(w http.ResponseWriter, r *http.Request) {
	u := s.charmUserFromRequest(w, r)
	path := filepath.Clean(pattern.Path(r.Context()))
	ms := r.URL.Query().Get("mode")
	m, err := strconv.ParseUint(ms, 10, 32)
	if err != nil {

func (s *HTTPServer) handleGetFile(w http.ResponseWriter, r *http.Request) {
	u := s.charmUserFromRequest(w, r)
	path := filepath.Clean(pattern.Path(r.Context()))
	f, err := s.cfg.FileStore.Get(u.CharmID, path)
	if errors.Is(err, fs.ErrNotExist) {
		s.renderCustomError(w, "file not found", http.StatusNotFound)

func (s *HTTPServer) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	u := s.charmUserFromRequest(w, r)
	path := filepath.Clean(pattern.Path(r.Context()))
	err := s.cfg.FileStore.Delete(u.CharmID, path)
	if err != nil {
		log.Printf("cannot delete file: %s", err)
