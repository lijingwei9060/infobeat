package command

import "time"

// Config command的配置信息
type Config struct {
	TimeOut time.Duration `config:"timeout" validate:"required,min=0,nonzero"` //超时时间
	Backoff int           `config:"backoff" validate:"required,min=1"`         //并发度
}
