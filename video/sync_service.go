package video

import (
	"context"
)

type SyncService interface {
	SyncChannel(ctx context.Context, channelId string) (int, error)
	SyncChannels(ctx context.Context, channelIds []string) (int, error)
	SyncPlaylist(ctx context.Context, playlistId string, level *int) (int, error)
	//   syncPlaylists(playlistIds string[], level? number) (int,error);
	GetSubscriptions(ctx context.Context, channelId string) ([]Channel, error)
}
