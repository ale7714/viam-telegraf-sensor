package telegrafsensor

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

const (
	telegrafConfPath = "/tmp/viam-telegraf.conf"
	baseTelegrafConf = "conf/base.conf"
)

//go:embed conf
var embeddedConfFS embed.FS

func newTelegrafConf(conf resource.Config, logger logging.Logger) error {
	newConf, err := resource.TransformAttributeMap[*Config](conf.Attributes)
	if err != nil {
		logger.Errorf("failed to configure sensor with %+v", conf)
		return err
	}
	logger.Debugf("configuring sensor with %+v", newConf)

	baseConfigData, err := embeddedConfFS.ReadFile(baseTelegrafConf)
	if err != nil {
		return fmt.Errorf("error reading base.conf: %v", err)
	}

	err = os.WriteFile(telegrafConfPath, baseConfigData, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	v := reflect.ValueOf(*newConf)
	typeOfS := v.Type()
	emptyInputs := true

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i).Interface()
		if skipSection, ok := fieldValue.(bool); ok && !skipSection {
			fieldName := typeOfS.Field(i).Name
			confName := strings.TrimPrefix(strings.ToLower(fieldName), "disable")
			configFileName := fmt.Sprintf("conf/inputs/%s.conf", confName)

			templateData, err := embeddedConfFS.ReadFile(configFileName)
			if err != nil {
				logger.Errorf("Skipping %s. Error reading %v", configFileName, err)
				continue
			}

			destFile, err := os.OpenFile(telegrafConfPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				logger.Errorf("Error opening file %s: %v", telegrafConfPath, err)
				continue
			}

			if _, err := destFile.Write(templateData); err != nil {
				logger.Errorf("Error writing config file %s: %v", telegrafConfPath, err)
			}
			destFile.Close()
			emptyInputs = false
			logger.Debugf("Added config section for %s metric", confName)
		}

	}

	if emptyInputs {
		return errors.New("all Telegraf input metrics disabled. At least one metric must be enabled")
	}

	return nil
}

type Config struct {
	resource.TriviallyValidateConfig
	DisableCpu       bool `json:"disable_cpu"`
	DisableDisk      bool `json:"disable_disk"`
	DisableDiskIo    bool `json:"disable_disk_io"`
	DisableKernel    bool `json:"disable_kernel"`
	DisableMem       bool `json:"disable_mem"`
	DisableNet       bool `json:"disable_net"`
	DisableNetstat   bool `json:"disable_netstat"`
	DisableProcesses bool `json:"disable_processes"`
	DisableSwap      bool `json:"disable_swap"`
	DisableSystem    bool `json:"disable_system"`
	DisableTemp      bool `json:"disable_temp"`
	DisableWireless  bool `json:"disable_wireless"`
}
