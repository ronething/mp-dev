package storage

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/ronething/mp-dev/config"

	"github.com/imroc/req"
	"github.com/ronething/mp-dev/storage/trie"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

var NetEaseApiHost string

func InitThirdService() {
	NetEaseApiHost = config.Config.GetString("third.netease.host")
}

type (
	CheckMusicResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	Song struct {
		Id  int64  `json:"id"`
		URL string `json:"url"`
	}
	Artist struct {
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}
	Album struct {
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}
	SongDetail struct {
		Id     int64    `json:"id"`
		Name   string   `json:"name"`
		Artist []Artist `json:"ar"`
		Album  Album    `json:"al"`
	}
	SongUrlResp struct {
		Data []Song `json:"data"`
		Code int64  `json:"code"`
	}
	SearchSong struct {
		Id     int64    `json:"id"`
		Name   string   `json:"name"`
		Artist []Artist `json:"artists"`
	}
	SubSearchResult struct {
		Songs     []SearchSong `json:"songs"`
		HasMore   bool         `json:"hasMore"`
		SongCount int64        `json:"songCount"`
	}
	SearchMusicResp struct {
		Result SubSearchResult `json:"result"`
		Code   int64           `json:"code"`
	}
)

//PlayMusicBySongId 通过 sid 播放音乐
func PlayMusicBySongId(c *trie.Context) (*message.Reply, error) {
	// TODO: goroutine 并发请求 不然怕超过 5s 限制
	var (
		resp  *req.Resp
		err   error
		avail bool
		url   string
	)
	sid := c.Params["sid"]
	//1、检查音乐是否可用
	if avail, err = checkMusicAvailable(sid); err != nil {
		return nil, err
	} else if !avail {
		return c.Text(fmt.Sprintf("%s 音乐不可用", sid)), nil
	}
	//2、获取详情 /song/detail?ids=:sid
	if resp, err = httpClient.Get(
		fmt.Sprintf("%s/song/detail", NetEaseApiHost),
		req.QueryParam{"ids": sid}); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("get song detail resp is %s", resp.String())
	var songsDetail struct {
		Songs []SongDetail `json:"songs"`
	}
	if err = resp.ToJSON(&songsDetail); err != nil {
		return nil, err
	}
	log.Debugf("get song detail resp 序列化之后 %+v", songsDetail)
	if len(songsDetail.Songs) == 0 {
		// TODO: 会出现这种情况吗
		return c.Text("没有歌曲详情"), nil
	}
	song := songsDetail.Songs[0]
	artist := song.Artist
	var songArtist string
	if len(artist) > 0 {
		songArtist = artist[0].Name
	}
	album := song.Album
	log.Debugf("get sond deatil album is %+v", album)
	//go func() { // TODO: 上传封面
	//
	//}()
	//3、获取音乐 url
	if url, err = getSongURLById(sid); err != nil {
		return nil, err
	} else if url == "" {
		return c.Text(fmt.Sprintf("获取歌曲 %s 链接失败", sid)), nil
	}
	// 4、返回
	tmpId := config.Config.GetString("wechat.ThumbId") //只需要给一个有效的 media ID 即可
	return &message.Reply{
		MsgType: message.MsgTypeMusic,
		MsgData: message.NewMusic(song.Name, songArtist, url, url, tmpId),
	}, nil
}

//PlayMusicByName
func PlayMusicByName(c *trie.Context) (*message.Reply, error) {
	var (
		resp *req.Resp
		err  error
	)
	name := c.Params["name"]
	if resp, err = httpClient.Get(
		fmt.Sprintf("%s/search", NetEaseApiHost),
		req.QueryParam{
			"keywords": name,
			"limit":    1,
			"offset":   0,
			"type":     1, //单曲
		}); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("play search music resp is %+v", resp.String())
	var searchResp SearchMusicResp
	if err = resp.ToJSON(&searchResp); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("play search music resp 序列化之后 %+v", searchResp)
	if searchResp.Code != 200 {
		return c.Text(fmt.Sprintf("搜索音乐 %s 失败", name)), nil
	}
	// 列出列表 sid name
	var (
		songId     int64
		songName   string
		songArtist string
		url        string
		avail      bool
	)
	songs := searchResp.Result.Songs
	if len(songs) == 0 {
		return c.Text(fmt.Sprintf("搜索不到 **%s** 相关的音乐", name)), nil
	}
	songId = songs[0].Id
	if len(songs[0].Artist) > 0 {
		songArtist = songs[0].Artist[0].Name
	} else {
		songArtist = "佚名"
	}
	songName = songs[0].Name
	sid := fmt.Sprintf("%d", songId)
	// check music
	if avail, err = checkMusicAvailable(sid); err != nil {
		return nil, err
	} else if !avail {
		return c.Text(fmt.Sprintf("%s 音乐不可用", sid)), nil
	}
	// get song url
	if url, err = getSongURLById(sid); err != nil {
		return nil, err
	} else if url == "" {
		return c.Text(fmt.Sprintf("获取歌曲 %s 链接失败", sid)), nil
	}
	// 4、返回
	tmpId := config.Config.GetString("wechat.ThumbId") //只需要给一个有效的 media ID 即可
	return &message.Reply{
		MsgType: message.MsgTypeMusic,
		MsgData: message.NewMusic(songName, songArtist, url, url, tmpId),
	}, nil
}

//SearchMusicByKeyword
func SearchMusicByKeyword(c *trie.Context) (*message.Reply, error) {
	var (
		resp   *req.Resp
		err    error
		offset int
	)
	limit := config.Config.GetInt("third.netease.searchLimit")
	keyword := c.Params["keywords"]
	page := c.Params["page"]
	log.Debugf("page is %v", page)
	// page 初始化
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return c.Text("请输入正确页码"), nil
	}
	offset = (pageInt - 1) * limit
	if resp, err = httpClient.Get(
		fmt.Sprintf("%s/search", NetEaseApiHost),
		req.QueryParam{
			"keywords": keyword,
			"limit":    limit,
			"offset":   offset,
		}); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("search music resp is %+v", resp.String())
	var searchResp SearchMusicResp
	if err = resp.ToJSON(&searchResp); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("search music resp 序列化之后 %+v", searchResp)
	if searchResp.Code != 200 {
		return c.Text(fmt.Sprintf("搜索音乐 %s 失败", keyword)), nil
	}
	// 列出列表 sid name
	var (
		result     bytes.Buffer
		songId     int64
		songName   string
		songArtist string
	)
	songs := searchResp.Result.Songs
	if len(songs) == 0 {
		return c.Text(fmt.Sprintf("搜索不到 **%s** 相关的音乐", keyword)), nil
	}
	result.WriteString(fmt.Sprintf("搜索 *%s* 第 %s 页结果\n", keyword, page))
	for i := 0; i < len(songs); i++ {
		songId = songs[i].Id
		songName = songs[i].Name
		if len(songs[i].Artist) > 0 {
			songArtist = songs[i].Artist[0].Name
		} else {
			songArtist = "佚名"
		}
		result.WriteString(fmt.Sprintf("- %v %v %v\n", songId, songName, songArtist))
	}
	return c.Text(result.String()), nil
}

//GetSongURL 获取歌曲的下载链接
func GetSongURL(c *trie.Context) (*message.Reply, error) {
	sid := c.Params["sid"]
	// check music
	avail, err := checkMusicAvailable(sid)
	if err != nil {
		return nil, err
	}
	if !avail {
		return c.Text(fmt.Sprintf("%s 音乐不可用", sid)), nil
	}
	// get song url
	url, err := getSongURLById(sid) // err 在这里指的是内部错误
	if err != nil {
		return nil, err
	}
	if url == "" {
		return c.Text(fmt.Sprintf("获取歌曲 %s 链接失败", sid)), nil
	}
	return c.Text(url), nil

}

func checkMusicAvailable(sid string) (bool, error) {
	var (
		resp *req.Resp
		err  error
	)
	//DONE: check/music?id=:sid
	if resp, err = httpClient.Get(
		fmt.Sprintf("%s/check/music", NetEaseApiHost),
		req.QueryParam{"id": sid}); err != nil {
		log.Error(err)
		return false, err
	}
	log.Debugf("check music resp is %+v", resp.String())
	var checkResp CheckMusicResp
	if err = resp.ToJSON(&checkResp); err != nil {
		log.Error(err)
		return false, err
	}
	log.Debugf("check music resp 序列化之后 %+v", checkResp)
	if !checkResp.Success { //音乐不可用
		log.Errorf(fmt.Sprintf("音乐不可用 %s", sid))
		return false, nil
	}
	return true, nil
}

func getSongURLById(sid string) (string, error) {
	var (
		resp *req.Resp
		err  error
	)
	//get url /song/url?id=:sid
	if resp, err = httpClient.Get(
		fmt.Sprintf("%s/song/url", NetEaseApiHost),
		req.QueryParam{"id": sid}); err != nil {
		log.Error(err)
		return "", err
	}
	log.Debugf("get song url resp is %+v", resp.String())
	var songResp SongUrlResp
	if err = resp.ToJSON(&songResp); err != nil {
		log.Error(err)
		return "", err
	}
	log.Debugf("get song url resp 序列化之后 %+v", songResp)
	if songResp.Code != 200 {
		return "", errors.New(fmt.Sprintf("获取歌曲 %s 链接失败", sid))
	}
	var url string
	if len(songResp.Data) > 0 {
		url = songResp.Data[0].URL
		log.Infof("url is %v", url)
	} else {
		return "", errors.New(fmt.Sprintf("获取歌曲 %s 链接失败", sid))
	}
	return url, nil
}
