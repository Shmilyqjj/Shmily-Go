package AppConfig

import (
	"time"
)

// 1. yaml配置文件获取

type Config struct {
	HttpClientConfig HttpClientConfig `yaml:"HttpClientConfig"`
	KafkaConfigs     []KafkaConfig    `yaml:"KafkaConfigs"`
	DorisConfigs     []DorisConfig    `yaml:"DorisConfigs"`
	EtlConfigs       []EtlConfig      `yaml:"EtlConfigs"`
	AlertWebHook     string           `yaml:"alert_web_hook"`
}

type HttpClientConfig struct {
	ReadTimeout         time.Duration `yaml:"read_timeout"`
	WriteTimeout        time.Duration `yaml:"write_timeout"`
	MaxIdleConnDuration time.Duration `yaml:"max_idle_conn_duration"`
	MaxConnsPerHost     int           `yaml:"max_conns_per_host"`
	Concurrency         int           `yaml:"http_concurrency"`
	DNSCacheDuration    time.Duration `yaml:"dns_cache_duration"`
}

type KafkaConfig struct {
	Cluster                string        `yaml:"cluster"`
	Brokers                string        `yaml:"brokers"`
	WriterBufferSize       int           `yaml:"writer_buffer_size"`
	WriterMaxRetries       int           `yaml:"writer_max_retries"`
	WriterBatchBytes       int64         `yaml:"writer_batch_bytes"`
	WriterBatchSize        int           `yaml:"writer_batch_size"`
	WriterBatchFrequency   time.Duration `yaml:"writer_batch_frequency"`
	WriterTimeout          time.Duration `yaml:"writer_timeout"`
	ReaderMaxBytes         int           `yaml:"reader_max_bytes"`
	ReaderQueueCapacity    int           `yaml:"reader_queue_capacity"`
	ReaderMaxWait          time.Duration `yaml:"reader_max_wait"`
	ReaderReadBatchTimeout time.Duration `yaml:"reader_read_batch_timeout"`
	ReaderMaxAttempts      int           `yaml:"reader_max_attempts"`
}

type DorisConfig struct {
	Cluster       string        `yaml:"cluster"`
	FENodes       string        `yaml:"fe_nodes"`
	Username      string        `yaml:"username"`
	Password      string        `yaml:"password"`
	BatchInterval time.Duration `yaml:"batch_interval"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
}

type EtlConfig struct {
	JobName                    string             `yaml:"job_name"`
	JobOwner                   string             `yaml:"job_owner"`
	SourceKafkaCluster         string             `yaml:"source_kafka_cluster"`
	SourceKafkaDataProtocol    string             `yaml:"source_kafka_data_protocol"`
	SourceKafkaDataExplodeKeys string             `yaml:"source_kafka_data_explode_keys"`
	SourceKafkaTopic           string             `yaml:"source_kafka_topic"`
	SourceKafkaGroup           string             `yaml:"source_kafka_group"`
	SinkDorisCluster           string             `yaml:"sink_doris_cluster"`
	SinkDorisTable             string             `yaml:"sink_doris_table"`
	SinkDorisColLower          bool               `yaml:"sink_doris_lowercase_column_enabled"`
	SinkDorisColumnAggr        DorisTableAggrType `yaml:"sink_doris_column_aggr"`
	SinkDorisColumnMapping     map[string]string  `yaml:"sink_doris_column_mapping"`
	SinkMaxRetries             int                `yaml:"sink_max_retries"`
	SinkRetryInterval          time.Duration      `yaml:"sink_retry_interval"`
	SinkBatchSize              int                `yaml:"sink_batch_size"`
	SinkBatchBytes             int                `yaml:"sink_batch_bytes"`
	SinkBatchInterval          time.Duration      `yaml:"sink_batch_interval"`
	SinkMaxFilterRatio         float64            `yaml:"sink_max_filter_ratio"`
}

type DorisTableAggrType struct {
	Enabled   bool              `yaml:"enabled"`
	AggrTypes map[string]string `yaml:"aggr_types"`
}

// 2. xml配置文件获取

type Users struct { // 结构体名为最外层xml的label
	Users []User `xml:"user"` // 这里为内层xml的label
}

type User struct {
	Id   int    `xml:"id"`
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}
