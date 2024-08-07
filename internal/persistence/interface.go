package persistence

import (
	"fmt"
	"time"

	"github.com/daguflow/dagu/internal/dag"
	"github.com/daguflow/dagu/internal/persistence/grep"
	"github.com/daguflow/dagu/internal/persistence/model"
)

var (
	ErrRequestIDNotFound = fmt.Errorf("request id not found")
	ErrNoStatusDataToday = fmt.Errorf("no status data today")
	ErrNoStatusData      = fmt.Errorf("no status data")
)

type DataStores interface {
	HistoryStore() HistoryStore
	DAGStore() DAGStore
	FlagStore() FlagStore
}

type HistoryStore interface {
	Open(dagFile string, t time.Time, requestID string) error
	Write(status *model.Status) error
	Close() error
	Update(dagFile, requestID string, st *model.Status) error
	ReadStatusRecent(dagFile string, n int) []*model.StatusFile
	ReadStatusToday(dagFile string) (*model.Status, error)
	FindByRequestID(dagFile string, requestID string) (*model.StatusFile, error)
	RemoveAll(dagFile string) error
	RemoveOld(dagFile string, retentionDays int) error
	Rename(oldName, newName string) error
}

type DAGStore interface {
	Create(name string, spec []byte) (string, error)
	Delete(name string) error
	List() (ret []*dag.DAG, errs []string, err error)
	GetMetadata(name string) (*dag.DAG, error)
	GetDetails(name string) (*dag.DAG, error)
	Grep(pattern string) (ret []*GrepResult, errs []string, err error)
	Load(name string) (*dag.DAG, error)
	Rename(oldID, newID string) error
	GetSpec(name string) (string, error)
	UpdateSpec(name string, spec []byte) error
	Find(name string) (*dag.DAG, error)
}

type GrepResult struct {
	Name    string
	DAG     *dag.DAG
	Matches []*grep.Match
}

type FlagStore interface {
	ToggleSuspend(id string, suspend bool) error
	IsSuspended(id string) bool
}
