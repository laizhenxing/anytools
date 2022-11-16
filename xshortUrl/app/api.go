package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"xshortUrl/app/db"
	"xshortUrl/tools"
)

// web 页面
func Home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	_ = t.Execute(w, nil)
}

func GetShortUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetShortUrl....")
	urlStr := r.FormValue("url_str")
	if urlStr == "" {
		fmt.Fprintf(w, "参数不存在")
		return
	}

	regexStr, err := regexp.Compile("^(http|https)*.*(cn|com|edu|gov|int|mil|net|org|biz|arpa|info).*$")
	if err != nil {
		fmt.Fprintf(w, "url不正确")
		return
	}
	if !regexStr.MatchString(urlStr) {
		fmt.Fprintf(w, "url格式不正确")
		return
	}

	// 规范化长链接，去除 Scheme 参数
	_url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	urlStr = _url.Host + _url.RequestURI()

	t := time.Now().UnixMilli()
	str := tools.Base62Encode(int(t))
	shortUrl := &db.ShortUrl{
		ShortCode: str,
		UrlStr: urlStr,
		Time: t,
	}
	err = shortUrl.Insert(db.MyDB)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	data := map[string]interface{}{
		"status": 200,
		"message": "Success",
		"data": map[string]string{
			"url": shortUrl.ShortCode,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		Home(w, r)
		return
	}
	w.Write(jsonData)
}

func ReIndex(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path
	rs := []rune(str)
	str = string(rs[1:])
	time, err := tools.Base62Decode(str)
	if err != nil {
		Home(w, r)
		return
	}
	shortUrl := &db.ShortUrl{
		Time: int64(time),
	}
	err = shortUrl.Select(db.MyDB)
	if err != nil {
		Home(w, r)
		return
	}

	w.Header().Set("Location", shortUrl.UrlStr)
	w.WriteHeader(302)
}

func Parse(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query().Get("short_url")
	if values == "" {
		w.Write([]byte("请输入正确的短链"))
		return
	}
	t, err := tools.Base62Decode(values)
	if err != nil {
		w.Write([]byte("无法解析短链，请重试"))
		return
	}
	shortUrl := &db.ShortUrl{
		Time: int64(t),
	}
	err = shortUrl.Select(db.MyDB)
	if err != nil {
		w.Write([]byte("无匹配，请重试"))
		return
	}
	w.Write([]byte("原始链接：" + shortUrl.UrlStr))
}
