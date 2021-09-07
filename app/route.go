package app

import (
	"context"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func Route(r *mux.Router, context context.Context, root Root) error {
	app, err := NewApp(context, root)
	if err != nil {
		return err
	}
	r.HandleFunc("/health", app.HealthHandler.Check).Methods(GET)

	r.HandleFunc("/tube/channel", app.SyncHandler.SyncChannel).Methods(POST)
	r.HandleFunc("/tube/channels", app.SyncHandler.SyncChannels).Methods(POST)
	r.HandleFunc("/tube/playlist", app.SyncHandler.SyncPlaylist).Methods(POST)
	r.HandleFunc("/tube/channels/subscriptions/{id}", app.SyncHandler.SyncSubctiption).Methods(GET)

	r.HandleFunc("/tube/channel/{params}", app.ClientHandler.GetChannel).Methods(GET)
	r.HandleFunc("/tube/channels/list/{params}", app.ClientHandler.GetChannels).Methods(GET)
	r.HandleFunc("/tube/playlist/{params}", app.ClientHandler.GetPlaylist).Methods(GET)
	r.HandleFunc("/tube/playlists/list/{params}", app.ClientHandler.GetPlaylists).Methods(GET)
	r.HandleFunc("/tube/video/{params}", app.ClientHandler.GetVideo).Methods(GET)
	r.HandleFunc("/tube/videos/list/{params}", app.ClientHandler.GetVideos).Methods(GET)
	r.HandleFunc("/tube/playlists", app.ClientHandler.GetChannelPlaylists).Methods(GET)
	r.HandleFunc("/tube/videos", app.ClientHandler.GetVideosFromChannelIdOrPlaylistId).Methods(GET)
	r.HandleFunc("/tube/category/{params}", app.ClientHandler.GetCategory).Methods(GET)
	r.HandleFunc("/tube/channels/search", app.ClientHandler.SearchChannel).Methods(GET)
	r.HandleFunc("/tube/playlists/search", app.ClientHandler.SearchPlaylists).Methods(GET)
	r.HandleFunc("/tube/videos/search", app.ClientHandler.SearchVideos).Methods(GET)
	r.HandleFunc("/tube/videos/related", app.ClientHandler.GetRelatedVideos).Methods(GET)
	r.HandleFunc("/tube/videos/popular", app.ClientHandler.GetPopularVideos).Methods(GET)

	r.HandleFunc("/channel/{id}", app.TubeHandler.GetChannel).Methods(GET)
	r.HandleFunc("/channels/{id}", app.TubeHandler.GetChannels).Methods(GET)
	r.HandleFunc("/playlist/{id}", app.TubeHandler.GetPlaylist).Methods(GET)
	r.HandleFunc("/playlists/{id}", app.TubeHandler.GetPlaylists).Methods(GET)
	r.HandleFunc("/channelplaylists/{id}", app.TubeHandler.GetChannelPlaylists).Methods(GET)
	r.HandleFunc("/playlistvideos/{id}", app.TubeHandler.GetPlaylistVideos).Methods(GET)
	r.HandleFunc("/videos/{id}", app.TubeHandler.GetVideos).Methods(GET)

	return err
}
