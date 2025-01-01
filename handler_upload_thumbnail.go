package main

import (
	"fmt"
	"io"
	"net/http"

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
	r.ParseMultipartForm(maxMemory)
	// Get the image data
	formData, formHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Get Image Data", err)
		return
	}
	// Read all the image data
	imgData, err := io.ReadAll(formData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Image Data Read", err)
		return
	}
	// Get video's metadata
	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Get Video Metadata", err)
		return
	}
	thumb := thumbnail{
		data:      imgData,
		mediaType: formHeader.Header.Get("Content-Type"),
	}
	videoThumbnails[video.ID] = thumb
	url := fmt.Sprintf("/api/thumbnail/%s", video.ID)
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
