package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func HttpGet(url string, header map[string]string, proxy func(ctx context.Context, network, addr string) (net.Conn, error)) ([]byte, error) {
	c := http.Client{
		Transport: &http.Transport{
			DialContext: proxy,
		},
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

var cfCache = make(map[string]string)

func CloudflareResolve(hostname string) (string, error) {
	if v, ok := cfCache[hostname]; ok {
		return v, nil
	}
	resp, err := http.Get("https://cloudflare-dns.com/dns-query?name=" + hostname + "&ct=application/dns-json&type=A&do=false&cd=false")
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	if err := json.Unmarshal(body, &m); err != nil {
		return "", err
	}
	tmp, ok := m["Answer"].([]interface{})
	if !ok {
		return "", errors.New("error asnwer")
	}
	if len(tmp) == 0 {
		return "", errors.New("0 answer")
	}
	for _, v := range tmp {
		if v.(map[string]interface{})["type"].(float64) == 1 {
			return v.(map[string]interface{})["data"].(string), nil
		}
	}
	return "", errors.New("0 answer")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func StringToInt(i string) int {
	ret, _ := strconv.Atoi(i)
	return ret
}

func StringToInt64(i string) int64 {
	ret, _ := strconv.ParseInt(i, 10, 64)
	return ret
}

func RegexMatch(content string, command string) []string {
	reg := regexp.MustCompile(command)
	return reg.FindStringSubmatch(content)
}

func RegexMatchs(content string, command string) [][]string {
	reg := regexp.MustCompile(command)
	return reg.FindAllStringSubmatch(content, -1)
}
