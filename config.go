package main

type Config struct {
	Ports struct {
		HTTP string `envconfig:"HTTP_PORT" default:":8080"`
	}
}

func (c Config) Validate() error {
	return nil
}
