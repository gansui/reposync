package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cidverse/cidverseutils/zerologconfig"
	"github.com/cidverse/reposync/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfg zerologconfig.LogConfig
var configFile string
var logFileDir string

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  `reposync`,
		Long: `A cli tool to mirror/sync projects onto the local file system (and/or merge content of specific folders to aggregate ie. doc files)`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 先配置日志级别和格式
			zerologconfig.Configure(cfg)

			// 配置日志文件输出
			if logFileDir != "" {
				// 获取日志文件目录
				dir := logFileDir
				if dir == "" {
					dir = config.GetConfigDir(configFile)
				}

				// 生成日志文件名
				logFileName := "reposync-" + time.Now().Format("20060102-1504") + ".log"
				logPath := filepath.Join(dir, logFileName)

				// 创建日志文件
				logFile, err := os.Create(logPath)
				if err != nil {
					log.Fatal().Err(err).Str("path", logPath).Msg("failed to create log file")
				}

				// 使用 MultiLevelWriter 同时输出到 stderr 和文件
				syncWriter := zerolog.SyncWriter(logFile)
				multiWriter := zerolog.MultiLevelWriter(
					zerolog.ConsoleWriter{Out: os.Stderr},
					syncWriter,
				)
				log.Logger = log.Output(multiWriter)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(0)
		},
	}

	cmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "log level - allowed: "+strings.Join(zerologconfig.ValidLogLevels, ","))
	cmd.PersistentFlags().StringVar(&cfg.LogFormat, "log-format", "color", "log format - allowed: "+strings.Join(zerologconfig.ValidLogFormats, ","))
	cmd.PersistentFlags().BoolVar(&cfg.LogCaller, "log-caller", false, "include caller in log functions")
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	cmd.PersistentFlags().StringVar(&logFileDir, "log-file", "", "save logs to file in specified directory (default: same dir as config.yaml)")

	cmd.AddCommand(indexCmd())
	cmd.AddCommand(cloneCmd())
	cmd.AddCommand(pullCmd())
	cmd.AddCommand(syncCmd())
	cmd.AddCommand(houseKeepingCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(versionCmd())

	return cmd
}

// Execute executes the root command.
func Execute() error {
	return rootCmd().Execute()
}
