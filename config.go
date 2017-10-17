package gateway

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"time"

	"gopkg.in/yaml.v2"
)

// TLSConfiguration defines a cert/key pair
type TLSConfiguration struct {
	CertificatePath string `yaml:"cert"`
	KeyPath         string `yaml:"key"`
}

// Address contains a valid URI scheme based address
type Address string

// URL will return the parsed URL from the address, or error if not valid
func (a Address) URL() (*url.URL, error) {
	return url.Parse(string(a))
}

// HTTPMonitorConfiguration defines the management website configuration
type HTTPMonitorConfiguration struct {
	Enabled bool              `yaml:"enabled"`
	Path    string            `yaml:"path"`
	Address Address           `yaml:"address"`
	TLS     *TLSConfiguration `yaml:"tls"`
}

// LoggingConfiguration defines the logging interface and outputs
type LoggingConfiguration struct {
	ParseRealIP bool `yaml:"real_ip"`
	Outputs     struct {
		Stdout struct {
			Format string `yaml:"format"`
		} `yaml:"std"`
		Systemd struct {
			Format string `yaml:"format"`
		} `yaml:"systemd"`
		Syslog struct {
			Address        Address `yaml:"address"`
			Facility       string  `yaml:"facility"`
			Format         string  `yaml:"format"`
			RFC            string  `yaml:"rfc"`
			StructuredData []struct {
				SDID   string            `yaml:"sdId"`
				Values map[string]string `yaml:"values"`
			} `yaml:"data"`
		} `yaml:"syslog"`
	} `yaml:"outputs"`
}

// MetricsConfiguration defines the metrics reporting formatting and outputs
type MetricsConfiguration struct {
	Prefix   string `yaml:"prefix"`
	Includes struct {
		Site   bool `yaml:"site_name"`
		Host   bool `yaml:"host_name"`
		Path   bool `yaml:"path"`
		Method bool `yaml:"method"`
	} `yaml:"includes"`
	TagSets map[string]string `yaml:"tag_sets"`
	Tags    []string          `yaml:"tags"`
	Address Address           `yaml:"address"`
	Format  string            `yaml:"format"`
}

// SiteConfiguration defines the hosted site
type SiteConfiguration struct {
	Headers struct {
		IncludeServer    bool              `yaml:"server"`
		Blacklist        []string          `yaml:"strip"`
		IncludeRequestID bool              `yaml:"request_id"`
		ExtraHeaders     map[string]string `yaml:"append"`
		IncludeDebug     bool              `yaml:"debug"`
	} `yaml:"headers"`
	Hosts   []string `yaml:"hosts"`
	Content struct {
		Path                            string            `yaml:"path"`
		DefaultDocument                 string            `yaml:"default"`
		EnableSinglePageApplicationMode bool              `yaml:"spa_mode"`
		Errors                          map[string]string `yaml:"errors"`
		Push                            struct {
			Enable   bool `yaml:"enable"`
			Tracking struct {
				Enabled    bool   `yaml:"enabled"`
				CookieName string `yaml:"name"`
			} `yaml:"cookie_tracker"`
		} `yaml:"push"`
		EnableCaching bool `yaml:"caching"`
	} `yaml:"content"`
	Listeners []struct {
		Address         Address `yaml:"address"`
		Force           bool    `yaml:"force"`
		StrictTransport struct {
			Age               time.Duration `yaml:"age"`
			IncludeSubdomains bool          `yaml:"sub_domains"`
			Preload           bool          `yaml:"preload"`
		} `yaml:"htst"`
		TLS *TLSConfiguration `yaml:"tls"`
	} `yaml:"listeners"`
	OpenAPI struct {
		Path   string `yaml:"path"`
		UIPath string `yaml:"ui"`
	} `yaml:"spec"`
	CORS struct {
		DisableAll          bool     `yaml:"disable"`
		Hijack              bool     `yaml:"handle"`
		Origins             []string `yaml:"origins"`
		Methods             []string `yaml:"methods"`
		RequestHeaders      []string `yaml:"request_headers"`
		ResponseHeaders     []string `yaml:"response_headers"`
		AllowAuthentication bool     `yaml:"authentication"`
	} `yaml:"cors"`
	Retry struct {
		Count   int           `yaml:"count"`
		Delay   time.Duration `yaml:"delay"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"retry"`
	Services []struct {
		Path           string        `yaml:"path"`
		Address        Address       `yaml:"address"`
		OpenAPI        string        `yaml:"spec"`
		ConnectTimeout time.Duration `yaml:"timeout_connect"`
		ReadTimeout    time.Duration `yaml:"timeout_read"`
	} `yaml:"services"`
	Discovery struct {
		Mode   string `yaml:"mode"`
		Consul struct {
			Address Address `yaml:"address"`
			Token   string  `yaml:"token"`
		} `yaml:"consul"`
		Docker struct {
			Address         Address `yaml:"address"`
			CertificatePath string  `yaml:"cert"`
			KeyPath         string  `yaml:"key"`
			CAPath          string  `yaml:"ca"`
		}
	} `yaml:"discovery"`
}

// Configuration details how the gateway actually runs
type Configuration struct {
	Monitoring struct {
		HTTP    HTTPMonitorConfiguration `yaml:"http"`
		Logging LoggingConfiguration     `yaml:"logging"`
		Metrics MetricsConfiguration     `yaml:"metrics"`
	} `yaml:"monitoring"`
	DNS struct {
		Address Address `yaml:"address"`
	} `yaml:"dns"`
	Site SiteConfiguration `yaml:"site"`
}

// YAML outputs the configuration in YAML format
func (c *Configuration) YAML() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to serialize configuration: %v", err).Error()
	}

	return string(data)
}

// LoadConfiguration Loads the configuration from the given reader
func LoadConfiguration(r io.Reader) (*Configuration, error) {
	config := DefaultConfiguration()

	readerContents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration from reader: %v", err)
	}

	if err := yaml.Unmarshal(readerContents, config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %v", err)
	}

	return config, nil
}

// DefaultConfiguration creates a default set of configuration for the gateway
//
// This configuration should work 100% out of the box without tweaking
func DefaultConfiguration() *Configuration {
	config := &Configuration{}

	config.Monitoring.HTTP.Enabled = true
	config.Monitoring.HTTP.Path = "./public/monitoring"
	config.Monitoring.HTTP.Address = "tcp://127.0.0.1:8080"

	config.Monitoring.Logging.ParseRealIP = true
	config.Monitoring.Logging.Outputs.Stdout.Format = "combined"
	config.Monitoring.Logging.Outputs.Systemd.Format = "combined"
	config.Monitoring.Logging.Outputs.Syslog.Format = "combined"
	config.Monitoring.Logging.Outputs.Syslog.Facility = "local7"
	config.Monitoring.Logging.Outputs.Syslog.RFC = "rfc5424"

	config.Monitoring.Metrics.Prefix = "gateway"
	config.Monitoring.Metrics.Includes.Host = true
	config.Monitoring.Metrics.Includes.Site = true
	config.Monitoring.Metrics.Includes.Path = true
	config.Monitoring.Metrics.Includes.Method = true

	config.Site.Headers.IncludeServer = true
	config.Site.Headers.IncludeRequestID = true
	config.Site.Headers.IncludeDebug = false

	config.Site.Content.Path = "./public/www"
	config.Site.Content.DefaultDocument = "index.html"
	config.Site.Content.EnableSinglePageApplicationMode = false
	config.Site.Content.EnableCaching = true

	config.Site.CORS.Hijack = true
	config.Site.CORS.DisableAll = true

	config.Site.Retry.Count = 5
	config.Site.Retry.Delay = time.Millisecond * 10
	config.Site.Retry.Timeout = time.Minute * 1

	config.Site.Discovery.Mode = "consul"
	config.Site.Discovery.Consul.Address = "tcp://localhost:8500"

	return config
}
