package jsondb

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/daguflow/dagu/internal/persistence/model"

	"github.com/daguflow/dagu/internal/util"
)

// Writer is the interface to write status to local file.
type writer struct {
	target  string
	dagFile string
	writer  *bufio.Writer
	file    *os.File
	mu      sync.Mutex
	closed  bool
}

// Open opens the writer.
func (w *writer) open() (err error) {
	_ = os.MkdirAll(filepath.Dir(w.target), 0755)
	w.file, err = util.OpenOrCreateFile(w.target)
	if err == nil {
		w.writer = bufio.NewWriter(w.file)
	}
	return
}

// Writer appends the status to the local file.
func (w *writer) write(st *model.Status) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	jsonb, _ := st.ToJSON()
	str := strings.ReplaceAll(string(jsonb), "\n", " ")
	str = strings.ReplaceAll(str, "\r", " ")
	_, err := w.writer.WriteString(str + "\n")
	util.LogErr("write status", err)
	return w.writer.Flush()
}

// Close closes the writer.
func (w *writer) close() (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.closed {
		err = w.writer.Flush()
		util.LogErr("flush file", err)
		util.LogErr("file sync", w.file.Sync())
		util.LogErr("file close", w.file.Close())
		w.closed = true
	}
	return err
}
