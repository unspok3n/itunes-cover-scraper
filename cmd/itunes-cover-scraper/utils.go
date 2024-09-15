package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	ErrExifMarkerNotFound    = errors.New("exif marker not found")
	ErrDownloadRequestFailed = errors.New("download request failed")
)

func DownloadFile(url string, filepath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrDownloadRequestFailed, resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func StripExif(buf []byte) ([]byte, error) {
	if buf[2] == 0xFF {
		switch buf[3] {
		case 0xE0:
			jfifSize := bytesToInt16(buf[4:6])
			if buf[jfifSize+4] == 0xFF && buf[jfifSize+5] == 0xE1 {
				return removeExifSegment(buf, int(jfifSize+4))
			}
			break
		case 0xE1:
			return removeExifSegment(buf, 2)
		}
	}
	return nil, ErrExifMarkerNotFound
}

func removeExifSegment(b []byte, exifMarkOffset int) ([]byte, error) {
	exifSize := bytesToInt16(b[exifMarkOffset+2 : exifMarkOffset+4])
	newImageSize := len(b) - int(exifSize)
	newBuf := make([]byte, 0, newImageSize)
	buffer := bytes.NewBuffer(newBuf)
	sum := exifMarkOffset + int(exifSize) + 2

	_, err := buffer.Write([]byte{0xFF, 0xD8})
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(b[sum:])
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func bytesToInt16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}
