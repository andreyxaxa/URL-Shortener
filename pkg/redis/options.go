package redis

import "time"

type Option func(*Client)

func DB(db int) Option {
	return func(c *Client) {
		c.db = db
	}
}

func MaxRetries(retries int) Option {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

func User(user string) Option {
	return func(c *Client) {
		c.user = user
	}
}

func Password(password string) Option {
	return func(c *Client) {
		c.password = password
	}
}

func DialTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.dialTimeout = timeout
	}
}

func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}
