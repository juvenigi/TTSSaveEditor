package properties

import (
	"os"
)

type Controller struct {
	propertiesFileLocation string
	applicationProperties  ApplicationProperties
}

// InitPropertiesController
// panic justification: this struct is intended to be a singleton that never gets restarted and there is no recoverable
// state if the properties are not loaded
func InitPropertiesController() Controller {
	cfgPath := getCfgFilePath()
	if cfgPath == "" {
		defaultCfg, err := WriteDefaultProperties()
		if err != nil {
			panic(err)
		}
		cfgPath = defaultCfg
	}

	propertiesFile, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer propertiesFile.Close()

	properties, err := newProperties(propertiesFile)
	if err != nil {
		panic(err)
	}

	return Controller{
		propertiesFileLocation: cfgPath,
		applicationProperties:  properties,
	}
}

func (p *Controller) GetApplicationProperties() ApplicationProperties {
	return p.applicationProperties
}
