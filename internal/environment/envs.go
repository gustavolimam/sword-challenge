package environment

import (
	"fmt"
	"os"
)

const (
	Port           = "PORT"
	TaskQueue      = "TASK_QUEUE"
	RabbitHost     = "RABBIT_HOST"
	RabbitUsername = "RABBIT_USERNAME"
	RabbitPassword = "RABBIT_PASSWORD"
	DBUser         = "DB_USER"
	DBPassword     = "DB_PASSWORD"
	DBHost         = "DB_HOST"
	DBName         = "DB_NAME"
)

func CheckEnvVars() {
	envVars := []string{
		Port,
		TaskQueue,
		RabbitHost,
		RabbitUsername,
		RabbitPassword,
		DBUser,
		DBPassword,
		DBHost,
		DBName,
	}

	for _, v := range envVars {
		fmt.Println("CHEGOU", v, os.Getenv(v))
		if os.Getenv(v) == "" {
			panic(fmt.Sprintf("env variable %s must be defined", v))
		}
	}
}
