package response

import (
	"fmt"
	"go-invoice/pkg/logger"
)

func Success(orderID, produk string, nominal int) {
	msg := fmt.Sprintf("[%s] %s | Rp%d", orderID, produk, nominal)
	logger.Success(msg)
}

func Error(orderID, context string, err error) {
	msg := fmt.Sprintf("[%s] %s", orderID, context)
	logger.Error(msg, err)
}

func Warn(orderID, message string) {
	msg := fmt.Sprintf("[%s] %s", orderID, message)
	logger.Warn(msg)
}

func Info(message string) {
	logger.Info(message)
}
