package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Moriartii/url-shortner-api/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	postgresUsername = "db_username"
	postgresPassword = "db_password"
	postgresHost     = "db_host"
	postgresDatabase = "db_database"
	postgresSslmode  = "db_sslmode"
)

var (
	Client *sql.DB

	username = os.Getenv(postgresUsername)
	password = os.Getenv(postgresPassword)
	host     = os.Getenv(postgresHost)
	database = os.Getenv(postgresDatabase)
	sslmode  = os.Getenv(postgresSslmode)

	log *zap.SugaredLogger
)

func init() {

	log = logger.GetLogger().Named("postgres (pg_init.go)").Sugar()

	datasourceName := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		username,
		password,
		host,
		database,
		sslmode,
	)

	var err error
	Client, err = sql.Open("postgres", datasourceName)
	if err != nil {
		log.Errorf("[DB-INIT] ERROR when trying to call sql.Open(): %s", err)
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		log.Errorf("[DB-INIT] ERROR when trying to call Client.Ping() for SQL DATABASE: %s", err)
		panic(err)
	}

	//TODO replace to conifg:
	// Maximum Idle Connections
	Client.SetMaxIdleConns(5)
	// Maximum Open Connections
	Client.SetMaxOpenConns(10)
	// Idle Connection Timeout
	Client.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	Client.SetConnMaxLifetime(30 * time.Second)

	log.Infof("[DB-INIT] SUCCESS initialized and configured SQL DATABASE")
}
