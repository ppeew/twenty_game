package global

import (
	"game_web/config"
	game_proto "game_web/proto/game"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DEBUG   bool
	MongoDB *mongo.Database

	NacosConfig              = &config.NacosConfig{}
	ServerConfig             = &config.ServerConfig{}
	GameSrvClient            game_proto.GameClient
	ConsulProcessWebServices []*api.CatalogService
)
