package objectives

import (
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/cmd/objectives/service"
	"logbook/internal/utilities/slicew/lines"
	"os"
	"path/filepath"

	"golang.org/x/exp/maps"
)

var configpath = filepath.Join(os.Getenv("WORKSPACE"), "backend/cmd/objectives/local.yml")

func migrate() error {
	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	err = database.RunMigration(srvcfg)
	if err != nil {
		return fmt.Errorf("queries.RunMigration: %w", err)
	}
	return nil
}

func Run() error {
	subcmd := os.Args[1]
	os.Args = os.Args[1:]

	handlers := map[string]func() error{
		"addbookmark":         addBookmark,
		"checkout":            checkout,
		"createsubtask":       createSubtask,
		"deletesubtask":       deleteSubtask,
		"getactiveversion":    getActiveVersion,
		"getmergedprops":      getMergedProps,
		"getobjectivehistory": getObjectiveHistory,
		"listbookmarks":       listBookmarks,
		"migrate":             migrate,
		"rockcreate":          rockCreate,
		"rockget":             rockGet,
		"viewbuilder":         viewBuilder,
		// "reattach":            reattach,
		// "reorder":             reorder,
	}
	handler, ok := handlers[subcmd]
	if !ok {
		return fmt.Errorf("handler not found: %s\n\navailable handlers:\n%s", subcmd, lines.Join(maps.Keys(handlers), "  "))
	}
	err := handler()
	if err != nil {
		return fmt.Errorf("%s: %w", subcmd, err)
	}

	return nil
}
