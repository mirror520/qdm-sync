package main

import (
	"log"
	"os"
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
	app := &cli.App{
		Name:        "QDMSync",
		Description: "QDMSync uses QDM API to sync e-commerce data to MongoDB, streamlining data management.",
		Commands: []*cli.Command{
			{
				Name:        "sync",
				Description: "Initiates data synchronization.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "path",
						Usage:   "Specifies the working directory",
						EnvVars: []string{"QDM_PATH"},
					},
					&cli.TimestampFlag{
						Name:     "start-time",
						Aliases:  []string{"start", "since"},
						Layout:   "2006-01-02T15:04:05",
						Timezone: time.Local,
						Required: true,
					},
					&cli.TimestampFlag{
						Name:     "end-time",
						Aliases:  []string{"end"},
						Layout:   "2006-01-02T15:04:05",
						Timezone: time.Local,
						Value:    cli.NewTimestamp(time.Now()),
					},
				},
				Action: synchronize,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Usage:   "Specifies the working directory",
				EnvVars: []string{"QDM_PATH"},
			},
			&cli.IntFlag{
				Name:    "port",
				Usage:   "Specifies the HTTP service port",
				Value:   8080,
				EnvVars: []string{"QDM_HTTP_PORT"},
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cli *cli.Context) error {
	return nil
}

func synchronize(cli *cli.Context) error {
	path := cli.String("path")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path = homeDir + "/.qdm-sync"
	}

	f, err := os.Open(path + "/config.yaml")
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

	ch, n, err := svc.Sync(start, end)
	if err != nil {
		return err
	}

	progress := mpb.New()
	defer progress.Shutdown()

	bar := progress.AddBar(n,
		mpb.PrependDecorators(
			decor.Name("syncronizing", decor.WCSyncSpaceR),
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
