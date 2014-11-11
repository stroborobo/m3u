// Copyright 2013 Björn Oelke
//
// Licensed under the ICS license, see LICENSE file.
//
// m3u spec: http://schworak.com/blog/e39/m3u-play-list-specification/

package m3u

import (
	"bytes"
	"fmt"
	"io"
	"bufio"
	"strings"
	"strconv"
)

type Song struct {
	Length int64
	Title string
	Path string
}

type Playlist []*Song

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

func Parse(r io.Reader) (Playlist, error) {
	b := bufio.NewReader(r)
	p := make(Playlist, 0, 10)
	s := new(Song)
	for {
		bytes, err := b.ReadBytes('\n')
		last := false
		if err == io.EOF {
			last = true
		} else if err != nil {
			return nil, err
		}
		line := strings.TrimSpace(string(bytes))
		length := len(line)
		if length >= 7 && line[0:7] == "#EXTINF" {
			infos := strings.SplitN(line[8:], ",", 2)
			length, err := strconv.ParseInt(infos[0], 10, 64)
			if err != nil {
				return nil, err
			}
			s.Length = length
			s.Title = infos[1]
		} else if length >= 1 && line[0] == '#' {
			continue
		} else if length != 0 {
			s.Path = line
			p = append(p, s)
			s = new(Song)
		}
		if last {
			break
		}
	}
	return p, nil
}
