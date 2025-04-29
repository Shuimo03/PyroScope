package config

// ref: https://github.com/prometheus/prometheus/blob/main/config/config.go
import (
	"github.com/alecthomas/units"
	"github.com/prometheus/common/model"
)

type Config struct {
	GlobalConfig GlobalConfig  `yaml:"global"`
	Runtime      RuntimeConfig `yaml:"runtime,omitempty"`
	Alert        AlertConfig   `yaml:"alert,omitempty"`
	Rule         RuleConfig    `yaml:"rule,omitempty"`
	Scrape       ScrapeConfig  `yaml:"scrape,omitempty"`
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

type AlertConfig struct {
}

type RuleConfig struct {
}

type ScrapeConfig struct {
}

func LoadConfig() (*Config, error) {
	return &Config{}, nil
}
