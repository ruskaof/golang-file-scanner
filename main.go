package main

import (
	"biocadTask/internal/preprocessing"
	"biocadTask/internal/queue/rabbit"
	"biocadTask/internal/scan"
	"biocadTask/internal/storage"
	"biocadTask/internal/web"
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"log"
)

//go:embed migrations/*sql
var embedMigrations embed.FS

func main() {
	log.Println("starting application")
	databaseInfo, filesInfo, rabbitInfo, err := readConfigs("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	databaseConnection, err := storage.SetupDbConnection(
		databaseInfo.username,
		databaseInfo.password,
		databaseInfo.host,
		databaseInfo.port,
		databaseInfo.name,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = databaseConnection.Close()
	}()

	err = runMigrations(databaseConnection)
	if err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	rabbitConnection, channel, queue, err := rabbit.SetupRabbit(
		rabbitInfo.username,
		rabbitInfo.password,
		rabbitInfo.host,
		rabbitInfo.port,
		rabbitInfo.queue,
	)

	defer func() {
		_ = channel.Close()
		_ = rabbitConnection.Close()
	}()

	fileQueue := rabbit.FileMessageQueue{
		Channel: channel,
		Queue:   queue,
	}

	messageDao := storage.NewPostgresFileDao(databaseConnection)
	errorDao := storage.NewPostgresErrorDao(databaseConnection)

	filePreprocessor := preprocessing.NewFilePreprocessor(
		filesInfo.preprocessed,
		messageDao,
		errorDao,
	)

	err = filePreprocessor.CreateOutputDir()

	if err != nil {
		log.Fatalf("could not create output dir: %v", err)
	}

	go func() {
		err = fileQueue.StartConsumer(filePreprocessor.PreprocessFile)
		if err != nil {
			log.Fatalf("could not start consumer: %v", err)
		}
	}()

	preprocessedFileDao := storage.NewPostgresPreprocessedFileDao(databaseConnection)

	go scan.StartScanner(filesInfo.input, fileQueue, preprocessedFileDao)

	handler := web.NewApiHandler(messageDao)
	err = handler.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type DatabaseInfo struct {
	name     string
	port     string
	host     string
	username string
	password string
}

type FilesInfo struct {
	input        string
	preprocessed string
}

type RabbitInfo struct {
	port     string
	host     string
	username string
	password string
	queue    string
}

func readConfigs(configPath string) (DatabaseInfo, FilesInfo, RabbitInfo, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return DatabaseInfo{}, FilesInfo{}, RabbitInfo{}, err
	}

	var databaseInfo = DatabaseInfo{
		name:     viper.GetString("database.name"),
		port:     viper.GetString("database.port"),
		host:     viper.GetString("database.host"),
		username: viper.GetString("database.username"),
		password: viper.GetString("database.password"),
	}

	var filesInfo = FilesInfo{
		input:        viper.GetString("files.input"),
		preprocessed: viper.GetString("files.preprocessed"),
	}

	var rabbitInfo = RabbitInfo{
		port:     viper.GetString("rabbit.port"),
		host:     viper.GetString("rabbit.host"),
		username: viper.GetString("rabbit.username"),
		password: viper.GetString("rabbit.password"),
		queue:    viper.GetString("rabbit.queue"),
	}

	return databaseInfo, filesInfo, rabbitInfo, nil
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
