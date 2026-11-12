package http

import (
	"strconv"
	"sync"
	"time"
)

type cachedPayload struct {
	data      interface{}
	expiresAt time.Time
}

var publicThemeResponseCache sync.Map

const publicThemeCacheTTL = 5 * time.Minute

func getCachedPublicThemePayload(key string) (interface{}, bool) {
	value, ok := publicThemeResponseCache.Load(key)
	if !ok {
		return nil, false
	}

	entry, ok := value.(cachedPayload)
	if !ok {
		publicThemeResponseCache.Delete(key)
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		publicThemeResponseCache.Delete(key)
		return nil, false
	}

	return entry.data, true
}

func setCachedPublicThemePayload(key string, payload interface{}) {
	publicThemeResponseCache.Store(key, cachedPayload{
		data:      payload,
		expiresAt: time.Now().Add(publicThemeCacheTTL),
	})
}

func invalidatePublicThemeCacheByRestaurant(restaurantID int64) {
	prefix := buildPublicThemeCachePrefix(restaurantID)
	publicThemeResponseCache.Range(func(key, _ interface{}) bool {
		cacheKey, ok := key.(string)
		if !ok {
			return true
		}
		if len(cacheKey) >= len(prefix) && cacheKey[:len(prefix)] == prefix {
			publicThemeResponseCache.Delete(cacheKey)
		}
		return true
	})
}

func buildPublicThemeCachePrefix(restaurantID int64) string {
	return "restaurant:" + strconv.FormatInt(restaurantID, 10) + ":"
}

func buildPublicThemeCacheKey(restaurantID int64, scope string) string {
	return buildPublicThemeCachePrefix(restaurantID) + scope
}
