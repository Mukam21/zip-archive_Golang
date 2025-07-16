package service

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"
	"zip-archive_Golang/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HandlerTk struct {
	Cfg *model.Config

	TasksWk map[string]*model.Task
	Mu      sync.Mutex

	ActiveSlot chan struct{}
}

func NewTkHandler(cfg *model.Config) *HandlerTk {
	return &HandlerTk{

		Cfg:     cfg,
		TasksWk: make(map[string]*model.Task),

		ActiveSlot: make(chan struct{}, cfg.MaxParallel),
	}
}

func (h *HandlerTk) CreateTk(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewString()

	h.Mu.Lock()

	h.TasksWk[id] = &model.Task{ID: id, Status: model.PendingTk, CreatedAt: time.Now()}

	h.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *HandlerTk) AddFileFlow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)

		return
	}

	h.Mu.Lock()

	task, ok := h.TasksWk[id]

	h.Mu.Unlock()

	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)

		return
	}

	task.Mu.Lock()

	defer task.Mu.Unlock()

	if len(task.URLs) >= h.Cfg.MaxFileTask {
		http.Error(w, "maximum files per task reached", http.StatusBadRequest)

		return
	}

	task.URLs = append(task.URLs, req.URL)

	if len(task.URLs) == h.Cfg.MaxFileTask {
		go h.processTk(task)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HandlerTk) GetStatusWk(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.Mu.Lock()

	task, ok := h.TasksWk[id]

	h.Mu.Unlock()

	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)

		return
	}
	task.Mu.Lock()
	defer task.Mu.Unlock()

	response := struct {
		Status string `json:"status"`

		Errors  []model.FileError `json:"errors,omitempty"`
		Archive string            `json:"archive,omitempty"`
	}{
		Status: string(task.Status),
		Errors: task.Errors,
	}
	if task.Status == model.CompletedTk {
		response.Archive = base64.StdEncoding.EncodeToString(task.Archive)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *HandlerTk) processTk(task *model.Task) {

	h.ActiveSlot <- struct{}{}
	defer func() { <-h.ActiveSlot }()

	task.Mu.Lock()
	task.Status = model.RunningTk

	task.Mu.Unlock()

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	client := &http.Client{Timeout: 30 * time.Second}

	for _, fileURL := range task.URLs {
		resp, err := client.Get(fileURL)

		if err != nil || resp.StatusCode != http.StatusOK {
			task.Errors = append(task.Errors, model.FileError{URL: fileURL, Message: "download failed"})
			continue
		}

		defer resp.Body.Close()

		ext := filepath.Ext(fileURL)
		mimeType := resp.Header.Get("Content-Type") //

		if !isAllowed(ext, mimeType, h.Cfg.AlExtens) {
			task.Errors = append(task.Errors, model.FileError{URL: fileURL, Message: "invalid file type"})
			continue
		}

		f, err := zipWriter.Create(filepath.Base(fileURL))

		if err != nil {
			task.Errors = append(task.Errors, model.FileError{URL: fileURL, Message: "zip creation failed"})
			continue
		}
		if _, err := io.Copy(f, resp.Body); err != nil {
			task.Errors = append(task.Errors, model.FileError{URL: fileURL, Message: "copy to zip failed"})
			continue
		}
	}

	if err := zipWriter.Close(); err != nil {

		task.Mu.Lock()
		task.Status = model.Failed

		task.Mu.Unlock()

		return
	}

	task.Mu.Lock()

	task.Archive = buf.Bytes()
	task.Status = model.CompletedTk

	task.Mu.Unlock()
}

func isAllowed(ext, mimeType string, allowed []string) bool {
	for _, a := range allowed {

		if ext == a || mime.TypeByExtension(a) == mimeType {
			return true
		}

	}

	return false
}
