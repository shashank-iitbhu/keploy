package generateConfig

import (
	"os"
	"os/exec"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
  "go.keploy.io/server/utils"
)

var Emoji = "\U0001F430" + " Keploy:"

type generatorConfig struct {
	logger *zap.Logger
	mutex  sync.Mutex
}

func NewGeneratorConfig(logger *zap.Logger) GeneratorConfig {
	return &generatorConfig{
		logger: logger,
		mutex:  sync.Mutex{},
	}
}

var config = `
record:
  path: ""
  # mandatory
  command: ""
  proxyport: 0
  containerName: ""
  networkName: ""
  delay: 5
  buildDelay: 30s
  passThroughPorts: []
  filters:
    ReqHeader: []
    urlMethods: {}
test:
  path: ""
  # mandatory
  command: ""
  proxyport: 0
  containerName: ""
  networkName: ""
  # example: "test-set-1": ["test-1", "test-2", "test-3"]
  tests:
  # to use globalNoise, please follow the guide at the end of this file.
  globalNoise:
    global:
      body: {}
      header: {}
  delay: 5
  buildDelay: 30s
  apiTimeout: 5
  passThroughPorts: []
  withCoverage: false
  coverageReportPath: ""
  `

func (g *generatorConfig) GenerateConfig(filePath string) {
	var node yaml.Node
	data := []byte(config)
	if err := yaml.Unmarshal(data, &node); err != nil {
		g.logger.Fatal("Unmarshalling failed %s", zap.Error(err))
	}
	results, err := yaml.Marshal(node.Content[0])
	if err != nil {
		g.logger.Fatal("Failed to marshal the config", zap.Error(err))
	}

  finalOutput := append(results, []byte(utils.ConfigGuide)...)

	err = os.WriteFile(filePath, finalOutput, os.ModePerm)
	if err != nil {
		g.logger.Fatal("Failed to write config file", zap.Error(err))
	}

	cmd := exec.Command("sudo", "chmod", "-R", "777", filePath)
	err = cmd.Run()
	if err != nil {
		g.logger.Error("failed to set the permission of config file", zap.Error(err))
		return
	}

	g.logger.Info("Config file generated successfully")
}
