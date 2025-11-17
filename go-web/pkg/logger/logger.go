package logger

import (
	"go.uber.org/zap"
)

func InitLogger()  {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Error("初始化日志记录器失败", zap.Error(err))
	}
}