package cmd

import (
	"github.com/cidverse/reposync/pkg/config"
	"github.com/cidverse/reposync/pkg/repository"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{},
		Short:   `synchronizes repositories: clones new repos and pulls updates for existing ones`,
		Run: func(cmd *cobra.Command, args []string) {
			// flags
			silent, err := cmd.Flags().GetBool("silent")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse silent flag")
			}
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to parse dry-run flag")
			}

			// 第一步：执行clone（同步仓库列表）
			log.Info().Msg("step 1: cloning repositories")
			cloneArgs := []string{}
			if dryRun {
				cloneArgs = append(cloneArgs, "--dry-run")
			}
			if silent {
				cloneArgs = append(cloneArgs, "--silent")
			}
			cloneCommand := cloneCmd()
			cloneCommand.SetArgs(cloneArgs)
			if err := cloneCommand.Execute(); err != nil {
				log.Error().Err(err).Msg("clone step failed")
			}

			if dryRun {
				return
			}

			// 第二步：执行pull（更新代码）
			log.Info().Msg("step 2: pulling updates")
			pullArgs := []string{}
			if silent {
				pullArgs = append(pullArgs, "--silent")
			}
			pullCommand := pullCmd()
			pullCommand.SetArgs(pullArgs)
			if err := pullCommand.Execute(); err != nil {
				log.Error().Err(err).Msg("pull step failed")
			}

			// 第三步：清理不存在的仓库
			log.Info().Msg("step 3: cleaning up stale entries")
			stateFile := config.StateFile()
			state, err := config.LoadState(stateFile)
			if err != nil {
				log.Error().Err(err).Msg("failed to load state for cleanup")
				return
			}

			cleaned := 0
			for id, r := range state.Repositories {
				if !repository.Exists(r.Directory) {
					log.Info().Str("repository", r.Directory).Msg("removing stale repository from state")
					delete(state.Repositories, id)
					cleaned++
				}
			}

			if cleaned > 0 {
				if err := config.SaveState(stateFile, state); err != nil {
					log.Error().Err(err).Msg("failed to save state after cleanup")
				}
				log.Info().Int("count", cleaned).Msg("cleaned up stale entries")
			}

			log.Info().Msg("sync completed")
		},
	}

	cmd.PersistentFlags().BoolP("dry-run", "d", false, "dry run")
	cmd.PersistentFlags().BoolP("silent", "s", false, "silent (omit stdout/stderr output) from cli commands")

	return cmd
}
