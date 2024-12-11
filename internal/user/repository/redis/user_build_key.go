package redis

import "fmt"

func (r implRedis) buildUserKey(userID string) string {
	return fmt.Sprintf("%s:%s", userKeyPrefix, userID)
}
