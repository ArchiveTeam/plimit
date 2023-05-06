package limitmgr

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (l *LimitManager) CollectGarbage(ctx context.Context) {
	t := time.Now().Unix()
	count, err := l.rdb.ZRemRangeByScore(ctx, l.connectionsKey, "-inf", fmt.Sprintf("%v", t)).Result()
	if err != nil {
		log.Fatalf("Failed to collect garbage: %s", err)
	} else {
		log.Printf("Cleaned up %v stale connections.", count)
	}
}
