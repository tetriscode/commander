package util

import (
	"os"

	"github.com/arichardet/grammar-log/logger"
)

var Log *logger.Logger

func init() {
	Log = logger.NewLogger("commander", os.Stdout)
}
