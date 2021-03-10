package x_http_client

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	LogOff = iota
	Error
	Warn
	Info
	Debug
)

// LogTag Tag for each level of log
var LogTag = []string{"[error]", "[warn]", "[info]", "[debug]"}

// HTTPTimeout defines HTTP timeout.
type HTTPTimeout struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	HeaderTimeout    time.Duration
	LongTimeout      time.Duration
	IdleConnTimeout  time.Duration
}

// HTTPMaxConns defines max idle connections and max idle connections per host
type HTTPMaxConns struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

// CredentialInf is interface for get AccessKeyID,AccessKeySecret,SecurityToken
type Credentials interface {
	GetAccessAppID() string
	GetAccessAppSecret() string
	GetSecurityToken() string
}

// CredentialInfBuild is interface for get CredentialInf
type CredentialsProvider interface {
	GetCredentials() Credentials
}

type defaultCredentials struct {
	config *Config
}

func (defCre *defaultCredentials) GetAccessAppID() string {
	return defCre.config.AccessAppID
}

func (defCre *defaultCredentials) GetAccessAppSecret() string {
	return defCre.config.AccessAppSecret
}
func (defCre *defaultCredentials) GetSecurityToken() string {
	return defCre.config.SecurityToken
}

type defaultCredentialsProvider struct {
	config *Config
}

func (defBuild *defaultCredentialsProvider) GetCredentials() Credentials {
	return &defaultCredentials{config: defBuild.config}
}

type Config struct {
	RetryTimes uint
	UserAgent  string
	Timeout    uint

	Endpoint        string // OSS endpoint
	AccessAppID     string // AccessId
	AccessAppSecret string // AccessKey
	SecurityToken   string // AccessKey

	AuthVersion AuthVersionType // AccessKey

	HTTPTimeout  HTTPTimeout  // HTTP timeout
	HTTPMaxConns HTTPMaxConns // Http max connections

	LocalAddr net.Addr

	CredentialsProvider CredentialsProvider

	AdditionalHeaders []string
	RedirectEnabled   bool

	MD5Threshold        int64
	IsEnableMD5         bool

	IsUseProxy          bool                // Flag of using proxy.
	ProxyHost           string              // Flag of using proxy host.
	IsAuthProxy         bool                // Flag of needing authentication.
	ProxyUser           string              // Proxy user
	ProxyPassword       string              // Proxy password


	LogLevel            int                 // Log level
	Logger              *log.Logger         // For write log
}

// WriteLog output log function
func (config *Config) WriteLog(LogLevel int, format string, a ...interface{}) {
	if config.LogLevel < LogLevel || config.Logger == nil {
		return
	}

	var logBuffer bytes.Buffer
	logBuffer.WriteString(LogTag[LogLevel-1])
	logBuffer.WriteString(fmt.Sprintf(format, a...))
	config.Logger.Printf("%s", logBuffer.String())
}

// for get Credentials
func (config *Config) GetCredentials() Credentials {
	return config.CredentialsProvider.GetCredentials()
}

func getDefaultConfig() *Config {
	config := &Config{}

	config.Endpoint = ""
	config.AccessAppID = ""
	config.AccessAppSecret = ""
	config.RetryTimes = 5
	config.UserAgent = userAgent()
	config.Timeout = 60 // Seconds
	config.SecurityToken = ""

	config.HTTPTimeout.ConnectTimeout = time.Second * 30   // 30s
	config.HTTPTimeout.ReadWriteTimeout = time.Second * 60 // 60s
	config.HTTPTimeout.HeaderTimeout = time.Second * 60    // 60s
	config.HTTPTimeout.LongTimeout = time.Second * 300     // 300s
	config.HTTPTimeout.IdleConnTimeout = time.Second * 50  // 50s
	config.HTTPMaxConns.MaxIdleConns = 100
	config.HTTPMaxConns.MaxIdleConnsPerHost = 100

	config.IsUseProxy = false
	config.ProxyHost = ""
	config.IsAuthProxy = false
	config.ProxyUser = ""
	config.ProxyPassword = ""

	config.MD5Threshold = 16 * 1024 * 1024 // 16MB
	config.IsEnableMD5 = false

	provider := &defaultCredentialsProvider{config: config}
	config.CredentialsProvider = provider
	config.AuthVersion = AuthV1
	return config
}
