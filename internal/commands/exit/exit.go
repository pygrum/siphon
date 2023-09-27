package exit

import (
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"os"
)

func ExitCmd() {
	logger.Notify("Goodbye!")
	_ = viper.WriteConfig()
	os.Exit(0)
}
