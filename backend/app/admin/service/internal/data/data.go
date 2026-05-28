package data

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tx7do/go-utils/captcha"
	"github.com/tx7do/go-utils/password"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"

	authenticationV1 "go-wind-admin/api/gen/go/authentication/service/v1"

	"go-wind-admin/pkg/oss"
	"go-wind-admin/pkg/serviceid"
)

func NewClientType() authenticationV1.ClientType {
	return authenticationV1.ClientType_admin
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(ctx *bootstrap.Context) (*redis.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil, func() {}, nil
	}

	l := ctx.NewLoggerHelper("redis/data/admin-service")

	cli := redisClient.NewClient(cfg.Data, l)

	return cli, func() {
		if err := cli.Close(); err != nil {
			l.Error(err)
		}
	}, nil
}

func NewMinIoClient(ctx *bootstrap.Context) *oss.MinIOClient {
	return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
}

// NewOSSClient 根据多供应商配置创建 OSS 客户端。
// 自动从 oss.yaml 加载 provider 字段，选择对应的供应商配置初始化客户端。
// 如果加载失败或供应商为 MinIO，则回退到默认的 MinIO 客户端。
func NewOSSClient(ctx *bootstrap.Context) *oss.MinIOClient {
	l := ctx.NewLoggerHelper("oss/data/admin-service")

	configDir := resolveConfigDir()
	if configDir == "" {
		l.Warn("Config directory not found, falling back to default MinIO client")
		return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
	}

	pc, err := oss.LoadOssProviderConfig(configDir)
	if err != nil {
		l.Warnf("Failed to load OssProviderConfig: %v, falling back to default MinIO client", err)
		return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
	}

	// MinIO 供应商直接使用 bootstrap 配置（保持原有行为）
	if pc.GetProvider() == 0 { // OSS_PROVIDER_MINIO
		return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
	}

	client, err := oss.NewOSSClientFromProviderConfig(pc, ctx.GetConfig().Oss, ctx.GetLogger())
	if err != nil {
		l.Errorf("Failed to create OSS client from provider config: %v, falling back to default MinIO client", err)
		return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
	}

	return client
}

// resolveConfigDir 从命令行参数或约定路径查找 configs 目录
func resolveConfigDir() string {
	// 1. 从 -conf 命令行参数解析配置目录
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-conf=") {
			return strings.TrimPrefix(arg, "-conf=")
		}
		if arg == "-conf" && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
	}

	// 2. 尝试约定路径
	candidates := []string{
		"configs",
		"./configs",
		"app/admin/service/configs",
	}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			return abs
		}
	}

	return ""
}

func NewPasswordCrypto() password.Crypto {
	crypto, err := password.CreateCrypto("bcrypt")
	if err != nil {
		panic(err)
	}
	return crypto
}

func NewCaptcha(rdb *redis.Client) *captcha.Captcha {
	captchaInstance := captcha.NewCaptcha(rdb,
		captcha.WithDriverType(captcha.DriverString),
		captcha.WithExpire(10*time.Minute),
		captcha.WithKeyPrefix(serviceid.ProjectName+":captcha"),
		captcha.WithStringCount(6),
		captcha.WithStringSource("ABCDEFGHJKLMNPQRSTUVWXYZ23456789"),
	)
	return captchaInstance
}
