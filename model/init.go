package model

import (
	"context"
	"os"

	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Database struct {
	Self *mongo.Client
}

var DB *Database

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// used for cli
func InitSelfDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(viper.GetString("db.url"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Errorf(err, "Database connection failed.")
		os.Exit(-1)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Errorf(err, "Database check connection failed.")
		os.Exit(-1)
	}

	log.Info("Connected to MongoDB!")

	return client
}

func GetSelfDB() *mongo.Client {
	return InitSelfDB()
}

func (db *Database) Init() {
	DB = &Database{
		Self: GetSelfDB(),
	}
}

func (db *Database) Close() {
	_ = DB.Self.Disconnect(context.TODO())
}
