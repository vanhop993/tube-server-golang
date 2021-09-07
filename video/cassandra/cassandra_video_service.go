package cassandra

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"go-service/video"
	"go-service/video/common"
	"go-service/video/youtube"
	"log"
	"strings"
)

type CassandraVideoService struct {
	cass         *gocql.ClusterConfig
	tubeCategory youtube.CategoryTubeClient
}

func NewCassandraVideoService(cass *gocql.ClusterConfig, tubeCategory youtube.CategoryTubeClient) *CassandraVideoService {
	return &CassandraVideoService{cass: cass, tubeCategory: tubeCategory}
}

func (c *CassandraVideoService) GetChannel(ctx context.Context, channelId string, fields []string) (*video.Channel, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	query := fmt.Sprintf(`Select %s from channel where id = ?`, strings.Join(fields, ","))
	q := session.Query(query, channelId)
	if q.Exec() != nil {
		return nil, q.Exec()
	}
	q.Iter()
	res, err := channelConvert(q.Iter())
	if err != nil {
		return nil, err
	}
	if len(res[0].ChannelList) > 0 {
		channels, err := c.GetChannels(ctx, res[0].ChannelList, []string{})
		if err != nil {
			return nil, err
		}
		res[0].Channels = *channels
	}
	return &res[0], nil
}

func (c *CassandraVideoService) GetChannels(ctx context.Context, ids []string, fields []string) (*[]video.Channel, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = "?"
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	query := fmt.Sprintf(`Select %s from channel where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	q := session.Query(query, cc...)
	if q.Exec() != nil {
		return nil, q.Exec()
	}
	res, err := channelConvert(q.Iter())
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *CassandraVideoService) GetPlaylist(ctx context.Context, id string, fields []string) (*video.Playlist, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	query := fmt.Sprintf(`Select %s from playlist where id = ?`, strings.Join(fields, ","))
	rows := session.Query(query, id)
	if rows.Exec() != nil {
		return nil, rows.Exec()
	}
	res, err := playlistConvert(rows.Iter())
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (c *CassandraVideoService) GetPlaylists(ctx context.Context, ids []string, fields []string) (*[]video.Playlist, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = "?"
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	query := fmt.Sprintf(`Select %s from playlist where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	rows := session.Query(query, cc...)
	if rows.Exec() != nil {
		return nil, rows.Exec()
	}
	result, err := playlistConvert(rows.Iter())
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *CassandraVideoService) GetVideo(ctx context.Context, id string, fields []string) (*video.Video, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	query := fmt.Sprintf(`Select %s from video where id = ?`, strings.Join(fields, ","))
	rows := session.Query(query, id)
	if rows.Exec() != nil {
		return nil, rows.Exec()
	}
	res, err := videoConvert(rows.Iter())
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (c *CassandraVideoService) GetVideos(ctx context.Context, ids []string, fields []string) (*[]video.Video, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = "?"
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	query := fmt.Sprintf(`Select %s from video where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	rows := session.Query(query, cc...)
	if rows.Exec() != nil {
		return nil, rows.Exec()
	}
	res, err := videoConvert(rows.Iter())
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *CassandraVideoService) GetChannelPlaylists(ctx context.Context, channelId string, max int, nextPageToken string, fields []string) (*video.ListResultPlaylist, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sort := map[string]interface{}{"field": `publishedat`, "reverse": true}
	must := map[string]interface{}{"type": "match", "field": "channelid", "value": fmt.Sprintf(`%s`, channelId)}
	a := map[string]interface{}{
		"filter": map[string]interface{}{
			"must": must,
		},
		"sort": sort,
	}
	queryObj, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf(`select %s from playlist where expr(playlist_index, '%s')`, strings.Join(fields, ","), queryObj)
	var query *gocql.Query
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	query = session.Query(sql).PageState(next).PageSize(max)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	iter := query.Iter()
	var res video.ListResultPlaylist
	res.NextPageToken = hex.EncodeToString(iter.PageState())
	res.List, err = playlistConvert(iter)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *CassandraVideoService) GetChannelVideos(ctx context.Context, channelId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sort := map[string]interface{}{"field": `publishedat`, "reverse": true}
	must := map[string]interface{}{"type": "match", "field": "channelid", "value": fmt.Sprintf(`%s`, channelId)}
	a := map[string]interface{}{
		"filter": map[string]interface{}{
			"must": must,
		},
		"sort": sort,
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	queryObj, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf(`select %s from video where expr(video_index, '%s')`, strings.Join(fields, ","), queryObj)
	var query *gocql.Query
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	query = session.Query(sql).PageState(next).PageSize(max)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	iter := query.Iter()
	var res video.ListResultVideos
	res.NextPageToken = hex.EncodeToString(iter.PageState())
	res.List, err = videoConvert(iter)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *CassandraVideoService) GetPlaylistVideos(ctx context.Context, playlistId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	const sql = `select videos from playlistVideo where id = ? `
	query := session.Query(sql, playlistId)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	var ids []string
	query.Iter().Scan(&ids)
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = "?"
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	queryV := fmt.Sprintf(`Select %s from video where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	rows := session.Query(queryV, cc...).PageState(next).PageSize(max)
	if rows.Exec() != nil {
		return nil, rows.Exec()
	}
	var res video.ListResultVideos
	res.List, err = videoConvert(rows.Iter())
	if err != nil {
		return nil, err
	}
	res.NextPageToken = hex.EncodeToString(rows.Iter().PageState())
	return &res, nil
}

func (c *CassandraVideoService) GetCagetories(ctx context.Context, regionCode string) (*video.Categories, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sql := `select * from category where id = ?`
	query := session.Query(sql, regionCode)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	var category video.Categories
	query.Iter().Scan(&category.Id, &category.Data)
	if category.Data == nil {
		res, er1 := c.tubeCategory.GetCagetories(regionCode)
		if er1 != nil {
			return nil, er1
		}
		query := "insert into category (id,data) values (?, ?)"
		err := session.Query(query, regionCode, res).Exec()
		if err != nil {
			return nil, err
		}
		result := video.Categories{
			Id:   regionCode,
			Data: *res,
		}
		return &result, nil
	}
	return &category, nil
}

func (c *CassandraVideoService) SearchChannel(ctx context.Context, channelSM video.ChannelSM, max int, nextPageToken string, fields []string) (*video.ListResultChannel, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sql, er1 := buildChannelSearch(channelSM, fields)
	if er1 != nil {
		return nil, er1
	}
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	query := session.Query(sql).PageState(next).PageSize(max)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	var res video.ListResultChannel
	res.List, err = channelConvert(query.Iter())
	if err != nil {
		return nil, err
	}
	res.NextPageToken = hex.EncodeToString(query.Iter().PageState())
	return &res, nil
}

func (c *CassandraVideoService) SearchPlaylists(ctx context.Context, playlistSM video.PlaylistSM, max int, nextPageToken string, fields []string) (*video.ListResultPlaylist, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sql, er1 := buildPlaylistSearch(playlistSM, fields)
	if er1 != nil {
		return nil, er1
	}
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	query := session.Query(sql).PageState(next).PageSize(max)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	var res video.ListResultPlaylist
	res.List, err = playlistConvert(query.Iter())
	if err != nil {
		return nil, err
	}
	res.NextPageToken = hex.EncodeToString(query.Iter().PageState())
	return &res, nil
}

func (c *CassandraVideoService) SearchVideos(ctx context.Context, itemSM video.ItemSM, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sql, er1 := buildVideosSearch(itemSM, fields)
	if er1 != nil {
		return nil, er1
	}
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	query := session.Query(sql).PageState(next).PageSize(max)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	var res video.ListResultVideos
	res.List, err = videoConvert(query.Iter())
	if err != nil {
		return nil, err
	}
	res.NextPageToken = hex.EncodeToString(query.Iter().PageState())
	return &res, nil
}

func (c *CassandraVideoService) Search(ctx context.Context, itemSM video.ItemSM, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	sql, er1 := buildVideosSearch(itemSM, fields)
	if er1 != nil {
		return nil, er1
	}
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	query := session.Query(sql).PageState(next).PageSize(max)
	if query.Exec() != nil {
		return nil, query.Exec()
	}
	var res video.ListResultVideos
	res.List, err = videoConvert(query.Iter())
	if err != nil {
		return nil, err
	}
	res.NextPageToken = hex.EncodeToString(query.Iter().PageState())
	return &res, nil
}

func (c *CassandraVideoService) GetRelatedVideos(ctx context.Context, videoId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	var a []string
	resVd, er1 := c.GetVideo(ctx, videoId, a)
	if er1 != nil {
		return nil, er1
	}
	if resVd == nil {
		return nil, errors.New("video don't exist")
	} else {
		var should []interface{}
		for _, v := range resVd.Tags {
			should = append(should, map[string]interface{}{"type": "contains", "field": "tags", "values": v})
		}
		not := map[string]interface{}{"type": "match", "field": "id", "value": videoId}
		sort := map[string]interface{}{"field": "publishedat", "reverse": true}
		fields = checkFields("publishedAt", fields)
		a := map[string]interface{}{
			"filter": map[string]interface{}{
				"should": should,
				"not":    not,
			},
			"sort": sort,
		}
		queryObj, err := json.Marshal(a)
		if err != nil {
			return nil, err
		}
		if len(fields) <= 0 {
			fields = append(fields, "*")
		}
		sql := fmt.Sprintf(`select %s from video where expr(video_index,'%s')`, strings.Join(fields, ","), queryObj)
		next, err := hex.DecodeString(nextPageToken)
		if err != nil {
			return nil, err
		}
		query := session.Query(sql).PageState(next).PageSize(max)
		var res video.ListResultVideos
		res.List, err = videoConvert(query.Iter())
		if err != nil {
			return nil, err
		}
		res.NextPageToken = hex.EncodeToString(query.Iter().PageState())
		return &res, nil
	}
}

func (c *CassandraVideoService) GetPopularVideos(ctx context.Context, regionCode string, categoryId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	session, er0 := c.cass.CreateSession()
	if er0 != nil {
		return nil, er0
	}
	query := []interface{}{}
	not := []interface{}{}
	if len(regionCode) > 0 {
		not = append(not, map[string]interface{}{"type": "contains", "field": "blockedregions", "values": regionCode})
	}
	if len(categoryId) > 0 {
		query = append(query, map[string]interface{}{"type": "match", "field": "categoryid", "value": categoryId})
		fields = checkFields("categoryId", fields)
	}
	sort := map[string]interface{}{"field": "publishedat", "reverse": true}
	fields = checkFields("publishedAt", fields)
	a := map[string]interface{}{
		"filter": map[string]interface{}{
			"not": not,
		},
		"query": query,
		"sort":  sort,
	}
	if len(not) == 0 {
		delete(a, "filter")
	}
	if len(query) == 0 {
		delete(a, "query")
	}
	queryObj, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	sql := fmt.Sprintf(`select %s from video where expr(video_index,'%s')`, strings.Join(fields, ","), queryObj)
	next, err := hex.DecodeString(nextPageToken)
	if err != nil {
		return nil, err
	}
	q := session.Query(sql).PageState(next).PageSize(max)
	var res video.ListResultVideos
	res.List, err = videoConvert(q.Iter())
	if err != nil {
		return nil, err
	}
	res.NextPageToken = hex.EncodeToString(q.Iter().PageState())
	return &res, nil
}

func buildChannelSearch(s video.ChannelSM, fields []string) (string, error) {
	should := []interface{}{}
	must := []interface{}{}
	not := []interface{}{}
	sort := []interface{}{}
	if len(s.Q) > 0 {
		should = append(should, map[string]interface{}{"type": "phrase", "field": "title", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "prefix", "field": "title", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "title", "value": fmt.Sprintf(`*%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "title", "value": fmt.Sprintf(`%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "phrase", "field": "description", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "prefix", "field": "description", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s*`, s.Q)})
	}
	if s.PublishedBefore != nil && s.PublishedAfter != nil {
		t1 := s.PublishedBefore.Format("2006-01-02 15:04:05")
		t2 := s.PublishedAfter.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "lower": t1, "upper": t2})
		fields = checkFields("publishedAt", fields)
	} else if s.PublishedAfter != nil {
		t2 := s.PublishedAfter.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "upper": t2})
		fields = checkFields("publishedAt", fields)
	} else if s.PublishedBefore != nil {
		t1 := s.PublishedBefore.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "lower": t1})
		fields = checkFields("publishedAt", fields)
	}
	if len(s.ChannelId) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "id", "value": s.ChannelId})
		fields = checkFields("id", fields)
	}
	if len(s.ChannelType) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "channeltype", "value": s.ChannelType})
		fields = checkFields("channelType", fields)
	}
	if len(s.TopicId) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "topicid", "value": s.TopicId})
		fields = checkFields("topicId", fields)
	}
	if len(s.RegionCode) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "country", "value": s.RegionCode})
		fields = checkFields("country", fields)
	}
	if len(s.RelevanceLanguage) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "relevancelanguage", "value": s.RelevanceLanguage})
		fields = checkFields("relevanceLanguage", fields)
	}
	if len(s.Sort) > 0 {
		sort = append(sort, map[string]interface{}{"field": strings.ToLower(s.Sort), "reverse": true})
		fields = checkFields(s.Sort, fields)
	}
	filter := map[string]interface{}{
		"should": should,
		"not":    not,
	}
	a := map[string]interface{}{
		"filter": filter,
		"query":  map[string]interface{}{"must": must},
		"sort":   sort,
	}
	if len(should) == 0 && len(not) == 0 {
		delete(a, "filter")
	} else {
		if len(should) == 0 {
			delete(filter, "should")
		}
		if len(not) == 0 {
			delete(filter, "not")
		}
	}
	if len(must) == 0 {
		delete(a, "query")
	}
	if len(sort) == 0 {
		delete(a, "sort")
	}
	queryObj, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	sql := fmt.Sprintf(`select %s from channel where expr(channel_index,'%s')`, strings.Join(fields, ","), queryObj)
	return sql, nil
}

func buildPlaylistSearch(s video.PlaylistSM, fields []string) (string, error) {
	var should []interface{}
	var must []interface{}
	var not []interface{}
	var sort []interface{}
	if len(s.Q) > 0 {
		should = append(should, map[string]interface{}{"type": "phrase", "field": "title", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "prefix", "field": "title", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "title", "value": fmt.Sprintf(`*%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "title", "value": fmt.Sprintf(`%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "phrase", "field": "description", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "prefix", "field": "description", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s*`, s.Q)})
	}
	if s.PublishedBefore != nil && s.PublishedAfter != nil {
		t1 := s.PublishedBefore.Format("2006-01-02 15:04:05")
		t2 := s.PublishedAfter.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "lower": t1, "upper": t2})
		fields = checkFields("publishedAt", fields)
	} else if s.PublishedAfter != nil {
		t2 := s.PublishedAfter.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "upper": t2})
		fields = checkFields("publishedAt", fields)
	} else if s.PublishedBefore != nil {
		t1 := s.PublishedBefore.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "lower": t1})
		fields = checkFields("publishedAt", fields)
	}
	if len(s.ChannelId) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "channelid", "value": s.ChannelId})
		fields = checkFields("channelId", fields)
	}
	if len(s.ChannelType) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "channeltype", "value": s.ChannelType})
		fields = checkFields("channelType", fields)
	}
	if len(s.RegionCode) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "country", "value": s.RegionCode})
		fields = checkFields("country", fields)
	}
	if len(s.RelevanceLanguage) > 0 {
		must = append(must, map[string]interface{}{"type": "match", "field": "relevancelanguage", "value": s.RelevanceLanguage})
		fields = checkFields("relevanceLanguage", fields)
	}
	if len(s.Sort) > 0 {
		sort = append(sort, map[string]interface{}{"field": strings.ToLower(s.Sort), "reverse": true})
	}
	filter := map[string]interface{}{
		"should": should,
		"not":    not,
	}
	a := map[string]interface{}{
		"filter": filter,
		"query":  map[string]interface{}{"must": must},
		"sort":   sort,
	}
	if len(should) == 0 && len(not) == 0 {
		delete(a, "filter")
	} else {
		if len(should) == 0 {
			delete(filter, "should")
		}
		if len(not) == 0 {
			delete(filter, "not")
		}
	}
	if len(must) == 0 {
		delete(a, "query")
	}
	if len(sort) == 0 {
		delete(a, "sort")
	}

	queryObj, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	sql := fmt.Sprintf(`select %s from playlist where expr(playlist_index,'%s')`, strings.Join(fields, ","), queryObj)
	log.Println(sql)
	return sql, nil
}

func buildVideosSearch(s video.ItemSM, fields []string) (string, error) {
	var should []interface{}
	var must []interface{}
	var not []interface{}
	var sort []interface{}
	if len(s.Duration) > 0 {
		switch s.Duration {
		case "short":
			must = append(must, map[string]interface{}{"type": "range", "field": "duration", "upper": "240"})
			break
		case "medium":
			must = append(must, map[string]interface{}{"type": "range", "field": "duration", "lower": "240", "upper": "1200"})
			break
		case "long":
			must = append(must, map[string]interface{}{"type": "range", "field": "duration", "lower": "1200"})
			break
		default:
			break
		}
		fields = checkFields("duration", fields)
	}
	if len(s.Q) > 0 {
		should = append(should, map[string]interface{}{"type": "phrase", "field": "title", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "prefix", "field": "title", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "title", "value": fmt.Sprintf(`*%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "title", "value": fmt.Sprintf(`%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "phrase", "field": "description", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "prefix", "field": "description", "value": fmt.Sprintf(`%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`%s*`, s.Q)})
		should = append(should, map[string]interface{}{"type": "wildcard", "field": "description", "value": fmt.Sprintf(`*%s*`, s.Q)})
	}
	if s.PublishedBefore != nil && s.PublishedAfter != nil {
		t1 := s.PublishedBefore.Format("2006-01-02 15:04:05")
		t2 := s.PublishedAfter.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "lower": t1, "upper": t2})
		fields = checkFields("publishedAt", fields)
	} else if s.PublishedAfter != nil {
		t2 := s.PublishedAfter.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "upper": t2})
		fields = checkFields("publishedAt", fields)
	} else if s.PublishedBefore != nil {
		t1 := s.PublishedBefore.Format("2006-01-02 15:04:05")
		must = append(must, map[string]interface{}{"type": "range", "field": "publishedat", "lower": t1})
		fields = checkFields("publishedAt", fields)
	}
	if len(s.RegionCode) > 0 {
		not = append(not, map[string]interface{}{"type": "match", "field": "blockedregions", "value": s.RegionCode})
	}
	if len(s.Sort) > 0 {
		sort = append(sort, map[string]interface{}{"field": strings.ToLower(s.Sort), "reverse": true})
		fields = checkFields(s.Sort, fields)
	}
	filter := map[string]interface{}{
		"should": should,
		"not":    not,
	}
	a := map[string]interface{}{
		"filter": filter,
		"query":  map[string]interface{}{"must": must},
		"sort":   sort,
	}
	if len(should) == 0 && len(not) == 0 {
		delete(a, "filter")
	} else {
		if len(should) == 0 {
			delete(filter, "should")
		}
		if len(not) == 0 {
			delete(filter, "not")
		}
	}
	if len(must) == 0 {
		delete(a, "query")
	}
	if len(sort) == 0 {
		delete(a, "sort")
	}
	queryObj, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	sql := fmt.Sprintf(`select %s from video where expr(video_index,'%s')`, strings.Join(fields, ","), queryObj)
	log.Println(sql)
	return sql, nil
}

func channelConvert(iter *gocql.Iter) (res []video.Channel, err error) {
	channel := video.Channel{}
	arrayColumn := iter.Columns()
	fieldDb := make([]string, len(arrayColumn))
	for i, v := range arrayColumn {
		fieldDb[i] = v.Name
	}
	column := common.ArrayValueForScansSQL(&channel, fieldDb, common.Cassandra)
	scanner := iter.Scanner()
	for scanner.Next() {
		err := scanner.Scan(column...)
		if err != nil {
			return nil, err
		}
		res = append(res, channel)
	}
	return
}

func playlistConvert(iter *gocql.Iter) (res []video.Playlist, err error) {
	playlist := video.Playlist{}
	arrayColumn := iter.Columns()
	fieldDb := make([]string, len(arrayColumn))
	for i, v := range arrayColumn {
		fieldDb[i] = v.Name
	}
	column := common.ArrayValueForScansSQL(&playlist, fieldDb, common.Cassandra)
	scanner := iter.Scanner()
	for scanner.Next() {
		err := scanner.Scan(column...)
		if err != nil {
			return nil, err
		}
		res = append(res, playlist)
	}
	return
}

func videoConvert(iter *gocql.Iter) (res []video.Video, err error) {
	video := video.Video{}
	arrayColumn := iter.Columns()
	fieldDb := make([]string, len(arrayColumn))
	for i, v := range arrayColumn {
		fieldDb[i] = v.Name
	}
	column := common.ArrayValueForScansSQL(&video, fieldDb, common.Cassandra)
	scanner := iter.Scanner()
	for scanner.Next() {
		err := scanner.Scan(column...)
		if err != nil {
			return nil, err
		}
		res = append(res, video)
	}
	return
}

func checkFields(check string, fields []string) []string {
	if len(fields) == 0 {
		return fields
	}
	flag := false
	for _, v := range fields {
		if v == check {
			flag = true
			break
		}
	}
	if !flag {
		fields = append(fields, check)
	}
	return fields
}
