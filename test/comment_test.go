package test

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/shortdaddy0711/go-rest-api/internal/comment"
	"github.com/stretchr/testify/assert"
)

func TestGetComment(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/api/comment")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestPostComment(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"slug": "/", "author": "terrence", "body": "hello world"}`).Post(BASE_URL + "/api/comment")

	var comment1 comment.Comment

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
	json.Unmarshal(resp.Body(), &comment1)
	log.Println(comment1.ID)

	id := strconv.FormatUint(uint64(comment1.ID), 10)
	resp, err = client.R().Get(BASE_URL + "/api/comment/" + id)
	if err != nil {
		t.Fail()
	}
	var comment2 comment.Comment

	json.Unmarshal(resp.Body(), &comment2)

	assert.Equal(t, "terrence", comment2.Author)
}
