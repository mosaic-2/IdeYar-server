package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mosaic-2/IdeYar-server/internal/sql/dbpkg"
	"github.com/mosaic-2/IdeYar-server/pkg/authpb"
)

type Server struct {
	authpb.UnimplementedAuthServer
	conn       *sql.DB
	query      *dbpkg.Queries
	hmacSecret []byte
}

func getQuery() (*dbpkg.Queries, *sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}
	q := dbpkg.New(conn)
	return q, conn, nil
}

func NewServer() (*Server, error) {
	q, conn, err := getQuery()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &Server{
		conn:       conn,
		query:      q,
		hmacSecret: []byte(os.Getenv("SECRET_KEY")),
	}, nil
}
