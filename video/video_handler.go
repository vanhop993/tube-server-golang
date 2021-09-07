package video

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ClientHandler struct {
	Client VideoService
}

func NewClientHandler(clientService VideoService) *ClientHandler {
	return &ClientHandler{Client: clientService}
}

func (c *ClientHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	s := strings.Split(params, "&")
	if len(s[0]) <= 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	var fields []string
	if len(s) > 1 && len(s[1]) > 0 {
		fields = strings.Split(s[1], ",")
	}
	res, err := c.Client.GetChannel(r.Context(), s[0], fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	s := strings.Split(params, "&")
	if len(s[0]) <= 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	arrayId := strings.Split(s[0], ",")
	var fields []string
	if len(s) > 1 && len(s[1]) > 0 {
		fields = strings.Split(s[1], ",")
	}
	res, err := c.Client.GetChannels(r.Context(), arrayId, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	s := strings.Split(params, "&")
	if len(s[0]) <= 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	var fields []string
	if len(s) > 1 && len(s[1]) > 0 {
		fields = strings.Split(s[1], ",")
	}
	res, err := c.Client.GetPlaylist(r.Context(), s[0], fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	s := strings.Split(params, "&")
	if len(s[0]) <= 0 {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}
	ids := strings.Split(s[0], ",")
	var fields []string
	if len(s) > 1 && len(s[1]) > 0 {
		fields = strings.Split(s[1], ",")
	}
	res, err := c.Client.GetPlaylists(r.Context(), ids, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	respond(w, res)
}

func (c *ClientHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	s := strings.Split(params, "&")
	if len(s[0]) <= 0 {
		http.Error(w, "Id cannot empty!", http.StatusBadRequest)
		return
	}
	var fields []string
	if len(s) > 1 && len(s[1]) > 0 {
		fields = strings.Split(s[1], ",")
	}
	res, err := c.Client.GetVideo(r.Context(), s[0], fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetVideos(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	s := strings.Split(params, "&")
	if len(s[0]) <= 0 {
		http.Error(w, "Ids cannot be empty!", http.StatusBadRequest)
		return
	}
	ids := strings.Split(s[0], ",")
	var fields []string
	if len(s) > 1 && len(s[1]) > 0 {
		fields = strings.Split(s[1], ",")
	}
	res, err := c.Client.GetVideos(r.Context(), ids, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetChannelPlaylists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	channelId := query.Get("channelId")
	if len(channelId) <= 0 {
		http.Error(w, "ChannelId cannot be empty", http.StatusBadRequest)
		return
	}
	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	res, er1 := c.Client.GetChannelPlaylists(r.Context(), channelId, limit, nextPageToken, fields)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetVideosFromChannelIdOrPlaylistId(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	playlistId := query.Get("playlistId")
	if len(playlistId) > 0 {
		limitString := query.Get("limit")
		var limit int
		if len(limitString) > 0 {
			res, err := strconv.Atoi(limitString)
			if err != nil {
				http.Error(w, "Limit is not number", http.StatusBadRequest)
			}
			limit = res
		} else {
			limit = 10
		}
		nextPageToken := query.Get("nextPageToken")
		if len(nextPageToken) <= 0 {
			nextPageToken = ""
		}
		var fields []string
		fieldsString := query.Get("fields")
		if len(fieldsString) > 0 {
			fields = strings.Split(fieldsString, ",")
		}
		res, er1 := c.Client.GetPlaylistVideos(r.Context(), playlistId, limit, nextPageToken, fields)
		if er1 != nil {
			http.Error(w, er1.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, res)
	} else {
		channelId := query.Get("channelId")
		if len(channelId) <= 0 {
			http.Error(w, "Require channelId or playlistId", http.StatusBadRequest)
			return
		}
		limitString := query.Get("limit")
		var limit int
		if len(limitString) > 0 {
			res, err := strconv.Atoi(limitString)
			if err != nil {
				http.Error(w, "Limit is not number", http.StatusBadRequest)
			}
			limit = res
		} else {
			limit = 10
		}
		nextPageToken := query.Get("nextPageToken")
		if len(nextPageToken) <= 0 {
			nextPageToken = ""
		}
		var fields []string
		fieldsString := query.Get("fields")
		if len(fieldsString) > 0 {
			fields = strings.Split(fieldsString, ",")
		}
		log.Println(channelId, limit, nextPageToken, fields)
		res, er1 := c.Client.GetChannelVideos(r.Context(), channelId, limit, nextPageToken, fields)
		if er1 != nil {
			http.Error(w, er1.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, res)
	}
}

func (c *ClientHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["params"]
	res, err := c.Client.GetCagetories(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) SearchChannel(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	var channelSM ChannelSM
	channelSM.Q = strings.TrimSpace(query.Get("q"))
	channelSM.ChannelId = strings.TrimSpace(query.Get("channelId"))
	channelSM.Sort = strings.TrimSpace(query.Get("sort"))
	layout := "2006-01-02T15:04:05Z"
	if query.Get("publishedAfter") != "" {
		t, err := time.Parse(layout, query.Get("publishedAfter"))
		if err != nil {
			http.Error(w, "publishedAfter is not time", http.StatusBadRequest)
			return
		}
		channelSM.PublishedAfter = &t
	}

	if query.Get("publishedAfter") != "" {
		t1, err := time.Parse(layout, query.Get("publishedBefore"))
		if err != nil {
			http.Error(w, "publishedBefore is not time", http.StatusBadRequest)
			return
		}
		channelSM.PublishedBefore = &t1
	}

	channelSM.RegionCode = query.Get("regionCode")
	log.Println(channelSM)
	res, er1 := c.Client.SearchChannel(r.Context(), channelSM, limit, nextPageToken, fields)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) SearchPlaylists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	var playlistSM PlaylistSM
	playlistSM.Q = strings.TrimSpace(query.Get("q"))
	playlistSM.ChannelId = strings.TrimSpace(query.Get("channelId"))
	playlistSM.Sort = strings.TrimSpace(query.Get("sort"))
	layout := "2006-01-02T15:04:05Z"
	if query.Get("publishedAfter") != "" {
		t, err := time.Parse(layout, query.Get("publishedAfter"))
		if err != nil {
			http.Error(w, "publishedAfter is not time", http.StatusBadRequest)
			return
		}
		playlistSM.PublishedAfter = &t
	}

	if query.Get("publishedAfter") != "" {
		t1, err := time.Parse(layout, query.Get("publishedBefore"))
		if err != nil {
			http.Error(w, "publishedBefore is not time", http.StatusBadRequest)
			return
		}
		playlistSM.PublishedBefore = &t1
	}

	playlistSM.RegionCode = query.Get("regionCode")
	log.Println(playlistSM)
	res, er1 := c.Client.SearchPlaylists(r.Context(), playlistSM, limit, nextPageToken, fields)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) SearchVideos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	var itemSM ItemSM
	itemSM.Q = strings.TrimSpace(query.Get("q"))
	itemSM.ChannelId = strings.TrimSpace(query.Get("channelId"))
	itemSM.Sort = strings.TrimSpace(query.Get("sort"))
	layout := "2006-01-02T15:04:05Z"
	if query.Get("publishedAfter") != "" {
		t, err := time.Parse(layout, query.Get("publishedAfter"))
		if err != nil {
			http.Error(w, "publishedAfter is not time", http.StatusBadRequest)
			return
		}
		itemSM.PublishedAfter = &t
	}

	if query.Get("publishedAfter") != "" {
		t1, err := time.Parse(layout, query.Get("publishedBefore"))
		if err != nil {
			http.Error(w, "publishedBefore is not time", http.StatusBadRequest)
			return
		}
		itemSM.PublishedBefore = &t1
	}
	if query.Get("duration") != "" {
		itemSM.Duration = query.Get("duration")
	}
	itemSM.RegionCode = query.Get("regionCode")
	log.Println(itemSM)
	res, er1 := c.Client.SearchVideos(r.Context(), itemSM, limit, nextPageToken, fields)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	var itemSM ItemSM
	itemSM.Q = strings.TrimSpace(query.Get("q"))
	itemSM.ChannelId = strings.TrimSpace(query.Get("channelId"))
	itemSM.Sort = strings.TrimSpace(query.Get("sort"))
	layout := "2006-01-02T15:04:05Z"
	if query.Get("publishedAfter") != "" {
		t, err := time.Parse(layout, query.Get("publishedAfter"))
		if err != nil {
			http.Error(w, "publishedAfter is not time", http.StatusBadRequest)
			return
		}
		itemSM.PublishedAfter = &t
	}

	if query.Get("publishedAfter") != "" {
		t1, err := time.Parse(layout, query.Get("publishedBefore"))
		if err != nil {
			http.Error(w, "publishedBefore is not time", http.StatusBadRequest)
			return
		}
		itemSM.PublishedBefore = &t1
	}
	if query.Get("duration") != "" {
		itemSM.Duration = query.Get("duration")
	}
	itemSM.RegionCode = query.Get("regionCode")
	log.Println(itemSM)
	res, er1 := c.Client.SearchVideos(r.Context(), itemSM, limit, nextPageToken, fields)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetRelatedVideos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	if len(id) <= 0 {
		http.Error(w, "id can not empty", http.StatusBadRequest)
	}
	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	res, err := c.Client.GetRelatedVideos(r.Context(), id, limit, nextPageToken, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func (c *ClientHandler) GetPopularVideos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	categoryId := query.Get("categoryId")
	regionCode := query.Get("regionCode")

	limitString := query.Get("limit")
	var limit int
	if len(limitString) > 0 {
		res, err := strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, "Limit is not number", http.StatusBadRequest)
		}
		limit = res
	} else {
		limit = 10
	}
	nextPageToken := query.Get("nextPageToken")
	if len(nextPageToken) <= 0 {
		nextPageToken = ""
	}
	var fields []string
	fieldsString := query.Get("fields")
	if len(fieldsString) > 0 {
		fields = strings.Split(fieldsString, ",")
	}
	res, err := c.Client.GetPopularVideos(r.Context(), regionCode, categoryId, limit, nextPageToken, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, res)
}

func respond(w http.ResponseWriter, result interface{}) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
