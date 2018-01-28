package main

import (
	"fmt"
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

func (m *Media) DisplayDateTaken() {
	tm, _ := m.exif.DateTime()
	if !tm.IsZero() {
		fmt.Printf("%v: Taken on %v\n", m.fpath, tm)
	}
}
