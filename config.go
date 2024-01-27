package sync

import "github.com/mirror520/qdm-sync/qdm"

type Config struct {
	QDM         qdm.Config  `yaml:"qdm"`
	Persistence Persistence `yaml:"persistence"`
}

type Persistence struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
}
