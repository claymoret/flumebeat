package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/claymoret/flumebeat/config"
)

type Flumebeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Flumebeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Flumebeat) Run(b *beat.Beat) error {
	logp.Info("flumebeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		for _, host := range bt.config.Hosts {
			url := host.GetMetricsUrl()
			logp.Info("Processing: %v", url)

			var mm MetricsMap
			if err := mm.Fetch(url); err != nil {
				logp.Err("An error occurred while fetching metrics: %v", err)
				continue
			}

			pm := mm.Parse()

			curr := common.Time(time.Now())
			for _, metric := range pm {
				metric["@timestamp"] = curr
				metric["from"] = host
				bt.client.PublishEvent(common.MapStr(metric))
			}

			logp.Info("%d Events sent for %v", len(pm), url)
		}
	}
}

func (bt *Flumebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
