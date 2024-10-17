package objectives

import (
	"context"
	"flag"
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/service"
	"logbook/internal/logger"
	"logbook/internal/utils/reflux"
	"logbook/models"
	"logbook/models/columns"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type bookmarkFlags struct {
	UserId     string
	SubjectOid string
	SubjectVid string
	Name       string
}

func addBookmark(l *logger.Logger) error {
	var flags bookmarkFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.StringVar(&flags.SubjectOid, "oid", "", "subject oid")
	flag.StringVar(&flags.SubjectVid, "vid", "", "subject vid")
	flag.StringVar(&flags.Name, "name", "My bookmark", "bookmark name")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("addBookmark"))

	err = a.AddBookmark(context.Background(), app.AddBookmarkParams{
		Actor:        columns.UserId(flags.UserId),
		Subject:      models.Ovid{columns.ObjectiveId(flags.SubjectOid), columns.VersionId(flags.SubjectVid)},
		BookmarkName: flags.Name,
	})
	if err != nil {
		return fmt.Errorf("a.Name: %w", err)
	}
	return nil
}

type checkoutFlags struct {
	UserId     string
	SubjectOid string
	SubjectVid string
	To         string
}

func checkout(l *logger.Logger) error {
	var flags checkoutFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.StringVar(&flags.SubjectOid, "oid", "", "subject oid")
	flag.StringVar(&flags.SubjectVid, "vid", "", "subject vid")
	flag.StringVar(&flags.To, "to", "", "target vid")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("checkout"))

	err = a.Checkout(context.Background(), app.CheckoutParams{
		User:    columns.UserId(flags.UserId),
		Subject: models.Ovid{columns.ObjectiveId(flags.SubjectOid), columns.VersionId(flags.SubjectVid)},
		To:      columns.VersionId(flags.To),
	})
	if err != nil {
		return fmt.Errorf("a.Checkout: %w", err)
	}
	return nil
}

type createSubtaskFlags struct {
	UserId    string
	ParentOid string
	ParentVid string
	Content   string
}

func createSubtask(l *logger.Logger) error {
	var flags createSubtaskFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.StringVar(&flags.ParentOid, "oid", "", "parent oid")
	flag.StringVar(&flags.ParentVid, "vid", "", "parent vid")
	flag.StringVar(&flags.Content, "content", "", "target vid")

	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("createSubtask"))

	obj, err := a.CreateSubtask(context.Background(), app.CreateSubtaskParams{
		Creator: columns.UserId(flags.UserId),
		Parent:  models.Ovid{columns.ObjectiveId(flags.ParentOid), columns.VersionId(flags.ParentVid)},
		Content: flags.Content,
	})
	if err != nil {
		return fmt.Errorf("a.CreateSubtask: %w", err)
	}
	fmt.Println(obj)
	return nil
}

type deleteSubtaskFlags struct {
	UserId string
	Oid    string
	Vid    string
}

func deleteSubtask(l *logger.Logger) error {
	var flags deleteSubtaskFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.StringVar(&flags.Oid, "oid", "", "subject oid")
	flag.StringVar(&flags.Vid, "vid", "", "subject vid")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("deleteSubtask"))

	err = a.DeleteSubtask(context.Background(), app.DeleteSubtaskParams{
		Subject: models.Ovid{columns.ObjectiveId(flags.Oid), columns.VersionId(flags.Vid)},
		Actor:   columns.UserId(flags.UserId),
	})
	if err != nil {
		return fmt.Errorf("a.DeleteSubtask: %w", err)
	}
	return nil
}

type getActiveVersionFlags struct {
	Oid string
}

func getActiveVersion(l *logger.Logger) error {
	var flags getActiveVersionFlags
	flag.StringVar(&flags.Oid, "oid", "", "")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("getActiveVersion"))

	vid, err := a.GetActiveVersion(context.Background(), columns.ObjectiveId(flags.Oid))
	if err != nil {
		return fmt.Errorf("a.GetActiveVersion: %w", err)
	}
	fmt.Println(vid)
	return nil
}

type getMergedPropsFlags struct {
	Oid string
	Vid string
}

func getMergedProps(l *logger.Logger) error {
	var flags getMergedPropsFlags
	flag.StringVar(&flags.Oid, "oid", "", "parent oid")
	flag.StringVar(&flags.Vid, "vid", "", "parent vid")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("getMergedProps"))

	mprops, err := a.GetMergedProps(context.Background(), models.Ovid{columns.ObjectiveId(flags.Oid), columns.VersionId(flags.Vid)})
	if err != nil {
		return fmt.Errorf("a.GetMergedProps: %w", err)
	}
	fmt.Println(mprops)
	return nil
}

type getObjectiveHistoryFlags struct {
	Oid string
	Vid string
}

func getObjectiveHistory(l *logger.Logger) error {
	var flags getObjectiveHistoryFlags
	flag.StringVar(&flags.Oid, "oid", "", "oid")
	flag.StringVar(&flags.Vid, "vid", "", "vid")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("getObjectiveHistory"))

	hist, err := a.GetObjectiveHistory(context.Background(), app.GetObjectiveHistoryParams{
		Subject: models.Ovid{Oid: columns.ObjectiveId(flags.Oid), Vid: columns.VersionId(flags.Vid)},
	})
	if err != nil {
		return fmt.Errorf("a.GetObjectiveHistory: %w", err)
	}
	for _, e := range hist {
		fmt.Println(e)
	}
	return nil
}

type listBookmarksFlags struct {
	UserId string
}

func listBookmarks(l *logger.Logger) error {
	var flags listBookmarksFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("listBookmarks"))

	bms, err := a.ListBookmarks(context.Background(), app.ListBookmarksParams{
		Viewer: columns.UserId(flags.UserId),
	})
	if err != nil {
		return fmt.Errorf("a.ListBookmarks: %w", err)
	}
	for _, e := range bms {
		fmt.Println(
			"Bid:", e.Bid,
			"Oid:", e.Oid,
			"IsRock:", e.IsRock,
			"Title:", e.Title,
			"CreatedAt:", e.CreatedAt,
		)
	}
	return nil
}

// type reattachFlags struct {
// }

// func reattach(l *logger.Logger) error {
// 	var flags reattachFlags
// 	flag.Parse()

// 	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
// 		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
// 	}

// 	srvcfg, err := service.ReadConfig(configpath)
// 	if err != nil {
// 		return fmt.Errorf("service.ReadConfig: %w", err)
// 	}
// 	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
// 	if err != nil {
// 		return fmt.Errorf("pgxpool.New: %w", err)
// 	}
// 	defer pool.Close()
// 	a := app.New(pool)

// 	err = a.Reattach(context.Background(), app.ReattachParams{})
// 	if err != nil {
// 		return fmt.Errorf("a.Reattach: %w", err)
// 	}
// 	return nil
// }

// type reorderFlags struct {
// }

// func reorder(l *logger.Logger) error {
// 	var flags reorderFlags
// 	flag.Parse()

// 	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
// 		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
// 	}

// 	srvcfg, err := service.ReadConfig(configpath)
// 	if err != nil {
// 		return fmt.Errorf("service.ReadConfig: %w", err)
// 	}
// 	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
// 	if err != nil {
// 		return fmt.Errorf("pgxpool.New: %w", err)
// 	}
// 	defer pool.Close()
// 	a := app.New(pool)

// 	err = a.Reorder(context.Background(), app.ReorderParams{})
// 	if err != nil {
// 		return fmt.Errorf("a.Reorder: %w", err)
// 	}
// 	return nil
// }

type rockCreateFlags struct {
	UserId string
}

func rockCreate(l *logger.Logger) error {
	var flags rockCreateFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("rockCreate"))

	err = a.RockCreate(context.Background(), columns.UserId(flags.UserId))
	if err != nil {
		return fmt.Errorf("a.RockCreate: %w", err)
	}
	return nil
}

type rockGetFlags struct {
	UserId string
}

func rockGet(l *logger.Logger) error {
	var flags rockGetFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		return fmt.Errorf("found zero values: %s", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("rockGet"))

	oid, err := a.RockGet(context.Background(), columns.UserId(flags.UserId))
	if err != nil {
		return fmt.Errorf("a.RockGet: %w", err)
	}
	fmt.Println(oid)
	return nil
}

type viewBuilderFlags struct {
	UserId        string
	RockOid       string
	RockVid       string
	Start, Length int
}

func viewBuilder(l *logger.Logger) error {
	var flags viewBuilderFlags
	flag.StringVar(&flags.UserId, "uid", "00000000-0000-0000-0000-000000000000", "")
	flag.StringVar(&flags.RockOid, "oid", "", "rock oid")
	flag.StringVar(&flags.RockVid, "vid", "", "rock vid")
	flag.IntVar(&flags.Start, "start", 0, "")
	flag.IntVar(&flags.Length, "length", 1, "")
	flag.Parse()

	if empties := reflux.FindZeroValues(flags); len(empties) > 0 {
		fmt.Printf("found zero values: %s\n", strings.Join(empties, ", "))
	}

	srvcfg, err := service.ReadConfig(configpath)
	if err != nil {
		return fmt.Errorf("service.ReadConfig: %w", err)
	}
	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	defer pool.Close()
	a := app.New(pool, l.Sub("viewBuilder"))

	d, err := a.ViewBuilder(context.Background(), app.ViewBuilderParams{
		Viewer: columns.UserId(flags.UserId),
		Root:   models.Ovid{columns.ObjectiveId(flags.RockOid), columns.VersionId(flags.RockVid)},
		Start:  flags.Start,
		Length: flags.Length,
	})
	if err != nil {
		return fmt.Errorf("ViewBuilder: %w", err)
	}
	for _, e := range d {
		fmt.Println(
			"Depth:", e.Depth,
			"ObjectiveType:", e.ObjectiveType,
			"Folded:", e.Folded,
			"Oid:", e.Oid,
			"Vid:", e.Vid,
		)
	}
	return nil
}
