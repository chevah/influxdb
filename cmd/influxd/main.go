package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/cmd/influxd/generate"
	"github.com/influxdata/influxdb/cmd/influxd/inspect"
	"github.com/influxdata/influxdb/cmd/influxd/launcher"
	_ "github.com/influxdata/influxdb/query/builtin"
	_ "github.com/influxdata/influxdb/tsdb/tsi1"
	_ "github.com/influxdata/influxdb/tsdb/tsm1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = fmt.Sprint(time.Now().UTC().Format(time.RFC3339))
)

var rootCmd = &cobra.Command{
	Use:   "influxd",
	Short: "Influx Server",
}

func init() {
	influxdb.SetBuildInfo(version, commit, date)
	viper.SetEnvPrefix("INFLUXD")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	rootCmd.InitDefaultHelpCmd()
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the influxd server version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("InfluxDB %s (git: %s) build_date: %s\n", version, commit, date)
		},
	})
	rootCmd.AddCommand(launcher.NewCommand())
	rootCmd.AddCommand(generate.Command)
	rootCmd.AddCommand(inspect.NewCommand())
}

// find determines the default behavior when running influxd.
// Specifically, find will return the influxd run command if no sub-command
// was specified.
func find(args []string) *cobra.Command {
	cmd, _, err := rootCmd.Find(args)
	if err == nil && cmd == rootCmd {
		// Execute the run command if no sub-command is specified
		return launcher.NewCommand()
	}

	return rootCmd
}

func main() {
	runtime.SetBlockProfileRate(int(time.Second))
	runtime.SetMutexProfileFraction(1)
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	
	cmd := find(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
