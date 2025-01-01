package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	// Parse Form data
	const maxMemory = 10 << 20
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Max Memory", err)
		return
	}
	// Get the image data
	formData, formHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Get Image Data", err)
		return
	}
	defer formData.Close()
	contentType := formHeader.Header.Get("Content-Type")
	// Read all the image data
	// imgData, err := io.ReadAll(formData)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Image Data Read", err)
	// 	return
	// }
	// Get video's metadata
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Get Video Metadata", err)
		return
	}
	// thumb := thumbnail{
	//     data: imgData,
	// 	mediaType: contentType,
	// }
	// videoThumbnails[video.ID] = thumb
	fileExt := strings.Split(contentType, "/")
	videoFile := fmt.Sprintf("%s.%s", videoID, fileExt[1])
	videoFilePath := filepath.Join(cfg.assetsRoot, videoFile)
	vFile, err := os.Create(videoFilePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Create File", err)
		return
	}
	defer vFile.Close()
	url := fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, videoFile)

	_, err = io.Copy(vFile, formData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Copy File", err)
		return
	}
	err = vFile.Sync()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Sync File", err)
		return
	}
	video.ThumbnailURL = &url
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Updating Video", err)
		return
	}
	updateV, err := cfg.db.GetVideo(video.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Get Updated Vedio Data", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updateV)
}
