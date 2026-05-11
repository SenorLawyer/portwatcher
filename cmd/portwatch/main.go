package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/SenorLawyer/portwatcher/internal/actions"
	"github.com/SenorLawyer/portwatcher/internal/app"
	"github.com/SenorLawyer/portwatcher/internal/config"
	dockerenricher "github.com/SenorLawyer/portwatcher/internal/docker"
	"github.com/SenorLawyer/portwatcher/internal/history"
	"github.com/SenorLawyer/portwatcher/internal/scanner"
	"github.com/SenorLawyer/portwatcher/internal/tui"
	"github.com/SenorLawyer/portwatcher/internal/version"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cfg := config.Default()

	cmd := &cobra.Command{
		Use:          "portwatch",
		Short:        "A fast TUI for local ports and processes",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := history.Open(cfg.HistoryPath, cfg.HistoryRetention)
			if err != nil {
				return err
			}
			defer store.Close()

			svc := app.New(scanner.New(dockerenricher.New(cfg.Docker)), store)
			model := tui.New(svc, actions.Real{}, cfg)
			program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
			_, err = program.Run()
			return err
		},
	}

	cmd.PersistentFlags().DurationVar(&cfg.Interval, "interval", cfg.Interval, "refresh interval")
	cmd.PersistentFlags().BoolVar(&cfg.Docker, "docker", cfg.Docker, "enrich rows with Docker container mappings")
	cmd.PersistentFlags().StringVar(&cfg.HistoryPath, "history", cfg.HistoryPath, "history JSONL path")
	cmd.PersistentFlags().DurationVar(&cfg.HistoryRetention, "history-retention", cfg.HistoryRetention, "history retention window")

	cmd.AddCommand(listCmd(&cfg), historyCmd(&cfg), versionCmd())
	return cmd
}

func listCmd(cfg *config.Config) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Print current local ports",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
			defer cancel()
			snap, err := scanner.New(dockerenricher.New(cfg.Docker)).Snapshot(ctx)
			if err != nil {
				return err
			}
			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(snap)
			}
			fmt.Printf("%-5s %-26s %-26s %-12s %-7s %-22s %s\n", "NET", "LOCAL", "REMOTE", "STATE", "PID", "PROCESS", "COMMAND")
			for _, row := range snap.Ports {
				fmt.Printf("%-5s %-26s %-26s %-12s %-7d %-22s %s\n", row.Protocol, row.Address, row.Remote, row.State, row.PID, row.Process, row.Command)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "emit JSON")
	return cmd
}

func historyCmd(cfg *config.Config) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Print port change history",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := history.Open(cfg.HistoryPath, cfg.HistoryRetention)
			if err != nil {
				return err
			}
			defer store.Close()
			events, err := store.ReadAll()
			if err != nil {
				return err
			}
			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(events)
			}
			for _, event := range events {
				fmt.Printf("%s %-7s %-5s %-22s pid=%d %s\n", event.At.Format(time.RFC3339), event.Type, event.Port.Protocol, event.Port.Address, event.Port.PID, event.Port.Command)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "emit JSON")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.String())
		},
	}
}
