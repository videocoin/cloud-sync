package service

import (
	"fmt"

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

func makePlaylistSegmentsKey(streamID string) string {
	return fmt.Sprintf("playlists/%s/segments", streamID)
}
