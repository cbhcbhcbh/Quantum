package infra

import (
	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	"github.com/gocql/gocql"
)

var CassandraSession *gocql.Session

func NewCassandraSession(config *config.Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(util.GetServerAddrs(config.Cassandra.Hosts)...)
	cluster.Port = config.Cassandra.Port
	cluster.Keyspace = config.Cassandra.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: 3,
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Cassandra.User,
		Password: config.Cassandra.Password,
	}
	cluster.DefaultIdempotence = false
	cluster.NumConns = 3
	CassandraSession, err := cluster.CreateSession()
	return CassandraSession, err
}
