package main

import (
	"github.com/rs/zerolog/log"
	"hson98/app-chat/config"
	"hson98/app-chat/database"
	"hson98/app-chat/internal/server"
	customvalidate "hson98/app-chat/pkg/customValidate"
	"hson98/app-chat/pkg/myjwt"
)

func main() {
	config, err := config.LoadConfig(".", "app")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	psqlDB := database.NewPostgresDB(config)

	redisClient := database.GetRedisClient(config.RedisHost, config.RedisPass)
	defer redisClient.Close()
	customvalidate.RegisterValidate()
	//khởi tạo JWT
	jwtMaker, err := myjwt.NewJwtMaker(config.SecretKeyJWT)
	if err != nil {
		panic(err)
	}
	s := server.NewServer(psqlDB, redisClient, jwtMaker, config)
	if err := s.Run(); err != nil {
		log.Fatal().Err(err).Msg("cannot run server")
	}

}
