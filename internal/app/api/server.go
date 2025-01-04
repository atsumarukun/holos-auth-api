package api

import (
	"context"
	"holos-auth-api/internal/pkg/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Serve() {
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	c := &mysql.Config{
		Addr:      config.MySQLHost + ":" + config.MySQLPort,
		User:      config.MySQLUser,
		Passwd:    config.MySQLPassword,
		DBName:    config.MySQLDatabase,
		Net:       "tcp",
		ParseTime: true,
	}
	db, err := sqlx.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	inject(db)

	r := gin.Default()
	registerRouter(r)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	<-ctx.Done()

	ctx, stop = context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln(err.Error())
	}
}
