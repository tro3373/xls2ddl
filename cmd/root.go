package cmd

import (
	"os"

	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xls2ddl",
	Short: "Generate MySql DDL from xls",
	Long: `Generate MySql DDL sql from defined in xls.

Define table cell in application config file(.xls2ddl)

## Sample .xls2ddl

IgnoreSheet:
  - 型一覧
  - '.*データ'
TableOption: 'ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci'
TableCell: W1
TableNameCell: W2
AdditionalOptionCell: BB1
ColumnStartRow: 6
ColumnCol: C
ColumnNameCol: K
TypeCol: S
DigitCol: W # 桁
DecimalCol: Y # 少数
NNCol: AA
PKCol: AD
DescriptionCol: AF
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, files []string) {
		err := Gen(config, files)
		if err != nil {
			log.Error("Failed to gen ddl", err)
			os.Exit(1)
		}
		log.Info("done")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string
var config Config

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xls2ddl)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	err := initConfigInner()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfigInner() error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".xls2ddl" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		// viper.AddConfigPath("..")
		viper.SetConfigName(".xls2ddl")
	}
	viper.SetConfigType("yml")

	viper.AutomaticEnv() // read in environment variables that match

	log.Debug(">>>>>> os.Args", os.Args)
	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Debugf("> Loading config from %s.\n", viper.ConfigFileUsed())
	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	log.Debug("> Loaded config:", config)
	return nil
}
