package config

type Config struct {
	Celeng  *Celeng
	Logger *Logger
	API    *API
}

type Celeng struct {
	Wallet              string
	Testnet             bool
	PoolURL             string
	Threads             int
	NonInteractive      bool
	DNS                 string
	IgnoreTLSValidation bool
}

type Logger struct {
	Debug     bool
	CLogLevel int8
}

type API struct {
	Transport string
	Listen    string
	Enabled   bool
}

// NewEmpty returns a new empty config
func NewEmpty() *Config {
	return &Config{
		Celeng:  &Celeng{},
		Logger: &Logger{},
		API:    &API{},
	}
}
