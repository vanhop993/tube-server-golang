package pq

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go-service/video"
	"go-service/video/common"
	"go-service/video/youtube"
	"strconv"
	"strings"
)

type PostgreVideoService struct {
	db           *sql.DB
	tubeCategory youtube.CategoryTubeClient
}

func NewPostgreVideoService(db *sql.DB, tubeCategory youtube.CategoryTubeClient) *PostgreVideoService {
	return &PostgreVideoService{db: db, tubeCategory: tubeCategory}
}

func (s *PostgreVideoService) GetChannel(ctx context.Context, channelId string, fields []string) (*video.Channel, error) {
	if len(fields) == 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`select %s from channel where id = $1`, strings.Join(fields, ","))
	query, err := s.db.QueryContext(ctx, strq, channelId)
	if err != nil {
		return nil, err
	}
	arr, err := channelResult(query)
	if err != nil {
		return nil, err
	}
	if len(arr[0].ChannelList) > 0 {
		channels, err := s.GetChannels(ctx, arr[0].ChannelList, []string{})
		if err != nil {
			return nil, err
		}
		arr[0].Channels = *channels
	}
	return &arr[0], nil
}

func (s *PostgreVideoService) GetChannels(ctx context.Context, ids []string, fields []string) (*[]video.Channel, error) {
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = fmt.Sprintf(`$%d`, i+1)
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`Select %s from channel where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	query, err := s.db.Query(strq, cc...)
	if err != nil {
		return nil, err
	}
	res, err := channelResult(query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *PostgreVideoService) GetPlaylist(ctx context.Context, id string, fields []string) (*video.Playlist, error) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`Select %s from playlist where id = $1`, strings.Join(fields, ","))
	query, err := s.db.Query(strq, id)
	if err != nil {
		return nil, err
	}
	res, err := playlistResult(query)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (s *PostgreVideoService) GetPlaylists(ctx context.Context, ids []string, fields []string) (*[]video.Playlist, error) {
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = fmt.Sprintf(`$%d`, i+1)
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`Select %s from playlist where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	query, err := s.db.Query(strq, cc...)
	if err != nil {
		return nil, err
	}
	res, err := playlistResult(query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *PostgreVideoService) GetVideo(ctx context.Context, id string, fields []string) (*video.Video, error) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`Select %s from video where id = $1`, strings.Join(fields, ","))
	query, err := s.db.Query(strq, id)
	if err != nil {
		return nil, err
	}
	res, err := videoResult(query)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (s *PostgreVideoService) GetVideos(ctx context.Context, ids []string, fields []string) (*[]video.Video, error) {
	question := make([]string, len(ids))
	cc := make([]interface{}, len(ids))
	for i, v := range ids {
		question[i] = fmt.Sprintf(`$%d`, i+1)
		cc[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`Select %s from video where id in (%s)`, strings.Join(fields, ","), strings.Join(question, ","))
	query, err := s.db.Query(strq, cc...)
	if err != nil {
		return nil, err
	}
	res, err := videoResult(query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *PostgreVideoService) GetChannelPlaylists(ctx context.Context, channelId string, max int, nextPageToken string, fields []string) (*video.ListResultPlaylist, error) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	next := getNext(nextPageToken)
	strq := fmt.Sprintf(`select %s from playlist where channelId=$1 order by publishedAt desc limit %d offset %s`, strings.Join(fields, ","), max, next)
	fmt.Println(strq)
	query, err := s.db.Query(strq, channelId)
	if err != nil {
		return nil, err
	}
	playlists, err := playlistResult(query)
	if err != nil {
		return nil, err
	}
	var res video.ListResultPlaylist
	res.List = playlists
	lenList := len(res.List)
	r, err := strconv.Atoi(next)
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) GetChannelVideos(ctx context.Context, channelId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	next := getNext(nextPageToken)
	strq := fmt.Sprintf(`select %s from video where channelId=$1 order by publishedAt desc limit %d offset %s`, strings.Join(fields, ","), max, next)
	query, err := s.db.Query(strq, channelId)
	if err != nil {
		return nil, err
	}
	videos, err := videoResult(query)
	if err != nil {
		return nil, err
	}
	var res video.ListResultVideos
	res.List = videos
	lenList := len(res.List)
	r, err := strconv.Atoi(next)
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) GetPlaylistVideos(ctx context.Context, playlistId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	strq1 := `select videos from playlistVideo where id = $1 `
	pv, err := s.db.Query(strq1, playlistId)
	if err != nil {
		return nil, err
	}
	var videoIds []string
	for pv.Next() {
		pv.Scan(pq.Array(&videoIds))
	}
	fmt.Println(videoIds)
	questions := make([]string, len(videoIds))
	values := make([]interface{}, len(videoIds))
	for i, v := range videoIds {
		questions[i] = fmt.Sprintf(`$%d`, i+1)
		values[i] = v
	}
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	next := getNext(nextPageToken)
	strq := fmt.Sprintf(`select %s from video where id in (%s) order by publishedAt desc limit %d offset %s`, strings.Join(fields, ","), strings.Join(questions, ","), max, next)
	fmt.Println(strq)
	query, err := s.db.Query(strq, values...)
	if err != nil {
		return nil, err
	}
	videos, err := videoResult(query)
	if err != nil {
		return nil, err
	}
	var res video.ListResultVideos
	res.List = videos
	lenList := len(res.List)
	r, err := strconv.Atoi(next)
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) GetCagetories(ctx context.Context, regionCode string) (*video.Categories, error) {
	sql := `select * from category where id = $1`
	query, err := s.db.Query(sql, regionCode)
	if err != nil {
		return nil, err
	}
	var category video.Categories
	for query.Next() {
		query.Scan(&category.Id, &category.Data)
	}
	if category.Data == nil {
		res, er1 := s.tubeCategory.GetCagetories(regionCode)
		if er1 != nil {
			return nil, er1
		}
		resBytes, err := json.Marshal(*res)
		if err != nil {
			return nil, err
		}
		query := "insert into category (id,data) values ($1, $2)"
		_, err = s.db.Exec(query, regionCode, pq.Array(resBytes))
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

func (s *PostgreVideoService) SearchChannel(ctx context.Context, channelSM video.ChannelSM, max int, nextPageToken string, fields []string) (*video.ListResultChannel, error) {
	next := getNext(nextPageToken)
	strq, statement := buildChannelQuery(channelSM, fields)
	strq = strq + fmt.Sprintf(` limit %d offset %s`, max, next)
	fmt.Println(strq)
	rows, err := s.db.Query(strq, statement...)
	if err != nil {
		return nil, err
	}
	var res video.ListResultChannel
	channels, err := channelResult(rows)
	if err != nil {
		return nil, err
	}
	res.List = channels
	lenList := len(channels)
	r, err := strconv.Atoi(next)
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) SearchPlaylists(ctx context.Context, playlistSM video.PlaylistSM, max int, nextPageToken string, fields []string) (*video.ListResultPlaylist, error) {
	next := strings.Split(nextPageToken, "|")
	strq, statement := buildPlaylistQuery(playlistSM, fields)
	strq = strq + fmt.Sprintf(` limit %d offset %s`, max, next[0])
	rows, err := s.db.Query(strq, statement...)
	if err != nil {
		return nil, err
	}
	var res video.ListResultPlaylist
	playlists, err := playlistResult(rows)
	if err != nil {
		return nil, err
	}
	res.List = playlists
	lenList := len(playlists)
	r, err := strconv.Atoi(next[0])
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) SearchVideos(ctx context.Context, itemSM video.ItemSM, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	next := strings.Split(nextPageToken, "|")
	strq, statement := buildVideoQuery(itemSM, fields)
	strq = strq + fmt.Sprintf(` limit %d offset %s`, max, next[0])
	rows, err := s.db.Query(strq, statement...)
	if err != nil {
		return nil, err
	}
	var res video.ListResultVideos
	videos, err := videoResult(rows)
	if err != nil {
		return nil, err
	}
	res.List = videos
	lenList := len(videos)
	r, err := strconv.Atoi(next[0])
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) Search(ctx context.Context, itemSM video.ItemSM, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	next := strings.Split(nextPageToken, "|")
	strq, statement := buildVideoQuery(itemSM, fields)
	strq = strq + fmt.Sprintf(` limit %d offset %s`, max, next[0])
	rows, err := s.db.Query(strq, statement...)
	if err != nil {
		return nil, err
	}
	var res video.ListResultVideos
	videos, err := videoResult(rows)
	if err != nil {
		return nil, err
	}
	res.List = videos
	lenList := len(videos)
	r, err := strconv.Atoi(next[0])
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func (s *PostgreVideoService) GetRelatedVideos(ctx context.Context, videoId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	var a []string
	resVd, er1 := s.GetVideo(ctx, videoId, a)
	if er1 != nil {
		return nil, er1
	}
	if resVd == nil {
		return nil, errors.New("video don't exist")
	} else {
		if len(fields) <= 0 {
			fields = append(fields, "*")
		}
		next := strings.Split(nextPageToken, "|")
		sql := fmt.Sprintf(`select %s from video order by publishedAt desc limit %d offset %s`, strings.Join(fields, ","), max, next[0])
		rows, err := s.db.Query(sql)
		if err != nil {
			return nil, err
		}
		var res video.ListResultVideos
		videos, err := videoResult(rows)
		if err != nil {
			return nil, err
		}
		res.List = videos
		lenList := len(videos)
		r, err := strconv.Atoi(next[0])
		if err != nil {
			return nil, errors.New("nextPageToken wrong")
		}
		res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
		return &res, nil
	}
}

func (s *PostgreVideoService) GetPopularVideos(ctx context.Context, regionCode string, categoryId string, max int, nextPageToken string, fields []string) (*video.ListResultVideos, error) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	next := strings.Split(nextPageToken, "|")
	var condition []string
	var args []interface{}
	var where string
	i := 0
	if len(regionCode) > 0 {
		condition = append(condition, fmt.Sprintf(`(blockedRegions is null or $%d != all(blockedRegions))`, i+1))
		args = append(args, regionCode)
	}
	if len(categoryId) > 0 {
		condition = append(condition, fmt.Sprintf(`(categoryId = $%d)`, i+1))
		args = append(args, categoryId)
	}
	if len(condition) > 0 {
		where += fmt.Sprintf(` where %s`, strings.Join(condition, " and "))
	}
	sql := fmt.Sprintf(`select %s from video %s order by publishedAt desc limit %d offset %s`, strings.Join(fields, ","), where, max, next[0])
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	var res video.ListResultVideos
	videos, err := videoResult(rows)
	if err != nil {
		return nil, err
	}
	res.List = videos
	lenList := len(videos)
	r, err := strconv.Atoi(next[0])
	if err != nil {
		return nil, errors.New("nextPageToken wrong")
	}
	res.NextPageToken = createNextPageToken(lenList, max, r, res.List[lenList-1].Id, "")
	return &res, nil
}

func channelResult(query *sql.Rows) ([]video.Channel, error) {
	var res []video.Channel
	channel := video.Channel{}
	fieldDb, err := query.Columns()
	if err != nil {
		return nil, err
	}
	columns := common.ArrayValueForScansSQL(&channel, fieldDb, common.Postgre)
	for query.Next() {
		err = query.Scan(columns...)
		if err != nil {
			return nil, err
		}
		res = append(res, channel)
	}
	return res, nil
}

func playlistResult(query *sql.Rows) ([]video.Playlist, error) {
	var res []video.Playlist
	var playlist video.Playlist
	fieldDb, err := query.Columns()
	if err != nil {
		return nil, err
	}
	columns := common.ArrayValueForScansSQL(&playlist, fieldDb, common.Postgre)
	for query.Next() {
		err := query.Scan(columns...)
		if err != nil {
			return nil, err
		}
		res = append(res, playlist)
	}
	return res, nil
}

func videoResult(query *sql.Rows) ([]video.Video, error) {
	var res []video.Video
	var video video.Video
	fieldDb, err := query.Columns()
	if err != nil {
		return nil, err
	}
	columns := common.ArrayValueForScansSQL(&video, fieldDb, common.Postgre)
	for query.Next() {
		err := query.Scan(columns...)
		if err != nil {
			return nil, err
		}
		res = append(res, video)
	}
	return res, nil
}

func createNextPageToken(lenList int, limit int, skip int, id string, name string) string {
	if len(name) <= 0 {
		name = "id"
	}
	if lenList < limit {
		return ""
	} else {
		//return list && list.length > 0 ? `${list[list.length - 1][name]}|${skip + limit}` : undefined;
		if lenList > 0 {
			return fmt.Sprintf(`%d|%s`, skip+limit, id)
		} else {
			return ""
		}
	}
}

func buildChannelQuery(s video.ChannelSM, fields []string) (string, []interface{}) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`select %s from channel`, strings.Join(fields, ","))
	var condition []string
	var params []interface{}
	i := 1
	if len(s.ChannelId) > 0 {
		params = append(params, s.ChannelId)
		condition = append(condition, fmt.Sprintf(`id = $%d`, i))
		i++
	}
	if len(s.RegionCode) > 0 {
		params = append(params, s.RegionCode)
		condition = append(condition, fmt.Sprintf(`country = $%d`, i))
		i++
	}
	if s.PublishedAfter != nil {
		params = append(params, s.PublishedAfter)
		condition = append(condition, fmt.Sprintf(`publishedAt <= $%d`, i))
		i++
	}
	if s.PublishedBefore != nil {
		params = append(params, s.PublishedBefore)
		condition = append(condition, fmt.Sprintf(`publishedAt > $%d`, i))
		i++
	}
	if len(s.Q) > 0 {
		q := "%" + s.Q + "%"
		params = append(params, q, q)
		condition = append(condition, fmt.Sprintf(`(title ilike $%d or description ilike $%d)`, i, i+1))
	}

	if len(condition) > 0 {
		cond := strings.Join(condition, " and ")
		strq += fmt.Sprintf(` where %s`, cond)
	}
	if len(s.Sort) > 0 {
		strq += fmt.Sprintf(` order by %s desc`, s.Sort)
	}
	return strq, params
}

func buildPlaylistQuery(s video.PlaylistSM, fields []string) (string, []interface{}) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`select %s from playlist`, strings.Join(fields, ","))
	var condition []string
	var params []interface{}
	i := 1
	if len(s.ChannelId) > 0 {
		params = append(params, s.ChannelId)
		condition = append(condition, fmt.Sprintf(`id = $%d`, i))
		i++
	}
	if s.PublishedAfter != nil {
		params = append(params, s.PublishedAfter)
		condition = append(condition, fmt.Sprintf(`publishedAt <= $%d`, i))
		i++
	}
	if s.PublishedBefore != nil {
		params = append(params, s.PublishedBefore)
		condition = append(condition, fmt.Sprintf(`publishedAt > $%d`, i))
		i++
	}
	if len(s.Q) > 0 {
		q := "%" + s.Q + "%"
		params = append(params, q, q)
		condition = append(condition, fmt.Sprintf(`(title ilike $%d or description ilike $%d)`, i, i+1))
	}
	if len(condition) > 0 {
		cond := strings.Join(condition, " and ")
		strq += fmt.Sprintf(` where %s`, cond)
	}
	if len(s.Sort) > 0 {
		strq += ` order by ${s.sort} desc`
	}
	return strq, params
}

func buildVideoQuery(s video.ItemSM, fields []string) (string, []interface{}) {
	if len(fields) <= 0 {
		fields = append(fields, "*")
	}
	strq := fmt.Sprintf(`select %s from video`, strings.Join(fields, ","))
	var condition []string
	var params []interface{}
	i := 1
	if s.PublishedAfter != nil {
		params = append(params, s.PublishedAfter)
		condition = append(condition, fmt.Sprintf(`publishedAt <= $%d`, i))
		i++
	}
	if s.PublishedBefore != nil {
		params = append(params, s.PublishedBefore)
		condition = append(condition, fmt.Sprintf(`publishedAt > $%d`, i))
		i++
	}
	if len(s.RegionCode) > 0 {
		params = append(params, s.RegionCode)
		condition = append(condition, fmt.Sprintf(`country = $%d`, i))
		i++
	}
	if len(s.Duration) > 0 {
		switch s.Duration {
		case "short":
			condition = append(condition, `duration <= 240`)
			break
		case "medium":
			condition = append(condition, `(duration > 240 and duration <= 1200)`)
			break
		case "long":
			condition = append(condition, `(duration > 1200)`)
			break
		default:
			break
		}
	}
	if len(s.Q) > 0 {
		q := "%" + s.Q + "%"
		params = append(params, q, q)
		condition = append(condition, fmt.Sprintf(`(title ilike $%d or description ilike $%d)`, i, i+1))
	}
	if len(condition) > 0 {
		cond := strings.Join(condition, " and ")
		strq += fmt.Sprintf(` where %s`, cond)
	}
	if len(s.Sort) > 0 {
		strq += ` order by ${s.sort} desc`
	}
	return strq, params
}

func getNext(nextPageToken string) (next string) {
	if len(nextPageToken) > 0 {
		next = strings.Split(nextPageToken, "|")[0]
	} else {
		next = "0"
	}
	return
}
