package app

import (
	"context"
	"database/sql"
	"github.com/core-go/health"
	cas "github.com/core-go/health/cassandra"
	mgo "github.com/core-go/health/mongo"
	s "github.com/core-go/health/sql"
	_ "github.com/lib/pq"
	"go-service/video"
	"go-service/video/cassandra"
	mgs "go-service/video/mongo"
	"go-service/video/pq"
	"go-service/video/sync"
	sc "go-service/video/sync_cassandra"
	sm "go-service/video/sync_mongo"
	sq "go-service/video/sync_pq"
	"go-service/video/youtube"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-service/video/handlers_test"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	SyncHandler   *sync.SyncHandler
	ClientHandler *video.ClientHandler
	TubeHandler   *handlers_test.TubeHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	var healthHandler *health.Handler
	var clientHandler *video.ClientHandler

	var syncHandler *sync.SyncHandler

	tubeCategory := youtube.NewCategoryTubeService(root.Key)

	tubeService := youtube.NewTubeService(root.Key)
	tubeHandler := handlers_test.NewTubeHandler(tubeService)

	switch root.OpenDb {
	case 1:
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(root.Mongo.Uri))
		if err != nil {
			return nil, err
		}

		mongoDb := client.Database(root.Mongo.Database)
		mongoChecker := mgo.NewHealthChecker(mongoDb)
		healthHandler = health.NewHandler(mongoChecker)
		channelCollectionName := "channel"
		channelSyncCollectionName := "channelSync"
		playlistCollectionName := "playlist"
		playlistVideoCollectionName := "playlistVideo"
		videoCollectionName := "video"
		categoryCollectionName := "category"
		repo := sm.NewMongoVideoRepository(mongoDb, channelCollectionName, channelSyncCollectionName, playlistCollectionName, playlistVideoCollectionName, videoCollectionName, categoryCollectionName)
		syncService := sync.NewDefaultSyncService(tubeService, repo)
		syncHandler = sync.NewSyncHandler(syncService)
		clientService := mgs.NewMongoVideoService(mongoDb, channelCollectionName, channelSyncCollectionName, playlistCollectionName, playlistVideoCollectionName, videoCollectionName, categoryCollectionName, *tubeCategory)
		clientHandler = video.NewClientHandler(clientService)
		break
	case 2:
		cassDb, err := Db(&root)
		if err != nil {
			return nil, err
		}
		casChecker := cas.NewHealthChecker(cassDb)
		healthHandler = health.NewHandler(casChecker)
		repo := sc.NewCassandraVideoRepository(cassDb)
		syncService := sync.NewDefaultSyncService(tubeService, repo)
		syncHandler = sync.NewSyncHandler(syncService)
		clientService := cassandra.NewCassandraVideoService(cassDb, *tubeCategory)
		clientHandler = video.NewClientHandler(clientService)
		break
	case 3:
		postgreDB, err := sql.Open(root.Postgre.Driver, root.Postgre.DataSourceName)
		if err != nil {
			return nil, err
		}
		sqlChecker := s.NewHealthChecker(postgreDB)
		healthHandler = health.NewHandler(sqlChecker)
		repo := sq.NewSQLVideoRepository(postgreDB)
		syncService := sync.NewDefaultSyncService(tubeService, repo)
		syncHandler = sync.NewSyncHandler(syncService)
		clientService := pq.NewPostgreVideoService(postgreDB, *tubeCategory)
		clientHandler = video.NewClientHandler(clientService)
		break
	default:
		break
	}

	return &ApplicationContext{
		HealthHandler: healthHandler,
		ClientHandler: clientHandler,
		SyncHandler:   syncHandler,
		TubeHandler:   tubeHandler,
	}, nil
}
