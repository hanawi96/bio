package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type CloudinaryResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func UploadToCloudinary(file multipart.File, filename string) (string, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	uploadPreset := os.Getenv("CLOUDINARY_UPLOAD_PRESET")

	if cloudName == "" {
		return "", fmt.Errorf("cloudinary cloud name not configured")
	}
	if uploadPreset == "" {
		return "", fmt.Errorf("cloudinary upload preset not configured")
	}

	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", cloudName)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	// Add upload preset (unsigned mode)
	writer.WriteField("upload_preset", uploadPreset)
	writer.WriteField("folder", "linkbio/thumbnails")

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("cloudinary upload failed (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result CloudinaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
