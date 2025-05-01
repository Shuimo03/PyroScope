package config

// ref: https://github.com/prometheus/prometheus/blob/main/config/config.go
import (
	"errors"
	"fmt"
	"net/url"
	"pyroscope/pkg/logger"
	_ "pyroscope/pkg/logger"
	"strings"
	"time"

	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/model/relabel"

	"github.com/alecthomas/units"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GlobalConfig  GlobalConfig    `yaml:"global"`
	Runtime       RuntimeConfig   `yaml:"runtime,omitempty"`
	Alert         AlertingConfig  `yaml:"alert,omitempty"`
	Rule          RuleConfig      `yaml:"rule,omitempty"`
	Scrape        ScrapeConfig    `yaml:"scrape,omitempty"`
	ScrapeConfigs []*ScrapeConfig `yaml:"scrape_configs,omitempty"`
}

type GlobalConfig struct {
	ScrapeInterval        model.Duration   `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout         model.Duration   `yaml:"scrape_timeout,omitempty"`
	EvaluationInterval    model.Duration   `yaml:"evaluation_interval,omitempty"`
	RuleQueryOffset       model.Duration   `yaml:"rule_query_offset,omitempty"`
	QueryLogFile          string           `yaml:"query_log_file,omitempty"`
	ScrapeFailureLogFile  string           `yaml:"scrape_failure_log_file,omitempty"`
	BodySizeLimit         units.Base2Bytes `yaml:"body_size_limit,omitempty"`
	SampleLimit           uint             `yaml:"sample_limit,omitempty"`
	LabelLimit            uint             `yaml:"label_limit,omitempty"`
	LabelNameLengthLimit  uint             `yaml:"label_name_length_limit,omitempty"`
	LabelValueLengthLimit uint             `yaml:"label_value_length_limit,omitempty"`
	KeepDroppedTargets    uint             `yaml:"keep_dropped_targets,omitempty"`
}

// 后续这里可以考虑更多的调优参数
type RuntimeConfig struct {
	GoGC int `yaml:"gogc,omitempty"`
}

// TODO: 后续需要考虑alert的配置
type AlertingConfig struct {
}

// TODO: 后续需要考虑rule的配置
type RuleConfig struct {
}

type ScrapeConfig struct {
	JobName                 string                  `yaml:"job_name"`
	HonorLabels             bool                    `yaml:"honor_labels,omitempty"`
	HonorTimestamps         bool                    `yaml:"honor_timestamps,omitempty"`
	ScrapeInterval          model.Duration          `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout           model.Duration          `yaml:"scrape_timeout,omitempty"`
	MetricsPath             string                  `yaml:"metrics_path,omitempty"`
	Scheme                  string                  `yaml:"scheme,omitempty"`
	Params                  url.Values              `yaml:"params,omitempty"`
	TargetLimit             int                     `yaml:"target_limit,omitempty"`
	LabelLimit              int                     `yaml:"label_limit,omitempty"`
	LabelNameLengthLimit    int                     `yaml:"label_name_length_limit,omitempty"`
	LabelValueLengthLimit   int                     `yaml:"label_value_length_limit,omitempty"`
	ServiceDiscoveryConfigs discovery.Configs       `yaml:"-"`
	HTTPClientConfig        config.HTTPClientConfig `yaml:",inline"`
	RelabelConfigs          []*relabel.Config       `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs    []*relabel.Config       `yaml:"metric_relabel_configs,omitempty"`
}

var (
	DefaultGlobalConfig = Config{
		GlobalConfig: GlobalConfig{
			ScrapeInterval: model.Duration(15 * time.Second),
			ScrapeTimeout:  model.Duration(15 * time.Second),
		},
		Scrape: DefaultScrapeConfig,
	}

	DefaultScrapeConfig = ScrapeConfig{
		MetricsPath: "/metrics",
		Scheme:      "http",
	}
)

func Load(s string) (*Config, error) {
	cfg := &Config{}
	*cfg = DefaultGlobalConfig
	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		logger.Error("Failed to parse YAML config",
			logger.Err(err),
			logger.String("content_preview", s[:min(len(s), 100)]))
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	return cfg, nil
}

func (sc *ScrapeConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*sc = DefaultScrapeConfig

	if err := discovery.UnmarshalYAMLWithInlineConfigs(sc, unmarshal); err != nil {
		return err
	}

	if len(sc.JobName) == 0 {
		return errors.New("job_name is empty")
	}

	if err := sc.HTTPClientConfig.Validate(); err != nil {
		return err
	}
	if len(sc.RelabelConfigs) == 0 {
		if err := checkStaticTargets(sc.ServiceDiscoveryConfigs); err != nil {
			return err
		}
	}

	for _, rlcfg := range sc.RelabelConfigs {
		if rlcfg == nil {
			return errors.New("empty or null target relabeling rule in scrape config")
		}
	}
	for _, rlcfg := range sc.MetricRelabelConfigs {
		if rlcfg == nil {
			return errors.New("empty or null metric relabeling rule in scrape config")
		}
	}
	return nil
}

func checkStaticTargets(configs discovery.Configs) error {
	for _, cfg := range configs {
		sc, ok := cfg.(discovery.StaticConfig)
		if !ok {
			continue
		}
		for _, tg := range sc {
			for _, t := range tg.Targets {
				if err := CheckTargetAddress(t[model.AddressLabel]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// CheckTargetAddress checks if target address is valid.
func CheckTargetAddress(address model.LabelValue) error {
	// For now check for a URL, we may want to expand this later.
	if strings.Contains(string(address), "/") {
		return fmt.Errorf("%q is not a valid hostname", address)
	}
	return nil
}
