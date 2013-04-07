// Copyright 2013 Björn Oelke
//
// Licensed under the ICS license, see LICENSE file.
//
// m3u spec: http://schworak.com/blog/e39/m3u-play-list-specification/

package m3u

import (
	"bytes"
	"fmt"
)

type Song struct {
	Length int
	Title string
	Path string
}

type Playlist []Song

func (s Song) String() (string, error) {
	if len(s.Path) == 0 {
		return "", fmt.Errorf("Path is empty, title: %s", s.Title)
	}
	l := s.Length
	if l < -1 || l == 0 {
		l = -1
	}
	ret := fmt.Sprintf("#EXTINF:%d,%s\n%s\n", l, s.Title, s.Path)
	return ret, nil
}

func (p Playlist) String() (string, error) {
	var b bytes.Buffer
	b.WriteString("#EXTM3U\n")
	for i, song := range p {
		str, err := song.String()
		if err != nil {
			return "", fmt.Errorf("At index %d: %v", i, err)
		}
		b.WriteString(str)
	}
	return b.String(), nil
}
