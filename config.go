package gateway

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Configuration details how the gateway actually runs
type Configuration struct {
	Monitoring struct {
		HTTP struct {
			Enabled bool   `yaml:"enabled"`
			Path    string `yaml:"path"`
			Address string `yaml:"address"`
			Port    int    `yaml:"port"`
			TLS     struct {
				CertificatePath string `yaml:"cert"`
				KeyPath         string `yaml:"key"`
			} `yaml:"tls"`
		} `yaml:"http"`
		Logging struct {
			ParseRealIP bool `yaml:"real_ip"`
			Outputs     struct {
				Stdout struct {
					Format string `yaml:"format"`
				} `yaml:"std"`
				Systemd struct {
					Format string `yaml:"format"`
				} `yaml:"systemd"`
				Syslog struct {
					Address        string `yaml:"address"`
					Facility       string `yaml:"facility"`
					Format         string `yaml:"format"`
					RFC            string `yaml:"rfc"`
					StructuredData []struct {
						SDID   string            `yaml:"sdId"`
						Values map[string]string `yaml:"values"`
					} `yaml:"data"`
				} `yaml:"syslog"`
			} `yaml:"outputs"`
		} `yaml:"logging"`
		Reporting struct {
			Prefix   string `yaml:"prefix"`
			Includes struct {
				Site   bool `yaml:"site_name"`
				Host   bool `yaml:"host_name"`
				Path   bool `yaml:"path"`
				Method bool `yaml:"method"`
			} `yaml:"includes"`
			TagSets map[string]string `yaml:"tag_sets"`
			Tags    []string          `yaml:"tags"`
			StatsD  struct {
				Address string `yaml:"address"`
			} `yaml:"statsd"`
			DataDog struct {
				Address string `yaml:"address"`
			} `yaml:"dogstatsd"`
			Telegraf struct {
				Address string `yaml:"address"`
			} `yaml:"telegraf"`
			StatSite struct {
				Address string `yaml:"address"`
			} `yaml:"statsite"`
		} `yaml:"reporting"`
	} `yaml:"monitoring"`
	DNS struct {
		Server string `yaml:"server"`
		Port   int    `yaml:"port"`
	} `yaml:"dns"`
	Site struct {
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
			Caching bool `yaml:"caching"`
		} `yaml:"content"`
		Listeners []struct {
			Address         string `yaml:"address"`
			Port            int    `yaml:"port"`
			Force           bool   `yaml:"force"`
			StrictTransport struct {
				Age               time.Duration `yaml:"age"`
				IncludeSubdomains bool          `yaml:"sub_domains"`
				Preload           bool          `yaml:"preload"`
			} `yaml:"htst"`
			TLS struct {
				CertificatePath string `yaml:"cert"`
				KeyPath         string `yaml:"key"`
			}
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
			Address        string        `yaml:"upstream"`
			OpenAPI        string        `yaml:"spec"`
			ConnectTimeout time.Duration `yaml:"timeout_connect"`
			ReadTimeout    time.Duration `yaml:"timeout_read"`
		} `yaml:"services"`
		Discovery struct {
			Mode   string `yaml:"mode"`
			Consul struct {
				Address string `yaml:"address"`
				Token   string `yaml:"token"`
			} `yaml:"consul"`
			Docker struct {
				Address         string `yaml:"address"`
				CertificatePath string `yaml:"cert"`
				KeyPath         string `yaml:"key"`
				CAPath          string `yaml:"ca"`
			}
		} `yaml:"discovery"`
	} `yaml:"site"`
}

func (c *Configuration) String() string {
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

// DefaultConfiguration creates a default set of configuration or the gateway
func DefaultConfiguration() *Configuration {
	config := &Configuration{}

	config.Monitoring.HTTP.Enabled = true
	config.Monitoring.HTTP.Path = "/var/www/gateway/admin"
	config.Monitoring.HTTP.Address = "0.0.0.0"
	config.Monitoring.HTTP.Port = 8080

	config.Monitoring.Logging.ParseRealIP = true
	config.Monitoring.Logging.Outputs.Stdout.Format = "common"
	config.Monitoring.Logging.Outputs.Systemd.Format = "combined"
	config.Monitoring.Logging.Outputs.Syslog.Format = "logfmt"
	config.Monitoring.Logging.Outputs.Syslog.Facility = "local7"
	config.Monitoring.Logging.Outputs.Syslog.RFC = "rfc5424"

	config.Monitoring.Reporting.Prefix = "gateway"
	config.Monitoring.Reporting.Includes.Host = true
	config.Monitoring.Reporting.Includes.Site = true
	config.Monitoring.Reporting.Includes.Path = true
	config.Monitoring.Reporting.Includes.Method = true

	return config
}
