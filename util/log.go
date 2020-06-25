package util

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
)

func FileLogHook(filename string) (func(zapcore.Entry) error, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewWriter(file)

	return func(e zapcore.Entry) error {
		_, err := buf.WriteString(fmt.Sprintf("%+v", e))
		if err != nil {
			return err
		}
		return buf.Flush()
	}, nil
}
