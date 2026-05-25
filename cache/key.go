package cache

import (
	"fmt"
	"strconv"
)

const (
	//RankKey 每日排名
	RankKey = "rank"
	//BookRank 图书排名
	BookRank = "BookRank"
	//Camera 相机排名
	CameraRank = "CameraRank"
)

func GetViewKey(id uint) string {
	return fmt.Sprintf("views:product:%s", strconv.Itoa(int(id)))
}

func GetProductLockKey(id uint) string {
	return fmt.Sprintf("lock:product:%d", id)
}
