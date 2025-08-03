package properties

import (
	"os"
)

type Controller struct {
	propertiesFileLocation string
	applicationProperties  ApplicationProperties
}

// NewPropertiesController
// panic justification: this struct is intended to be a singleton that never gets restarted and there is no recoverable
// state if the properties are not loaded
func NewPropertiesController() Controller {
	cfgPath := getCfgFilePath()
	if cfgPath == "" {
		defaultCfg, err := writeDefaultProperties()
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

// SetApplicationProperties todo: remove if this is not needed for wails
func (p *Controller) SetApplicationProperties(applicationProperties ApplicationProperties) {
	p.applicationProperties = applicationProperties
}

func GetDefaultApplicationProperties() (ApplicationProperties, error) {
	var blank ApplicationProperties
	gameDir, err := getDefaultGameDir()
	if err != nil {
		return blank, err
	}
	packData, err := createDefaultPackData()
	if err != nil {
		return blank, err
	}

	return ApplicationProperties{
		GameDir:     gameDir,
		PackDataDir: packData,
	}, nil
}
