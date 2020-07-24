package models

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
	DbType   string `json:"dbType"`
	Port     string `json:"port"`
	Host     string `json:"host"`
}
