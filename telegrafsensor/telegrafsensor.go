package telegrafsensor

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strings"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	Model = resource.NewModel("aleparedes", "viam-sensor", "telegrafsensor")
)

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newSensor,
		})
}

func newSensor(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (sensor.Sensor, error) {
	// newConfig, err := resource.NativeConfig[*Config](conf)
	// if err != nil {
	// 	return nil, err
	// }
	ts := TelegrafSensor{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
		// Attributes necessary for this sensor component config
	}

	return &ts, nil
}

type TelegrafSensor struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

type Config struct {
}

type Metric struct {
	Name      string                 `json:"name"`
	Fields    map[string]interface{} `json:"fields"`
	Tags      map[string]interface{} `json:"tags"`
	Timestamp uint64                 `json:"timestamp"`
}

func (cfg *Config) Validate(path string) ([]string, error) {
	return nil, nil
}

func (ts *TelegrafSensor) Readings(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	var readings []Metric

	// telegraf must be configure to output in json format
	cmd := exec.Command("telegraf", "--config", "/opt/homebrew/etc/telegraf.conf", "--once")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		ts.logger.Errorw("Error executing Telegraf", "error", err)
		return nil, err
	}

	output := out.String()
	ts.logger.Debugw("Telegraf Output", "json", output)

	for _, mline := range strings.Split(output, "\n") {
		if mline == "" {
			continue
		}

		var metric Metric
		err := json.Unmarshal([]byte(mline), &metric)
		if err != nil {
			ts.logger.Errorw("Error parsing reading", "error", err, "input", mline)
		}

		readings = append(readings, metric)
	}

	result["metrics"] = readings
	ts.logger.Debugw("Readingds result", "result", result)

	return result, nil
}
