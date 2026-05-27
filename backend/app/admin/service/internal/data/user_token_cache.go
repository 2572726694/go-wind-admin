package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authenticationV1 "go-wind-admin/api/gen/go/authentication/service/v1"
)

const (
	// AccessTokenKeyFormat 访问令牌键格式 at:{ct}:{uid}:{jti}
	AccessTokenKeyFormat = "at:%d:%d:%s"
	// AccessTokenKeyPrefixFormat 访问令牌键前缀格式 at:{ct}:{uid}:
	AccessTokenKeyPrefixFormat = "at:%d:%d:"

	// RefreshTokenKeyFormat 刷新令牌键格式 rt:{ct}:{uid}:{jti}
	RefreshTokenKeyFormat = "rt:%d:%d:%s"
	// RefreshTokenKeyPrefixFormat 刷新令牌键前缀格式 rt:{ct}:{uid}:
	RefreshTokenKeyPrefixFormat = "rt:%d:%d:"

	// BlacklistKeyFormat 访问令牌黑名单键格式 bl:{jti}
	BlacklistKeyFormat = "bl:%s"
)

// UserTokenCache 用户令牌缓存
type UserTokenCache struct {
	log *log.Helper
	rdb *redis.Client
}

func NewUserTokenCache(ctx *bootstrap.Context, rdb *redis.Client) *UserTokenCache {
	utc := &UserTokenCache{
		rdb: rdb,
		log: ctx.NewLoggerHelper("user-token/cache"),
	}
	return utc
}

// AddTokenPair 添加令牌对
func (r *UserTokenCache) AddTokenPair(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
	accessToken string,
	refreshToken string,
	accessTokenExpires time.Duration,
	refreshTokenExpires time.Duration,
) error {
	pipe := r.rdb.TxPipeline()

	atKey := r.makeAccessTokenKey(clientType, userId, jti)
	pipe.Set(ctx, atKey, accessToken, accessTokenExpires)

	rtKey := r.makeRefreshTokenKey(clientType, userId, jti)
	pipe.Set(ctx, rtKey, refreshToken, refreshTokenExpires)

	_, err := pipe.Exec(ctx)
	return err
}

// AddAccessToken 添加访问令牌
func (r *UserTokenCache) AddAccessToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
	accessToken string,
	expires time.Duration,
) error {
	key := r.makeAccessTokenKey(clientType, userId, jti)
	return r.set(ctx, key, accessToken, expires)
}

// AddRefreshToken 添加刷新令牌
func (r *UserTokenCache) AddRefreshToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
	refreshToken string,
	expires time.Duration,
) error {
	key := r.makeRefreshTokenKey(clientType, userId, jti)
	return r.set(ctx, key, refreshToken, expires)
}

// AddBlockedAccessToken 添加被阻止的访问令牌
func (r *UserTokenCache) AddBlockedAccessToken(ctx context.Context, jti string, reason string, expires time.Duration) error {
	key := r.makeBlacklistKey(jti)
	return r.set(ctx, key, reason, expires)
}

// GetAccessTokens 获取访问令牌
func (r *UserTokenCache) GetAccessTokens(ctx context.Context, clientType authenticationV1.ClientType, userId uint32) []string {
	prefix := r.makeAccessTokenPrefix(clientType, userId)
	return r.scanValues(ctx, prefix)
}

// GetRefreshTokens 获取刷新令牌
func (r *UserTokenCache) GetRefreshTokens(ctx context.Context, clientType authenticationV1.ClientType, userId uint32) []string {
	prefix := r.makeRefreshTokenPrefix(clientType, userId)
	return r.scanValues(ctx, prefix)
}

// RevokeToken 移除所有令牌
func (r *UserTokenCache) RevokeToken(ctx context.Context, clientType authenticationV1.ClientType, userId uint32) error {
	var err error
	if err = r.RevokeUserAllAccessToken(ctx, clientType, userId); err != nil {
		r.log.Errorf("remove user access token failed: [%v]", err)
	}

	if err = r.RevokeUserAllRefreshToken(ctx, clientType, userId); err != nil {
		r.log.Errorf("remove user refresh token failed: [%v]", err)
	}

	return err
}

func (r *UserTokenCache) RevokeTokenByJti(ctx context.Context, clientType authenticationV1.ClientType, userId uint32, jti string) error {
	var err error
	if err = r.RevokeAccessToken(ctx, clientType, userId, jti); err != nil {
		r.log.Errorf("remove user access token failed: [%v]", err)
	}

	if err = r.RevokeRefreshToken(ctx, clientType, userId, jti); err != nil {
		r.log.Errorf("remove user refresh token failed: [%v]", err)
	}

	return err
}

// RevokeAccessToken 移除访问令牌
func (r *UserTokenCache) RevokeAccessToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
) error {
	key := r.makeAccessTokenKey(clientType, userId, jti)
	return r.del(ctx, key)
}

// RevokeRefreshToken 移除刷新令牌
func (r *UserTokenCache) RevokeRefreshToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
) error {
	key := r.makeRefreshTokenKey(clientType, userId, jti)
	return r.del(ctx, key)
}

// RevokeBlockedAccessToken 撤销被阻止的访问令牌
func (r *UserTokenCache) RevokeBlockedAccessToken(ctx context.Context, jti string) error {
	key := r.makeBlacklistKey(jti)
	return r.del(ctx, key)
}

// IsValidAccessToken 访问令牌是否有效
func (r *UserTokenCache) IsValidAccessToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
	uploadedToken string,
) (bool, error) {
	key := r.makeAccessTokenKey(clientType, userId, jti)

	storedToken, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil // 令牌不存在或已过期
	}
	if err != nil {
		return false, err
	}

	return storedToken == uploadedToken, nil
}

func (r *UserTokenCache) IsExistAccessTokenByJti(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
) (bool, error) {
	key := r.makeAccessTokenKey(clientType, userId, jti)
	return r.exists(ctx, key), nil
}

// IsExistAccessToken 访问令牌是否存在（按令牌值查找，返回 jti）
func (r *UserTokenCache) IsExistAccessToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	uploadedToken string,
) (exist bool, jti string, err error) {
	prefix := r.makeAccessTokenPrefix(clientType, userId)
	return r.scanFindValue(ctx, prefix, uploadedToken)
}

// IsExistRefreshToken 刷新令牌是否存在（按令牌值查找，返回 jti）
func (r *UserTokenCache) IsExistRefreshToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	uploadedToken string,
) (exist bool, jti string, err error) {
	prefix := r.makeRefreshTokenPrefix(clientType, userId)
	return r.scanFindValue(ctx, prefix, uploadedToken)
}

// IsValidRefreshToken 刷新令牌是否有效
func (r *UserTokenCache) IsValidRefreshToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
	jti string,
	uploadedToken string,
) (bool, error) {
	key := r.makeRefreshTokenKey(clientType, userId, jti)

	storedToken, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil // 令牌不存在或已过期
	}
	if err != nil {
		return false, err
	}

	return storedToken == uploadedToken, nil
}

// IsBlockedAccessToken 访问令牌是否被阻止
func (r *UserTokenCache) IsBlockedAccessToken(ctx context.Context, jti string) bool {
	key := r.makeBlacklistKey(jti)
	return r.exists(ctx, key)
}

// RevokeUserAllAccessToken 删除用户所有访问令牌
func (r *UserTokenCache) RevokeUserAllAccessToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
) error {
	prefix := r.makeAccessTokenPrefix(clientType, userId)
	return r.delByPattern(ctx, prefix)
}

// RevokeUserAllRefreshToken 删除用户所有刷新令牌
func (r *UserTokenCache) RevokeUserAllRefreshToken(
	ctx context.Context,
	clientType authenticationV1.ClientType,
	userId uint32,
) error {
	prefix := r.makeRefreshTokenPrefix(clientType, userId)
	return r.delByPattern(ctx, prefix)
}

// makeAccessTokenKey 生成访问令牌键 at:{ct}:{uid}:{jti}
func (r *UserTokenCache) makeAccessTokenKey(clientType authenticationV1.ClientType, userId uint32, jti string) string {
	return fmt.Sprintf(AccessTokenKeyFormat, clientType.Number(), userId, jti)
}

// makeAccessTokenPrefix 生成访问令牌键前缀 at:{ct}:{uid}:
func (r *UserTokenCache) makeAccessTokenPrefix(clientType authenticationV1.ClientType, userId uint32) string {
	return fmt.Sprintf(AccessTokenKeyPrefixFormat, clientType.Number(), userId)
}

// makeRefreshTokenKey 生成刷新令牌键 rt:{ct}:{uid}:{jti}
func (r *UserTokenCache) makeRefreshTokenKey(clientType authenticationV1.ClientType, userId uint32, jti string) string {
	return fmt.Sprintf(RefreshTokenKeyFormat, clientType.Number(), userId, jti)
}

// makeRefreshTokenPrefix 生成刷新令牌键前缀 rt:{ct}:{uid}:
func (r *UserTokenCache) makeRefreshTokenPrefix(clientType authenticationV1.ClientType, userId uint32) string {
	return fmt.Sprintf(RefreshTokenKeyPrefixFormat, clientType.Number(), userId)
}

// makeBlacklistKey 生成黑名单键
func (r *UserTokenCache) makeBlacklistKey(jti string) string {
	return fmt.Sprintf(BlacklistKeyFormat, jti)
}

// extractJtiFromKey 从键中提取 jti，如 "at:0:1:abc-def" → "abc-def"
func (r *UserTokenCache) extractJtiFromKey(key, prefix string) string {
	return strings.TrimPrefix(key, prefix)
}

func (r *UserTokenCache) set(ctx context.Context, key string, value string, expires time.Duration) error {
	if err := r.rdb.Set(ctx, key, value, expires).Err(); err != nil {
		r.log.Errorf("set key[%s] failed: %v", key, err)
		return err
	}
	return nil
}

func (r *UserTokenCache) get(ctx context.Context, key string) string {
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ""
		}

		r.log.Errorf("get key[%s] failed: %v", key, err)
		return ""
	}
	return result
}

// del 删除键
func (r *UserTokenCache) del(ctx context.Context, key string) error {
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		r.log.Errorf("del key[%s] failed: %v", key, err)
		return err
	}
	return nil
}

func (r *UserTokenCache) exists(ctx context.Context, key string) bool {
	n, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		r.log.Errorf("exists key[%s] failed: %v", key, err)
		return false
	}
	return n > 0
}

// scanKeys 使用 SCAN 命令扫描匹配前缀的所有键
func (r *UserTokenCache) scanKeys(ctx context.Context, prefix string) []string {
	var keys []string
	pattern := prefix + "*"
	var cursor uint64
	for {
		var batch []string
		var err error
		batch, cursor, err = r.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			r.log.Errorf("scan pattern[%s] failed: %v", pattern, err)
			return keys
		}
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}
	return keys
}

// scanValues 扫描前缀匹配的所有键并返回值
func (r *UserTokenCache) scanValues(ctx context.Context, prefix string) []string {
	keys := r.scanKeys(ctx, prefix)
	if len(keys) == 0 {
		return []string{}
	}

	var values []string
	for _, key := range keys {
		val, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue // 已过期，跳过
			}
			r.log.Errorf("get key[%s] failed: %v", key, err)
			continue
		}
		values = append(values, val)
	}
	return values
}

// scanFindValue 扫描前缀匹配的所有键，查找值等于 targetValue 的条目，返回 (exist, jti, err)
func (r *UserTokenCache) scanFindValue(ctx context.Context, prefix, targetValue string) (bool, string, error) {
	keys := r.scanKeys(ctx, prefix)
	for _, key := range keys {
		val, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			r.log.Errorf("get key[%s] failed: %v", key, err)
			return false, "", err
		}
		if val == targetValue {
			jti := r.extractJtiFromKey(key, prefix)
			return true, jti, nil
		}
	}
	return false, "", nil
}

// delByPattern 按前缀批量删除所有匹配的键
func (r *UserTokenCache) delByPattern(ctx context.Context, prefix string) error {
	keys := r.scanKeys(ctx, prefix)
	if len(keys) == 0 {
		return nil
	}
	if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
		r.log.Errorf("del keys by prefix[%s] failed: %v", prefix, err)
		return err
	}
	return nil
}
