package models

type Config struct {
	Port  string         `json:"port"`
	DbCfg DatabaseConfig `json:"database"`
}
