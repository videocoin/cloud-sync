package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-redis/redis"
)

type Datastore struct {
	cli *redis.Client
}

func NewDatastore(cli *redis.Client) (*Datastore, error) {
	ds := &Datastore{
		cli: cli,
	}

	err := cli.Ping().Err()
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *Datastore) AddSegment(streamID string, segmentNum int, duration float64) error {
	k := makePlaylistSegmentsKey(streamID)
	segment := &Segment{Num: segmentNum, Duration: duration}
	segmentJSON, err := json.Marshal(segment)
	if err != nil {
		return err
	}
	err = ds.cli.SAdd(k, string(segmentJSON)).Err()
	return err
}

func (ds *Datastore) GetSegments(streamID string) ([]*Segment, error) {
	k := makePlaylistSegmentsKey(streamID)
	segments := []*Segment{}
	segmentsJSON, err := ds.cli.SMembers(k).Result()
	for _, segmentJSON := range segmentsJSON {
		segment := &Segment{}
		err := json.Unmarshal([]byte(segmentJSON), segment)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	sort.Sort(ByNum(segments))

	return segments, err
}

func (ds *Datastore) GetMaxSegment(streamID string) (int, error) {
	k := makePlaylistSegmentsKey(streamID)
	maxNum := 0
	segments, err := ds.cli.SMembers(k).Result()
	for _, segmentStr := range segments {
		num, err := strconv.Atoi(segmentStr)
		if err == nil {
			if maxNum < num {
				maxNum = num
			}
		}
	}
	return maxNum, err
}

func makePlaylistSegmentsKey(streamID string) string {
	return fmt.Sprintf("playlists/%s/segments", streamID)
}
