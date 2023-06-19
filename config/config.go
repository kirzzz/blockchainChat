package config

import (
	"blockchainStorage/common"
	"encoding/json"
	"os"
)

type Config struct {
	Port              int           `json:"port"`
	DataBasePath      string        `json:"dbPath"`
	Nodes             []common.Node `json:"nodes"`
	Difficulty        int           `json:"difficulty"`
	BlockReward       int           `json:"blockReward"`
	GenesisBlockNonce int64         `json:"genesisBlockNonce"`
}

func LoadConfig(filePath string) (*Config, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
