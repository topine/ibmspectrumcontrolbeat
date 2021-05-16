package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/topine/ibmspectrumcontrolbeat/config"
	"github.com/topine/ibmspectrumcontrolbeat/ibmspectrum"
)

// ibmspectrumcontrolbeat configuration.
type ibmspectrumcontrolbeat struct {
	done           chan struct{}
	config         config.Config
	client         beat.Client
	spectrumClient ibmspectrum.Client
}

// New creates an instance of ibmspectrumcontrolbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	if err := c.GetMetricsConf(); err != nil {
		return nil, fmt.Errorf("Error reading config metrics file: %v", err)
	}

	log := logp.NewLogger("ibmspectrumcontrolbeat")
	ibmSpectrumClient := ibmspectrum.NewClient(log, *c.MetricsConfig, c.Username, c.Password, c.BaseURL)
	bt := &ibmspectrumcontrolbeat{
		done:           make(chan struct{}),
		config:         c,
		spectrumClient: *ibmSpectrumClient,
	}
	return bt, nil
}

// Run starts ibmspectrumcontrolbeat.
func (bt *ibmspectrumcontrolbeat) Run(b *beat.Beat) error {
	logp.Info("ibmspectrumcontrolbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		collectedMetrics, err := bt.spectrumClient.CollectFromStorage("svc.*|SVC.*|p.*")
		if err != nil {
			return err
		}

		events, err := bt.mapToBeatEvents(*collectedMetrics, b)
		if err != nil {
			return err
		}

		for _, event := range events {
			bt.client.Publish(event)
			logp.Info("Event sent")
		}
	}
}

// Stop stops ibmspectrumcontrolbeat.
func (bt *ibmspectrumcontrolbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *ibmspectrumcontrolbeat) mapToBeatEvents(collectedMetrics ibmspectrum.CollectedStorageMetrics,
	b *beat.Beat) ([]beat.Event, error) {
	var events []beat.Event

	for _, metric := range collectedMetrics.Metrics {
		fields := common.MapStr{
			"type": b.Info.Name,
		}

		for _, sMetric := range metric.StorageSystemMetrics {
			for x := 1; x <= len(sMetric.Current); x++ {
				// get the latest available metric
				current := sMetric.Current[len(sMetric.Current)-x]
				if current.Y == nil {
					continue
				}
				fields["metricID"] = current.Y
			}
		}

		event := beat.Event{
			Timestamp: time.Now(),
			Fields:    fields,
		}

		events = append(events, event)
	}

	return events, nil
}
