package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

// ChainConfig holds configuration for a blockchain network
type ChainConfig struct {
	ChainID       int64             `json:"chainId"`
	Name          string            `json:"name"`
	RPCUrls       []string          `json:"rpcUrls"`
	WSUrls        []string          `json:"wsUrls"` // WebSocket URLs
	Contracts     ContractAddresses `json:"contracts"`
	BlockTime     int               `json:"blockTime"`     // seconds
	Confirmations int               `json:"confirmations"` // blocks
	StartBlock    uint64            `json:"startBlock"`    // Block to start indexing from
}

// ContractAddresses holds deployed contract addresses
type ContractAddresses struct {
	CTFExchange       string `json:"ctfExchange"`
	ConditionalTokens string `json:"conditionalTokens"`
}

// Config holds all chain configurations
type Config struct {
	Chains map[string]*ChainConfig `json:"chains"`
}

// LoadConfig loads chain configuration from JSON file
func LoadConfig(filepath string) (*Config, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// GetChain returns configuration for a specific chain
func (c *Config) GetChain(name string) (*ChainConfig, error) {
	chain, ok := c.Chains[name]
	if !ok {
		return nil, fmt.Errorf("chain %s not found in config", name)
	}
	return chain, nil
}

// GetCTFExchangeAddress returns the CTFExchange contract address as common.Address
func (cc *ChainConfig) GetCTFExchangeAddress() common.Address {
	return common.HexToAddress(cc.Contracts.CTFExchange)
}

// GetConditionalTokensAddress returns the ConditionalTokens contract address
func (cc *ChainConfig) GetConditionalTokensAddress() common.Address {
	return common.HexToAddress(cc.Contracts.ConditionalTokens)
}

// GetAllContractAddresses returns all contract addresses as a slice
func (cc *ChainConfig) GetAllContractAddresses() []common.Address {
	return []common.Address{
		cc.GetCTFExchangeAddress(),
		cc.GetConditionalTokensAddress(),
	}
}

// GetAllContractAddressStrings returns all contract addresses as strings
func (cc *ChainConfig) GetAllContractAddressStrings() []string {
	return []string{
		cc.Contracts.CTFExchange,
		cc.Contracts.ConditionalTokens,
	}
}
