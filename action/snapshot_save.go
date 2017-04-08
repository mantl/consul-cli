package action

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type snapshotSave struct {
	*config
}

func SnapshotSaveAction() Action {
	return &snapshotSave{
		config: &gConfig,
	}
}

func (s *snapshotSave) CommandFlags() *flag.FlagSet {
	return s.newFlagSet(FLAG_DATACENTER, FLAG_STALE)
}

func (s *snapshotSave) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Path to snapshot file must be specifie")
	case len(args) > 1:
		return fmt.Errorf("Only one file path allowed")
	}
	filePath := args[0]

	snap, err := s.newSnapshot()
	if err != nil {
		return err
	}

	queryOpts := s.queryOptions()

	reader, _, err := snap.Save(queryOpts)
	if err != nil {
		return err
	}
	defer reader.Close()

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}
