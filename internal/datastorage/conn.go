package datastorage

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"l0/utils"
	"log"
	"os"
	"time"
)

type ConnectionConfig struct {
	user, pass, host, port, db string
	maxTry                     int
	defaultSleep               int
}

func (cc *ConnectionConfig) GettingEnv() {
	cc.user = os.Getenv("PSQL_DB_USER")
	cc.pass = os.Getenv("PSQL_DB_PASS")
	cc.host = os.Getenv("PSQL_DB_HOST")
	cc.port = os.Getenv("PSQL_DB_PORT")
	cc.db = os.Getenv("PSQL_DB_NAME")
	cc.maxTry = utils.GetIntEnv("PSQL_CONN_MAX_TRY", 5)
	cc.defaultSleep = utils.GetIntEnv("PSQL_CONN_MAX_DELAY", 3)
}

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func CreateClient(ctx context.Context, cc ConnectionConfig) (err error, pool *pgxpool.Pool) {
	cc.GettingEnv()
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cc.user, cc.pass, cc.host, cc.port, cc.db)
	fmt.Println(connString)

	//Подключение в цикле под развертывание в контейнере
	err = utils.TryToConn(func() error {
		ctx, escape := context.WithTimeout(ctx, time.Duration(cc.defaultSleep)*time.Second)
		defer escape()
		pool, err = pgxpool.Connect(ctx, connString)
		if err != nil {
			log.Fatalf("error tries to connection: %s", err)
			return err
		}
		return nil
	}, cc.maxTry, time.Duration(cc.defaultSleep)*time.Second)

	return
}
