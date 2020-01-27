package service

import (
	"fmt"
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

func (ds *Datastore) AddSegment(streamID string, segmentNum int) error {
	k := makePlaylistSegmentsKey(streamID)
	err := ds.cli.SAdd(k, segmentNum).Err()
	return err
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
