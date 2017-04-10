package action

import (
	"flag"
	"fmt"
	"os"
)

type snapshotRestore struct {
	*config
}

func SnapshotRestoreAction() Action {
	return &snapshotRestore{
		config: &gConfig,
	}
}

func (s *snapshotRestore) CommandFlags() *flag.FlagSet {
	return s.newFlagSet(FLAG_DATACENTER)
}

func (s *snapshotRestore) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Path to snapshot file must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one snapshot file can be specified")
	}
	filePath := args[0]

	snap, err := s.newSnapshot()
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writeOpts := s.writeOptions()

	return snap.Restore(writeOpts, f)
}
