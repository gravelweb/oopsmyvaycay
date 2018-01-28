package main

import (
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

type Media struct {
	fpath string
	file  *os.File
	exif  *exif.Exif
}

func NewMedia(fpath string) (*Media, error) {
	file, err := os.Open(fpath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	x, err := exif.Decode(file)
	if err != nil {
		return nil, err
	}

	return &Media{fpath, file, x}, nil
}

func (m *Media) DateStamp() string {
	tm, _ := m.exif.DateTime()
	return tm.Format("2006-01")
}
