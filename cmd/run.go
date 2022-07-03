package cmd

import (
	"github.com/aliakbariaa1996/Calculate-Deliver-To-Destination/config"
	loggerx "github.com/aliakbariaa1996/Calculate-Deliver-To-Destination/internal/common/log"
	"github.com/aliakbariaa1996/Calculate-Deliver-To-Destination/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"os"
)

func init() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// Server is exported to make it graceful stop inside the main:
// cmd.Server.GracefulStop() for reconfiguration and code profiling.
var Server *grpc.Server

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "Run totem user profile",
	Long:  `Run totem user profile`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().String("port", "5050", "HTTP server listen address")
		cmd.Flags().String("config", "", "config file if present")
		cmd.Flags().String("deliver_man_loc", "", "config file if present")
		err := cmd.ParseFlags(args)
		if err != nil {
			return err
		}
		configFlag := cmd.Flags().Lookup("config")
		if configFlag != nil {
			configFilePath := configFlag.Value.String()
			if configFilePath != "" {
				viper.SetConfigFile(configFilePath)
				err := viper.ReadInConfig()
				if err != nil {
					return err
				}
			}
		}
		err = viper.BindPFlags(cmd.Flags())
		if err != nil {
			return err
		}
		return nil
	},
	RunE: runCmdE,
}

func intConfig() *config.Config {
	return &config.Config{
		Port:          viper.GetString("port"),
		DeliverManLoc: viper.GetString("deliver_man_loc"),
	}
}

func runCmdE(cmd *cobra.Command, args []string) error {
	cfg := intConfig()
	logger, err := loggerx.New("", "")

	if err = server.RunServer(cfg, logger); err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(runCMD)
}
