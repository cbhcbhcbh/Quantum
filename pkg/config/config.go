package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Chat      *ChatConfig      `mapstructure:"chat"`
	User      *UserConfig      `mapstructure:"user"`
	Forwarder *ForwarderConfig `mapstructure:"forwarder"`
	Redis     *RedisConfig     `mapstructure:"redis"`
	Kafka     *KafkaConfig     `mapstructure:"kafka"`
	Cassandra *CassandraConfig `mapstructure:"cassandra"`
}

type ChatConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
			Swag    bool
		}
	}
	Grpc struct {
		Server struct {
			Port string
		}
		Client struct {
			User struct {
				Endpoint string
			}
			Forwarder struct {
				Endpoint string
			}
		}
	}
	Subscriber struct {
		Id string
	}
	Message struct {
		MaxNum        int64
		PaginationNum int
		MaxSizeByte   int64
	}
	JWT struct {
		Secret           string
		ExpirationSecond int64
	}
}

type UserConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
			Swag    bool
		}
	}
	Grpc struct {
		Server struct {
			Port string
		}
	}
	OAuth struct {
		Google struct {
			ClientID     string
			ClientSecret string
			RedirectUrl  string
			Scopes       []string
		}
	}
}

type ForwarderConfig struct {
	Grpc struct {
		Server struct {
			Port string
		}
	}
}

type KafkaConfig struct {
	Addrs   string
	Version string
}

type RedisConfig struct {
	Password                string
	Addrs                   string
	ExpirationHour          int64
	MinIdleConn             int
	PoolSize                int
	ReadTimeoutMilliSecond  int64
	WriteTimeoutMilliSecond int64
}

type CassandraConfig struct {
	Hosts    string
	Port     int
	User     string
	Password string
	Keyspace string
}

func setDefault() {
	viper.SetDefault("chat.http.server.port", "5001")
	viper.SetDefault("chat.http.server.maxConn", 200)
	viper.SetDefault("chat.http.server.swag", false)
	viper.SetDefault("chat.grpc.server.port", "4000")
	viper.SetDefault("chat.grpc.client.user.endpoint", "localhost:4001")
	viper.SetDefault("chat.grpc.client.forwarder.endpoint", "localhost:4002")
	viper.SetDefault("chat.subscriber.id", "rc.msg."+os.Getenv("HOSTNAME"))
	viper.SetDefault("chat.message.maxNum", 5000)
	viper.SetDefault("chat.message.paginationNum", 5000)
	viper.SetDefault("chat.message.maxSizeByte", 4096)
	viper.SetDefault("chat.jwt.secret", "replaceme")
	viper.SetDefault("chat.jwt.expirationSecond", 86400)

	viper.SetDefault("forwarder.grpc.server.port", "4002")

	viper.SetDefault("kafka.addrs", "localhost:9092")
	viper.SetDefault("kafka.version", "1.0.0")

	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.addrs", "localhost:6379")
	viper.SetDefault("redis.expirationHour", 24)
	viper.SetDefault("redis.minIdleConn", 16)
	viper.SetDefault("redis.poolSize", 64)
	viper.SetDefault("redis.readTimeoutMilliSecond", 3000)
	viper.SetDefault("redis.writeTimeoutMilliSecond", 3000)

	viper.SetDefault("cassandra.hosts", "localhost")
	viper.SetDefault("cassandra.port", 9042)
	viper.SetDefault("cassandra.user", "cassandra")
	viper.SetDefault("cassandra.password", "cassandra")
	viper.SetDefault("cassandra.keyspace", "randomchat")

}

func NewConfig() (*Config, error) {
	setDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
