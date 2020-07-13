package core

import (
	"github.com/corrots/go-crontab/common"
)

func GetLogList(name string, skip, limit int64) ([]common.Log, error) {
	mongo, err := common.NewMongo()
	if err != nil {
		return nil, err
	}
	return mongo.GetLogs(name, skip, limit)
}
