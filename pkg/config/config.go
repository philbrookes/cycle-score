package config

import (
	"strconv"
)

type Config struct {
	clientSecret string
	clientId int
	port int
	cookieLifetime int
	oauthCallbackUrl string
	allowedMethods []string
	allowedOrigins []string
}

//GetConfig returns a populated config object
func GetConfig() *Config {
	config := Config{
		clientSecret: "32f695b3558f66b57fdb8a1dbcb46e296842a25e",
		clientId: 17697,
		port: 8080,
		cookieLifetime: 1000,
		oauthCallbackUrl: "http://www.cycle-score.com:8080/api/auth/callback",
		allowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		allowedOrigins: []string{"http://localhost:8080", "http://127.0.0.1:8080"},
	}
	return &config
}

func (c *Config) GetClientId() (int) {
	return c.clientId
}

func (c *Config) GetClientSecret() (string) {
	return c.clientSecret
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) GetPortListenerStr() string {
	return ":" + strconv.Itoa(c.port)
}

func (c *Config) GetOAuthCallbackUrl() string {
	return c.oauthCallbackUrl
}

func (c *Config) GetAllowedMethods() []string {
	return c.allowedMethods
}

func (c *Config) GetAllowedOrigins() []string {
	return c.allowedOrigins
}

func (c *Config) GetCookieLifetime() int {
	return c.cookieLifetime
}