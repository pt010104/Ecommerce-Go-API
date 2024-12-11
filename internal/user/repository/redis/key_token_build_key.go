package redis

import "fmt"

func (r implRedis) buildSecretKeyKey(userID, sessionID string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", userKeyPrefix, userID, sessionKeyPrefix, sessionID, secretKeyPrefix)
}
