package oss

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	ossconfig "go-wind-admin/api/gen/go/oss_config/v1"
)

// ossYAMLWrapper 用于解析 oss.yaml 中嵌套的多供应商配置。
// YAML 结构示例：
//
//	oss:
//	  provider: "minio"
//	  minio:
//	    endpoint: "minio:9000"
//	    ...
//	  aliyun:
//	    endpoint: "oss-cn-hangzhou.aliyuncs.com"
//	    ...
type ossYAMLWrapper struct {
	OSS ossYAMLBody `yaml:"oss"`
}

type ossYAMLBody struct {
	Provider string        `yaml:"provider"`
	Minio    *s3ConfigYAML `yaml:"minio"`
	Aliyun   *s3ConfigYAML `yaml:"aliyun"`
	Tencent  *s3ConfigYAML `yaml:"tencent"`
	Qiniu    *s3ConfigYAML `yaml:"qiniu"`
	Huawei   *s3ConfigYAML `yaml:"huawei"`
	Baidu    *s3ConfigYAML `yaml:"baidu"`
	AWS      *s3ConfigYAML `yaml:"aws"`
	Google   *s3ConfigYAML `yaml:"google"`
	Azure    *s3ConfigYAML `yaml:"azure"`
}

type s3ConfigYAML struct {
	Endpoint     string `yaml:"endpoint"`
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	Token        string `yaml:"token"`
	UseSSL       bool   `yaml:"use_ssl"`
	Region       string `yaml:"region"`
	UploadHost   string `yaml:"upload_host"`
	DownloadHost string `yaml:"download_host"`
	Bucket       string `yaml:"bucket"`
}

// toProto 将 YAML 配置转为 protobuf S3CompatibleConfig
func (s *s3ConfigYAML) toProto() *ossconfig.S3CompatibleConfig {
	if s == nil {
		return nil
	}
	return &ossconfig.S3CompatibleConfig{
		Endpoint:     s.Endpoint,
		AccessKey:    s.AccessKey,
		SecretKey:    s.SecretKey,
		Token:        s.Token,
		UseSsl:       s.UseSSL,
		Region:       s.Region,
		UploadHost:   s.UploadHost,
		DownloadHost: s.DownloadHost,
		Bucket:       s.Bucket,
	}
}

// parseProviderType 将供应商名称字符串转为 OssProviderType 枚举
func parseProviderType(name string) ossconfig.OssProviderType {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "minio", "":
		return ossconfig.OssProviderType_OSS_PROVIDER_MINIO
	case "aliyun", "alibaba", "alicloud":
		return ossconfig.OssProviderType_OSS_PROVIDER_ALIYUN
	case "tencent", "cos":
		return ossconfig.OssProviderType_OSS_PROVIDER_TENCENT
	case "qiniu", "kodo":
		return ossconfig.OssProviderType_OSS_PROVIDER_QINIU
	case "huawei", "obs":
		return ossconfig.OssProviderType_OSS_PROVIDER_HUAWEI
	case "baidu", "bos":
		return ossconfig.OssProviderType_OSS_PROVIDER_BAIDU
	case "aws", "s3":
		return ossconfig.OssProviderType_OSS_PROVIDER_AWS
	case "google", "gcs":
		return ossconfig.OssProviderType_OSS_PROVIDER_GOOGLE
	case "azure", "blob":
		return ossconfig.OssProviderType_OSS_PROVIDER_AZURE
	default:
		return ossconfig.OssProviderType_OSS_PROVIDER_MINIO
	}
}

// LoadOssProviderConfig 从 oss.yaml 文件加载多供应商配置。
// configDir 为配置文件所在目录（如 configs/），函数会读取该目录下的 oss.yaml。
// 如果文件不存在或解析失败，返回 nil 和错误信息。
func LoadOssProviderConfig(configDir string) (*ossconfig.OssProviderConfig, error) {
	path := filepath.Join(configDir, "oss.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read oss.yaml: %w", err)
	}

	var wrapper ossYAMLWrapper
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse oss.yaml: %w", err)
	}

	body := wrapper.OSS
	pc := &ossconfig.OssProviderConfig{
		Provider: parseProviderType(body.Provider),
		Minio:    body.Minio.toProto(),
		Aliyun:   body.Aliyun.toProto(),
		Tencent:  body.Tencent.toProto(),
		Qiniu:    body.Qiniu.toProto(),
		Huawei:   body.Huawei.toProto(),
		Baidu:    body.Baidu.toProto(),
		Aws:      body.AWS.toProto(),
		Google:   body.Google.toProto(),
		Azure:    body.Azure.toProto(),
	}

	return pc, nil
}

// BuildBootstrapOSS 从 S3CompatibleConfig 构建兼容的 conf.OSS（用于保持下游代码兼容）。
// 当使用非 MinIO 供应商时，此函数将供应商配置映射到 conf.OSS_MinIO 结构。
func BuildBootstrapOSS(s3cfg *ossconfig.S3CompatibleConfig) *conf.OSS {
	if s3cfg == nil {
		return &conf.OSS{}
	}
	return &conf.OSS{
		Minio: &conf.OSS_MinIO{
			Endpoint:     s3cfg.GetEndpoint(),
			AccessKey:    s3cfg.GetAccessKey(),
			SecretKey:    s3cfg.GetSecretKey(),
			Token:        s3cfg.GetToken(),
			UseSsl:       s3cfg.GetUseSsl(),
			UploadHost:   s3cfg.GetUploadHost(),
			DownloadHost: s3cfg.GetDownloadHost(),
		},
	}
}
