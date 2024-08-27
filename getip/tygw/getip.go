package tygw

import (
	"encoding/json"
	"fmt"
	"github.com/attson/ddns/getip/common"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func init() {
	common.RegisterGetIp("tygw", GetIp)
}

func retryOnRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("too many redirects")
	}

	if len(via) > 0 {
		lastRequest := via[len(via)-1]
		redirectURL := req.URL

		if lastRequest.URL.Host == redirectURL.Host && lastRequest.URL.Path == redirectURL.Path {
			query := lastRequest.URL.Query()
			redirectURL.RawQuery = query.Encode()
		}
	}

	for _, cookie := range via[0].Cookies() {
		req.AddCookie(cookie)
	}

	return nil
}

func login(conf map[string]interface{}) (string, error) {
	data := url.Values{}
	data.Set("username", conf["username"].(string))
	data.Set("psd", conf["password"].(string))
	request, err := http.NewRequest("POST", conf["base_url"].(string)+"cgi-bin/luci/", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Origin", conf["base_url"].(string))
	request.Header.Add("Referer", conf["base_url"].(string)+"/cgi-bin/luci")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar:           cookieJar,
		CheckRedirect: retryOnRedirect,
	}

	_, err = client.Do(request)
	if err != nil {
		return "", err
	}

	for _, cookie := range cookieJar.Cookies(request.URL) {
		if cookie.Name == "sysauth" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("login failed")
}

func GetIp(conf map[string]interface{}) (string, error) {
	auth, err := login(conf)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("GET", conf["base_url"].(string)+"cgi-bin/luci/admin/settings/gwinfo?get=part&_=0.6667750936618431", nil)
	if err != nil {
		return "", err
	}

	request.AddCookie(&http.Cookie{Name: "sysauth", Value: auth})
	request.AddCookie(&http.Cookie{Name: "/", Value: request.Host})

	request.Header.Add("Referer", conf["base_url"].(string)+"cgi-bin/luci/admin/home")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{})

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		all, _ := io.ReadAll(resp.Body)
		fmt.Println(string(all))
		return "", err
	}

	return result["WANIP"].(string), nil
}
