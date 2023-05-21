package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"bytes"
	"io"
)

type Category struct {
	Id int `json:"id"`
	Name string `json:"name"`
	IsCustom string `json:"is_custom"`
	PictureUrl string `json:"picture_url"`
}

type Hashtag struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Mentions int `json:"mentions"`
}

type PostToCategory struct {
	Id int `json:"id"`
	IdPost int `json:"id_post"`
	IdCategory int `json:"id_category"`
}

type PostToHashtag struct {
	Id int `json:"id"`
	IdPost int `json:"id_post"`
	IdHashtag int `json:"id_hashtag"`
}

func makeRequest(client *http.Client, method string, uri string, reader io.Reader) *http.Response {
	uri = strings.Replace(uri, " ", "%20", -1)
	Log("Processing PPS-API", method, "-", uri)
	req, _ := http.NewRequest(method, uri, reader)
	req.Header.Set("Authorization", "Bearer " + apiToken)

	resp, err := client.Do(req)
	if err != nil {
		Log("Error occured while PPS-API", method, "-", uri)
	}

	return resp
}

func FetchCategoryFromAPI(category string) int {
	client := http.Client{}
	resp := makeRequest(&client, "GET", apiAddr + "/rest/categories?name=" + category, nil)

	var instance []Category
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyByte, &instance)
	if len(instance) < 1 {
		return 1
	} else {
		return instance[0].Id
	}
}

func FetchHashtagsFromAPI(hashtagsOneline string) []int {
	hashtags := strings.Split(hashtagsOneline, " ")

	client := http.Client{}
	ret := []int{}

	for _, hashtag := range(hashtags) {
		for {
			resp := makeRequest(&client, "GET", apiAddr + "/rest/hashtags?name=" + hashtag, nil)

			bodyByte, _ := ioutil.ReadAll(resp.Body)
			var instance []Hashtag
			json.Unmarshal(bodyByte, &instance)
			
			if len(instance) < 1 {
				payload, _ := json.Marshal(Hashtag{0, hashtag, 0})
				resp = makeRequest(&client, "POST", apiAddr + "/rest/hashtags", bytes.NewReader(payload))
			} else {
				ret = append(ret, instance[0].Id)
				break
			}
		}
	}

	return ret
}

func PostPushToAPI(post Post) {
	Log("Pushing PPS-API:", post.Title, "- Author:", post.IdAuthor)

	categoryId := FetchCategoryFromAPI(post.Category)
	hashtagIds := FetchHashtagsFromAPI(post.Hashtags)

	client := http.Client{}

	for {
		resp := makeRequest(&client, "GET", apiAddr + "/rest/posts?title=" + post.Title + "&id_author=" + strconv.Itoa(post.IdAuthor), nil)

		bodyByte, _ := ioutil.ReadAll(resp.Body)
		Log(string(bodyByte))
		var posts []Post
		json.Unmarshal(bodyByte, &posts)
		Log(posts)

		if len(posts) < 1 {
			payload, _ := json.Marshal(post)
			resp = makeRequest(&client, "POST", apiAddr + "/rest/posts", bytes.NewReader(payload))
		} else {
			payload, _ := json.Marshal(PostToCategory{0, posts[0].Id, categoryId})
			makeRequest(&client, "POST", apiAddr + "/rest/post_to_category", bytes.NewReader(payload))

			for _, hashtagId := range(hashtagIds) {
				payload, _ = json.Marshal(PostToHashtag{0, posts[0].Id, hashtagId})
				makeRequest(&client, "POST", apiAddr + "/rest/post_to_hashtag", bytes.NewReader(payload))
			}
			break
		}
	}
}
