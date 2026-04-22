package config

// DBConfig holds database configuration
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// NewDBConfig returns default database configuration
func NewDBConfig() *DBConfig {
	return &DBConfig{
		User:     "root",
		Password: "1234",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "finance_db",
	}
}
