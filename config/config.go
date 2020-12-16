package config

import (
	"encoding/json"
	"github.com/secure-for-ai/secureai-microsvs/cache"
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/secure-for-ai/secureai-microsvs/session"
	"github.com/secure-for-ai/secureai-microsvs/snowflake"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	MongoDB   db.MongoDBConf
	Redis     cache.RedisConf
	Session   session.HybridStoreConf
	Snowflake snowflake.NodeConf
	AppInfo   util.AppInfo
}

var Conf *Config
var MongoDBClient *db.MongoDBClient
var RedisClient *cache.RedisClient
var SessionStore *session.HybridStore
var SnowflakeNode *snowflake.Node

var confPathPrefix = defaultConfPath("config")

func init() {
	log.Println("Begin init")

	initConf()
	log.Println(Conf)

	intiMongoDB()
	initRedis()
	initSession()
	initSnowflake()

	log.Println("Over init")
	log.Println("Env: ", Conf.AppInfo.Env)
}

func defaultConfPath(dir string) string {
	wdPath, err := os.Getwd()

	if err != nil {
		log.Panic(err)
		return ""
	}

	s := path.Join(wdPath, dir)
	return s
}

func initConf() {
	log.Println("Begin init default config")

	Conf = &Config{}
	fileName := "default.json"

	if v, ok := os.LookupEnv("CONFIG_PATH"); ok {
		confPathPrefix = v
	}

	// read default config
	defaultConfFilePath := path.Join(confPathPrefix, fileName)
	data, err := ioutil.ReadFile(defaultConfFilePath)

	if err != nil {
		log.Println("config-initConf: read default.json error")
		log.Panic(err)
		return
	}

	err = json.Unmarshal(data, Conf)
	if err != nil {
		log.Println("config-initConf: unmarshal default.json error")
		log.Panic(err)
		return
	}

	// read env and config path
	if v, ok := os.LookupEnv("ENV"); ok {
		fileName = v + ".json"
	}

	if fileName != "default.json" {
		// read env config
		data, err = ioutil.ReadFile(path.Join(confPathPrefix, fileName))
		if err != nil {
			log.Println("config-initConf: read [env].json error")
			log.Panic(err)
			return
		}

		err = json.Unmarshal(data, Conf)
		if err != nil {
			log.Println("config-initConf: unmarshal [env].json error")
			log.Panic(err)
			return
		}
	}

	log.Println("Over init default config")
}

func intiMongoDB() {
	log.Println("Begin init mongoDB")

	client, err := db.NewMongoDB(Conf.MongoDB)

	if err != nil {
		log.Println("unable to init mongoDB")
		log.Panic(err)
		return
	}

	MongoDBClient = client

	log.Println("Over init mongoDB")
}

func initRedis() {
	log.Println("Begin init redis")

	var err error
	RedisClient, err = cache.NewRedisClient(Conf.Redis)
	if err != nil {
		log.Println("unable to init redis")
		log.Panic(err)
		return
	}

	log.Println("Over init redis")
}

func initSession() {
	log.Println("Begin init session store")

	storage := session.RedisStoreEngine{
		RedisClient: RedisClient,
		Serializer:  session.GobSerializer{},
		Prefix:      "session_",
		IDGenerator: session.SUIDInt64Generator{},
	}
	SessionStore = session.NewSessionStore(&storage, &Conf.Session)

	log.Println("Over init session store")
}

func initSnowflake() {
	log.Println("Begin init snowflake node")

	var err error
	// Todo Add service to get snowflake id
	SnowflakeNode, err = snowflake.NewNode(0, &Conf.Snowflake)

	if err != nil {
		log.Println("unable to init snowflake node")
		log.Panic(err)
	}

	log.Println("Over init snowflake node")
}
