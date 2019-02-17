package config

type Config struct {
	MySQL   MySQLConfig
	Redis   RedisConfig
	Http    HttpConfig
	Page    []WebSite
	Crawler CrawlerConfig
}

type MySQLConfig struct {
	Username       string
	PasswordEnvKey string
	Host           string
	Port           int
	DatabaseName   string
	Connection     string
}

type RedisConfig struct {
	Connection string
	Host       string
	Port       int
}

type HttpConfig struct {
	ProxyUrl string
}

type WebSite struct {
	Namespace        string
	StartPage     string
	AllowedDomain string
}

type CrawlerConfig struct {
	GenerateRequestLimit int
	HeartBeat            int
	DownloaderNum        int
}
