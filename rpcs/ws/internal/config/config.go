package config

type Config struct {
	Name     string
	ListenOn string
	Pattern  string
	Jwt      struct {
		AccessSecret string
		AccessExpire int64
	}
}
