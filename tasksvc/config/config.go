package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	SetupConfigPath()
}

func SetupConfigPath() {
	localConfigFileName := "config.yml" // Local config file
	// name should be set in the environment variable during execution in production and test
	_, testEnvFound := os.LookupEnv("tasksvc_TESTENV")
	_, prodEnvFound := os.LookupEnv("tasksvc_PRODENV")

	configFileName, configFileFound := os.LookupEnv("tasksvc_CONF")

	if !testEnvFound && !prodEnvFound {
		if !configFileFound {
			configFileName = localConfigFileName
		}
	} else if !configFileFound {
		panic("FATAL :: Config file env path tasksvc_CONF not set in production or test environment. Can't start service")
	}

	viper.SetConfigFile(configFileName)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("%s :: %s", fmt.Sprintf("FATAL :: Could not read config file [%s]. Can't start service. Check cxpsvc_CONF environment variable and the config file format are correct", configFileName), err.Error()))
	}
}

func (cfg *StartupConfig) ReadConfig() {
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalln("Cannot decode config file: ", err)
	}

	if len(cfg.TaskDB.User) == 0 && len(os.Getenv("TASKDB_MYSQL_USER")) > 0 {
		cfg.TaskDB.User = os.Getenv("TASKDB_MYSQL_USER")
	}

	if len(cfg.TaskDB.User) == 0 {
		log.Fatalln("TaskDB user not set in environment variable TASKDB_MYSQL_USER, can't start service")
	}

	if len(cfg.TaskDB.Password) == 0 && len(os.Getenv("TASKDB_MYSQL_PASSWORD")) > 0 {
		cfg.TaskDB.Password = os.Getenv("TASKDB_MYSQL_PASSWORD")
	}

	if len(cfg.TaskDB.Password) == 0 {
		log.Fatalln("TaskDB user not set in environment variable TASKDB_MYSQL_PASSWORD, can't start service")
	}
}
