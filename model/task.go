package model

import (
	"sync"
	"time"
)

type TaskStatus string

const (
	PendingTk   TaskStatus = "pending"
	RunningTk   TaskStatus = "running"
	CompletedTk TaskStatus = "completed"
	Failed      TaskStatus = "failed"
)

type Task struct {
	ID string `json:"id"`

	URLs    []string    `json:"urls"`
	Archive []byte      `json:"archive_data,omitempty"`
	Errors  []FileError `json:"errors,omitempty"`

	CreatedAt time.Time  `json:"-"`
	Status    TaskStatus `json:"status"`
	Mu        sync.Mutex `json:"-"`
}

type FileError struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}
