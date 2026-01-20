package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"gopkg.in/yaml.v3"

	"github.com/mirror520/qdm-sync/persistence/mongo"
	"github.com/mirror520/qdm-sync/qdm"

	sync "github.com/mirror520/qdm-sync"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	path := homeDir + "/.qdm-sync"

	app := &cli.App{
		Name:        "QDMSync",
		Description: "QDMSync uses QDM API to sync e-commerce data to MongoDB, streamlining data management.",
		Commands: []*cli.Command{
			{
				Name:        "sync",
				Description: "Initiates data synchronization.",
				Subcommands: []*cli.Command{
					{
						Name: "orders",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "path",
								Usage:   "Specifies the working directory",
								EnvVars: []string{"QDM_PATH"},
								Value:   path,
							},
							&cli.TimestampFlag{
								Name:     "start-time",
								Aliases:  []string{"start", "since"},
								Layout:   time.RFC3339,
								Timezone: time.Local,
								Required: true,
							},
							&cli.TimestampFlag{
								Name:     "end-time",
								Aliases:  []string{"end"},
								Layout:   time.RFC3339,
								Timezone: time.Local,
								Value:    cli.NewTimestamp(time.Now()),
							},
						},
						Action: syncOrders,
					},
					{
						Name: "customers",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "path",
								Usage:   "Specifies the working directory",
								EnvVars: []string{"QDM_PATH"},
								Value:   path,
							},
							&cli.TimestampFlag{
								Name:     "start-time",
								Aliases:  []string{"start", "since"},
								Layout:   time.RFC3339,
								Timezone: time.Local,
								Required: true,
							},
							&cli.TimestampFlag{
								Name:     "end-time",
								Aliases:  []string{"end"},
								Layout:   time.RFC3339,
								Timezone: time.Local,
								Value:    cli.NewTimestamp(time.Now()),
							},
						},
						Action: syncCustomers,
					},
					{
						Name: "customer-groups",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "path",
								Usage:   "Specifies the working directory",
								EnvVars: []string{"QDM_PATH"},
								Value:   path,
							},
						},
						Action: syncCustomerGroups,
					},
				},
			},
		},
		Action: cli.ShowAppHelp,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func syncOrders(cli *cli.Context) error {
	filepath := filepath.Join(cli.String("path"), "config.yaml")

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	var cfg *sync.Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}

	qdm, err := qdm.NewService(cfg.QDM)
	if err != nil {
		return err
	}
	defer qdm.Close()

	repo, err := mongo.NewOrderRepository(cfg.Persistence)
	if err != nil {
		return err
	}
	defer repo.Disconnected()

	svc := sync.NewService(qdm, repo)
	defer svc.Close()

	start := *cli.Timestamp("start-time")
	end := time.Now()
	if endTS := cli.Timestamp("end-time"); endTS != nil {
		end = *endTS
	}

	ch, n, err := svc.SyncOrders(start, end)
	if err != nil {
		return err
	}

	progress := mpb.New()
	defer progress.Shutdown()

	bar := progress.AddBar(n,
		mpb.PrependDecorators(
			decor.Name("synchronizing", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)

	for p := range ch {
		bar.SetCurrent(p.Current)

		if bar.Completed() {
			break
		}
	}

	progress.Wait()

	return nil
}

func syncCustomers(cli *cli.Context) error {
	filepath := filepath.Join(cli.String("path"), "config.yaml")

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	var cfg *sync.Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}

	qdm, err := qdm.NewService(cfg.QDM)
	if err != nil {
		return err
	}
	defer qdm.Close()

	repo, err := mongo.NewOrderRepository(cfg.Persistence)
	if err != nil {
		return err
	}
	defer repo.Disconnected()

	svc := sync.NewService(qdm, repo)
	defer svc.Close()

	start := *cli.Timestamp("start-time")
	end := time.Now()
	if endTS := cli.Timestamp("end-time"); endTS != nil {
		end = *endTS
	}

	ch, n, err := svc.SyncCustomers(start, end)
	if err != nil {
		return err
	}

	progress := mpb.New()
	defer progress.Shutdown()

	bar := progress.AddBar(n,
		mpb.PrependDecorators(
			decor.Name("synchronizing", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)

	for p := range ch {
		bar.SetCurrent(p.Current)

		if bar.Completed() {
			break
		}
	}

	progress.Wait()

	return nil
}

func syncCustomerGroups(cli *cli.Context) error {
	filepath := filepath.Join(cli.String("path"), "config.yaml")

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	var cfg *sync.Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}

	qdm, err := qdm.NewService(cfg.QDM)
	if err != nil {
		return err
	}
	defer qdm.Close()

	repo, err := mongo.NewOrderRepository(cfg.Persistence)
	if err != nil {
		return err
	}
	defer repo.Disconnected()

	svc := sync.NewService(qdm, repo)
	defer svc.Close()

	n, err := svc.SyncCustomerGroups()
	if err != nil {
		return err
	}

	fmt.Println("record stored: " + strconv.Itoa(n))
	return nil
}
