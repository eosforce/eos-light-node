package eosio

import (
	"github.com/eosforce/eos-light-node/config"
)

// Genesis eosforce genesis file struct
type Genesis struct {
	InitialTimestamp     string `json:"initial_timestamp"`
	InitialKey           string `json:"initial_key"`
	InitialConfiguration struct {
		MaxBlockNetUsage               int `json:"max_block_net_usage"`
		TargetBlockNetUsagePct         int `json:"target_block_net_usage_pct"`
		MaxTransactionNetUsage         int `json:"max_transaction_net_usage"`
		BasePerTransactionNetUsage     int `json:"base_per_transaction_net_usage"`
		NetUsageLeeway                 int `json:"net_usage_leeway"`
		ContextFreeDiscountNetUsageNum int `json:"context_free_discount_net_usage_num"`
		ContextFreeDiscountNetUsageDen int `json:"context_free_discount_net_usage_den"`
		MaxBlockCPUUsage               int `json:"max_block_cpu_usage"`
		TargetBlockCPUUsagePct         int `json:"target_block_cpu_usage_pct"`
		MaxTransactionCPUUsage         int `json:"max_transaction_cpu_usage"`
		MinTransactionCPUUsage         int `json:"min_transaction_cpu_usage"`
		MaxTransactionLifetime         int `json:"max_transaction_lifetime"`
		DeferredTrxExpirationWindow    int `json:"deferred_trx_expiration_window"`
		MaxTransactionDelay            int `json:"max_transaction_delay"`
		MaxInlineActionSize            int `json:"max_inline_action_size"`
		MaxInlineActionDepth           int `json:"max_inline_action_depth"`
		MaxAuthorityDepth              int `json:"max_authority_depth"`
	} `json:"initial_configuration"`
}

// NewGenesisFromFile new genesis from file
func NewGenesisFromFile(path string) (*Genesis, error) {
	res := &Genesis{}
	err := config.LoadJSONFile(path, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
