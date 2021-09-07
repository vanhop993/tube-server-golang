package app

import (
	"github.com/gocql/gocql"
	"log"
	"time"
)

const (
	Keyspace = `tube`

	CreateKeyspace = `create keyspace if not exists tube with replication = {'class':'SimpleStrategy', 'replication_factor':1}`

	CreateChannelTable = `
					CREATE TABLE IF NOT EXISTS tube.channel (
							id varchar,
							count float,
							country varchar,
							customUrl varchar,
							description varchar,
							favorites varchar,
							highThumbnail varchar,
							itemCount float,
							likes varchar,
							localizedDescription varchar,
							localizedTitle varchar,
							mediumThumbnail varchar,
							playlistCount float,
							playlistItemCount float,
							playlistVideoCount float,
							playlistVideoItemCount float,
							publishedAt timestamp,
							thumbnail varchar,
							lastUpload timestamp,
							title varchar,
							uploads varchar, 
							PRIMARY KEY(id )
					)`
	CreateChannelSyncTable = `
					CREATE TABLE IF NOT EXISTS tube.channelSync (
							id varchar,
							synctime timestamp,
							uploads varchar, 
							PRIMARY KEY(id ) 
					)`
	CreatePlaylistTable = `
					CREATE TABLE IF NOT EXISTS tube.playlist (
							id varchar,
							channelId varchar,
							channelTitle varchar,
							count float,
							itemCount float,
							description varchar,
							highThumbnail varchar,
							localizedDescription varchar,
							localizedTitle varchar,
							maxresThumbnail varchar,
							mediumThumbnail varchar,
							publishedAt timestamp,
							standardThumbnail varchar,
							thumbnail varchar,
							title varchar, 
							PRIMARY KEY(id )
					)`
	CreatePlaylistVideoTable = `
					CREATE TABLE IF NOT EXISTS tube.playlistvideo (
							id varchar,
							videos list<varchar>, 
							PRIMARY KEY(id )
					)`
	CreateVideoTable = `
					CREATE TABLE IF NOT EXISTS tube.video (
							id varchar,
							caption varchar,
							categoryId varchar,
							channelId varchar,
							channelTitle varchar,
							defaultAudioLanguage varchar,
							defaultLanguage varchar,
							definition float,
							description varchar,
							dimension varchar,
							duration float,
							highThumbnail varchar,
							licensedContent boolean,
							liveBroadcastContent varchar,
							localizedDescription varchar,
							localizedTitle varchar,
							maxresThumbnail varchar,
							mediumThumbnail varchar,
							projection varchar,
							publishedAt timestamp,
							standardThumbnail varchar,
							tags list<varchar>,
							thumbnail varchar,
							title varchar,
							blockedRegions list<varchar>,
							allowedRegions list<varchar>, 
							PRIMARY KEY(id) 
					)`
	CreaateCategoryType = `CREATE TYPE IF NOT EXISTS tube.categorisType (id varchar,title varchar,assignable boolean,channelId varchar)`
	CreateCategoryTable = `CREATE TABLE IF NOT EXISTS tube.category (
							id varchar,
							data list<frozen<categorisType>>, 
							PRIMARY KEY(id ))`
	CreateVideoLuceneIndex    = `CREATE CUSTOM INDEX IF NOT EXISTS video_index ON tube.video (title) USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {'refresh_seconds': '1','schema': '{fields: {"id":{"type":"text"},"caption":{"type":"boolean"},"categoryid":{"type":"text"},"channelid":{"type":"text"},"channeltitle":{"type":"text"},"defaultaudiolanguage":{"type":"text"},"defaultlanguage":{"type":"text"},"definition":{"type":"float"},"description":{"type":"text"},"dimension":{"type":"text"},"duration":{"type":"float"},"highthumbnail":{"type":"float"},"licensedcontent":{"type":"boolean"},"livebroadcastcontent":{"type":"text"},"localizeddescription":{"type":"text"},"localizedtitle":{"type":"text"},"maxresthumbnail":{"type":"text"},"mediumthumbnail":{"type":"text"},"projection":{"type":"text"},"publishedat":{"type":"date","pattern":"yyyy-MM-dd HH:mm:ss"},"standardthumbnail":{"type":"text"},"blockedregions":{"type":"string"},"tags":{"type":"string"},"thumbnail":{"type":"text"},"title":{"type":"string"}}}'}`
	CreateChannelLuceneIndex  = `CREATE CUSTOM INDEX IF NOT EXISTS channel_index  ON tube.channel (title) USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {'refresh_seconds': '1','schema': '{fields: {"id":{"type":"text"},"count":{"type":"float"},"country":{"type":"text"},"customurl":{"type":"text"},"description":{"type":"text"},"favorites":{"type":"text"},"highthumbnail":{"type":"text"},"itemcount":{"type":"float"},"likes":{"type":"text"},"localizeddescription":{"type":"text"},"localizedtitle":{"type":"text"},"mediumthumbnail":{"type":"text"},"playlistcount":{"type":"float"},"playlistitemcount":{"type":"float"},"playlistvideocount":{"type":"float"},"playlistvideoitemcount":{"type":"float"},"publishedat":{"type":"date","pattern":"yyyy-MM-dd HH:mm:ss"},"thumbnail":{"type":"text"},"lastupload":{"type":"date","pattern":"yyyy-MM-dd HH:mm:ss"},"title":{"type":"text"},"uploads":{"type":"text"}}}'}`
	CreatePlaylistLuceneIndex = `CREATE CUSTOM INDEX IF NOT EXISTS playlist_index ON tube.playlist (title) USING 'com.stratio.cassandra.lucene.Index' WITH OPTIONS = {'refresh_seconds': '1','schema': '{fields: {"id":{"type":"text"},"channelid":{"type":"text"},"channeltitle":{"type":"text"},"count":{"type":"float"},"itemcount":{"type":"float"},"description":{"type":"text"},"highthumbnail":{"type":"text"},"localizeddescription":{"type":"text"},"localizedtitle":{"type":"text"},"maxresthumbnail":{"type":"text"},"mediumthumbnail":{"type":"text"},"publishedat":{"type":"date","pattern":"yyyy-MM-dd HH:mm:ss"},"standardthumbnail":{"type":"text"},"thumbnail":{"type":"text"},"title":{"type":"text"}}}'};`
)

func Db(root *Root) (*gocql.ClusterConfig, error) {
	cluster := gocql.NewCluster(root.Cassandra.Uri)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	log.Println(root.Cassandra.Username, root.Cassandra.Password)
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: root.Cassandra.Username, Password: root.Cassandra.Password}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	err = session.Query(CreateKeyspace).Exec()
	if err != nil {
		return nil, err
	}

	// create table
	cluster.Keyspace = Keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return cluster, nil
}
