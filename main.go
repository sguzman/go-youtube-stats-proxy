package main

import (
    "encoding/json"
    "fmt"
    "github.com/valyala/fasthttp"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
)

type ChannelListType struct {
    Kind string `json:"kind"`
    Etag string `json:"etag"`
    Id string `json:"id"`
    PageInfo PageInfoType `json:"pageInfo"`
    Items []ItemType `json:"items"`
}

type PageInfoType struct {
    TotalResults uint8 `json:"totalResults"`
    ResultsPerPage uint8 `json:"resultsPerPage"`
}

type ItemType struct {
    Kind string `json:"kind"`
    Etag string `json:"etag"`
    Id string `json:"id"`
    Statistics StatisticsType `json:"statistics"`
}

type StatisticsType struct {
    ViewCount string `json:"viewCount"`
    CommentCount string `json:"commentCount"`
    SubscriberCount string `json:"subscriberCount"`
    HiddenSubscriberCount bool `json:"hiddenSubscriberCount"`
    VideoCount string `json:"videoCount"`
}

func main() {
    key, exists := os.LookupEnv("PORT")
    if exists {
        port := ":" + key
        err := fasthttp.ListenAndServe(port, handle)
        if err != nil {
            panic(err)
        }
    } else {
        err := fasthttp.ListenAndServe(":8888", handle)
        if err != nil {
            panic(err)
        }
    }
}

func handle(ctx *fasthttp.RequestCtx) {
    var path = string(ctx.URI().Path())
    var split = strings.Split(path, "/")
    if len(split) == 3 {
        key := split[1]
        ids := split[2]
        fmt.Println(key, len(ids))

        var url =
            fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=statistics&id=%s&key=%s", ids, key)

        resp, err1 := http.Get(url)
        if err1 != nil {
            panic(err1)
        }

        bytes, err2 := ioutil.ReadAll(resp.Body)
        if err2 != nil {
            panic(err2)
        }

        var data ChannelListType

        err3 := json.Unmarshal(bytes, &data)
        if err3 != nil {
            panic(err3)
        }

        body := make([]string, 0)
        for _, item := range data.Items {
            subs := fmt.Sprintf("subscribers{channel=\"%s\"} %s", item.Id, item.Statistics.SubscriberCount)
            views := fmt.Sprintf("views{channel=\"%s\"} %s", item.Id, item.Statistics.ViewCount)
            videos := fmt.Sprintf("videos{channel=\"%s\"} %s", item.Id, item.Statistics.VideoCount)

            body = append(body, subs)
            body = append(body, views)
            body = append(body, videos)
        }

        bodyStr := strings.Join(body, "\n")
        count, err4 := fmt.Fprintln(ctx, bodyStr)
        if err4 != nil {
            panic(err4)
        }

        if count != len(bodyStr) + 1 {
            panic(fmt.Sprintf("Did not write all bytes: %d\n", count))
        }
    } else {
        fmt.Print(path)
        _, _ = fmt.Fprintln(ctx, "Hello, world!")
    }

}