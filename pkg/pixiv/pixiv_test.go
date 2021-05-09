package pixiv

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPixiv(t *testing.T) {
	pixiv := NewPixiv("", http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	})

	list, err := pixiv.RankList("weekly", 1)
	fmt.Println(list, err)
}
