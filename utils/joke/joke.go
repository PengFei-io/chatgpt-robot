package joke

import (
	"fmt"
	"github.com/duke-git/lancet/v2/netutil"
	"strings"
)

// GetJokeList 对接的api https://www.mxnzp.com/doc/list
func GetJokeList() string {

	jokeApi := "https://www.mxnzp.com/api/jokes/list/random?app_id=qmnwjnimnkpnonom&app_secret=Nm5ZYk1JQVB3TW5lTittU0l1dzYzUT09"
	request := &netutil.HttpRequest{
		RawURL: jokeApi,
		Method: "GET",
	}

	httpClient := netutil.NewHttpClient()
	resp, err := httpClient.SendRequest(request)
	if err != nil || resp.StatusCode != 200 {
		_ = fmt.Errorf("GetJokeList Error: %v", err)
		return ""
	}
	ret := make(map[string]any)
	err = httpClient.DecodeResponse(resp, &ret)
	if err != nil {
		return ""
	}
	jokeSlice := make([]string, 0)
	contentSlice := ret["data"]
	for seq, v := range contentSlice.([]interface{}) {
		content := v.(map[string]interface{})["content"]
		jokeSlice = append(jokeSlice, fmt.Sprintf("笑话%d:%s \n", seq+1, content.(string)))
	}
	return strings.Join(jokeSlice, "\n")
}
