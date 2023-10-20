package vkapi

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type uploadServer struct {
	Url     string `json:"upload_url"`
	AlbumID int    `json:"album_id"`
	GroupID int    `json:"group_id"`
}

type uploadPhotoWallResponse struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type uploadServerResponse struct {
	Response uploadServer `json:"response"`
}

type savedPhoto struct {
	ID      int `json:"id"`
	OwnerID int `json:"owner_id"`
}

type savedPhotoResponse struct {
	Response []savedPhoto `json:"response"`
}

func (c *Client) getMessagesUploadServer(peerID int) (url string, albumID, groupID int) {
	r := c.Request("photos.getMessagesUploadServer", "peer_id="+strconv.Itoa(peerID))
	response := uploadServerResponse{}
	err := json.Unmarshal(r, &response)
	if err != nil {
		return "", 0, 0
	}
	CheckError(err)
	return response.Response.Url, response.Response.AlbumID, response.Response.GroupID
}

func (c *Client) saveMessagesPhoto(photo, hash string, server int) (mediaID, ownerID int) {
	params := fmt.Sprintf("photo=%s&server=%d&hash=%s", photo, server, hash)
	r := c.Request("photos.saveMessagesPhoto", params)
	resp := savedPhotoResponse{}
	err := json.Unmarshal(r, &resp)
	CheckError(err)
	if len(resp.Response) == 0 {
		return 0, 0
	}
	return resp.Response[0].ID, resp.Response[0].OwnerID
}

func UploadPhotoToServer(reader io.Reader, url string) (server int, photo, hash string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("photo", "test.png")
	CheckError(err)

	_, err = io.Copy(fileWriter, reader)
	CheckError(err)

	contentType := bodyWriter.FormDataContentType()
	err = bodyWriter.Close()
	CheckError(err)

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return 0, "", ""
	}
	CheckError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)

	var uploaded uploadPhotoWallResponse
	err = json.Unmarshal(body, &uploaded)
	if err != nil {
		log.Println(string(body))
		return 0, "", ""
	}
	CheckError(err)
	return uploaded.Server, uploaded.Photo, uploaded.Hash
}

func (c *Client) UploadPhotoToMessages(reader io.Reader, peerID int) (ownerID, mediaID int) {
	url, _, _ := c.getMessagesUploadServer(peerID)
	if url == "" {
		return 0, 0
	}

	server, photo, hash := UploadPhotoToServer(reader, url)
	if server == 0 {
		return 0, 0
	}
	mediaID, ownerID = c.saveMessagesPhoto(photo, hash, server)
	return
}

func (c *Client) UploadPhotoToMessagesFromPath(path string, peerID int) (ownerID, mediaID int) {
	file, err := os.Open(path)
	defer file.Close()
	CheckError(err)

	return c.UploadPhotoToMessages(file, peerID)
}
