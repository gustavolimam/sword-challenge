package models

import "fmt"

type Credential struct {
	Host     string
	User     string
	Password string
	Vhost    *string
}

func (credential Credential) GetConnectionString() string {

	return fmt.Sprintf("amqp://%s:%s@%s", credential.User, credential.Password, credential.Host)
}
