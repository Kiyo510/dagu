package jsondb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/daguflow/dagu/internal/persistence/model"

	"github.com/daguflow/dagu/internal/dag"
	"github.com/daguflow/dagu/internal/dag/scheduler"
	"github.com/stretchr/testify/require"
)

func TestWriteStatusToFile(t *testing.T) {
	tmpDir, db := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := &dag.DAG{
		Name:     "test_write_status",
		Location: "test_write_status.yaml",
	}
	dw, file, err := db.newWriter(d.Location, time.Now(), "request-id-1")
	require.NoError(t, err)
	require.NoError(t, dw.open())
	defer func() {
		_ = dw.close()
		_ = db.RemoveOld(d.Location, 0)
	}()

	status := model.NewStatus(d, nil, scheduler.StatusRunning, 10000, nil, nil)
	status.RequestID = fmt.Sprintf("request-id-%d", time.Now().Unix())
	require.NoError(t, dw.write(status))
	require.Regexp(t, ".*test_write_status.*", file)

	dat, err := os.ReadFile(file)
	require.NoError(t, err)

	r, err := model.StatusFromJSON(string(dat))
	require.NoError(t, err)

	require.Equal(t, d.Name, r.Name)

	err = dw.close()
	require.NoError(t, err)

	// TODO: fixme
	oldS := d.Location
	newS := "text_write_status_new.yaml"

	oldDir := db.getDirectory(oldS, prefix(oldS))
	newDir := db.getDirectory(newS, prefix(newS))
	require.DirExists(t, oldDir)
	require.NoDirExists(t, newDir)

	err = db.Rename(oldS, newS)
	require.NoError(t, err)
	require.NoDirExists(t, oldDir)
	require.DirExists(t, newDir)

	ret := db.ReadStatusRecent(newS, 1)
	require.Equal(t, 1, len(ret))
	require.Equal(t, status.RequestID, ret[0].Status.RequestID)
}

func TestWriteStatusToExistingFile(t *testing.T) {
	tmpDir, db := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := &dag.DAG{Name: "test_append_to_existing", Location: "test_append_to_existing.yaml"}
	dw, file, err := db.newWriter(d.Location, time.Now(), "request-id-1")
	require.NoError(t, err)
	require.NoError(t, dw.open())

	status := model.NewStatus(d, nil, scheduler.StatusCancel, 10000, nil, nil)
	status.RequestID = "request-id-test-write-status-to-existing-file"
	require.NoError(t, dw.write(status))
	_ = dw.close()

	data, err := db.FindByRequestID(d.Location, status.RequestID)
	require.NoError(t, err)
	require.Equal(t, data.Status.Status, scheduler.StatusCancel)
	require.Equal(t, file, data.File)

	dw = &writer{target: file}
	require.NoError(t, dw.open())
	status.Status = scheduler.StatusSuccess
	require.NoError(t, dw.write(status))
	_ = dw.close()

	data, err = db.FindByRequestID(d.Location, status.RequestID)
	require.NoError(t, err)
	require.Equal(t, data.Status.Status, scheduler.StatusSuccess)
	require.Equal(t, file, data.File)
}
