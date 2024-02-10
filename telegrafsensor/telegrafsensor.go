package telegrafsensor

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"reflect"
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

	// telegraf must be configure to output in json format
	cmd := exec.Command("telegraf", "--config", "/opt/homebrew/etc/telegraf.conf", "--once")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		ts.logger.Errorw("Error executing Telegraf", "error", err)
		return nil, err
	}

	for _, mline := range strings.Split(out.String(), "\n") {
		if mline == "" {
			continue
		}

		var metric Metric
		err := json.Unmarshal([]byte(mline), &metric)
		if err != nil {
			ts.logger.Errorw("Error parsing reading", "error", err, "input", mline)
		}

		if metric.Name == "disk" {
			if _, ok := result["disk"]; !ok {
				result["disk"] = map[string]interface{}{}
			}
			path := metric.Tags["path"].(string)
			result["disk"].(map[string]interface{})[path] = toMap(metric)
			continue
		}

		if _, ok := result[metric.Name]; ok {
			mergeMetrics(result[metric.Name].(map[string]interface{}), metric, ts.logger)
		} else {
			result[metric.Name] = toMap(metric)
		}
	}

	ts.logger.Debugw("Readingds result", "result", result)

	return result, nil
}

func toMap(m Metric) map[string]interface{} {
	metric := map[string]interface{}{}

	metric["fields"] = m.Fields
	metric["tags"] = m.Tags
	metric["timestamp"] = m.Timestamp

	return metric
}

func mergeMetrics(old map[string]interface{}, new Metric, logger logging.Logger) map[string]interface{} {
	values := reflect.ValueOf(new)
	types := reflect.TypeOf(new)
	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i).Interface()
		name := strings.ToLower(types.Field(i).Name)
		if name == "name" || reflect.DeepEqual(old[name], value) {
			continue
		}

		switch reflect.TypeOf(value).String() {
		case "map[string]interface {}":
			for k, v := range value.(map[string]interface{}) {
				if _, ok := old[name].(map[string]interface{})[k]; ok {
					logger.Debugw("Duplicate internal key", "name", name, "key", k, "v", v, "old", old, "new", new, "value", value)
					continue
				}

				old[name].(map[string]interface{})[k] = v
			}
		default:
			logger.Debugw("Duplicate but not a map", "name", name, "value", value, "type", reflect.TypeOf(value).String())
			continue
		}
	}

	return old
}
