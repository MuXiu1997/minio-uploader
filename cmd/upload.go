package cmd

import (
	"fmt"
	"minio-uploader/internal/minioclient"
	"minio-uploader/internal/uploader"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	setUsageTemplate(cmdUpload)

	registerServerFlags(cmdUpload)
	rootCmd.AddCommand(cmdUpload)
}

var cmdUpload = &cobra.Command{
	Use:                   "upload [OPTIONS] FILE [FILE...]",
	DisableFlagsInUseLine: true,
	Aliases:               []string{"up"},
	Short:                 "Upload images",
	RunE: func(cmd *cobra.Command, args []string) error {
		v := viper.GetViper()
		folder := v.GetString("folder")

		minioClient, err := minioclient.NewMinioClient(v)
		if err != nil {
			return err
		}
		uploader, err := uploader.NewUploader(v, minioClient)
		if err != nil {
			return err
		}

		files := args
		for _, file := range files {
			returnUrl, e := uploader.Upload(folder, file)
			if e != nil {
				return e
			}
			fmt.Println(returnUrl)
		}

		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		BindFlags(cmd)
		if err := initConfig(); err != nil {
			return err
		}
		if len(args) == 0 {
			return fmt.Errorf(`"%s %s" requires at least 1 argument`, cmd.Root().Name(), cmd.Name())
		}
		return nil
	},
}

func registerServerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config", "c", "config.yml", "config file")
	cmd.PersistentFlags().StringP("folder", "f", "", "upload to this folder")
}

func BindFlags(cmd *cobra.Command) {
	_ = viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("folder", cmd.PersistentFlags().Lookup("folder"))
}

func initConfig() error {
	configFile := viper.GetString("config")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return fmt.Errorf("%s is not found: %v", configFile, err)
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("cannot read config file: %v", err)
		}
	}
	return nil
}
