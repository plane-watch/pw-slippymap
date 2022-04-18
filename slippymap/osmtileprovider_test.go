package slippymap

import "testing"

func TestGetTileAddress(t *testing.T) {

	TileProvider := &OSMTileProvider{}

	// check correct a... url is returned
	url, err := TileProvider.GetTileAddress(1, 2, 3)
	if err != nil {
		t.Error(err)
	}
	if url != "http://a.tile.openstreetmap.org/3/1/2.png" {
		t.Errorf("Expected http://a.tile.openstreetmap.org/3/1/2.png, got: %s", url)
	}

	// check correct b... url is returned
	url, err = TileProvider.GetTileAddress(1, 2, 3)
	if err != nil {
		t.Error(err)
	}
	if url != "http://b.tile.openstreetmap.org/3/1/2.png" {
		t.Errorf("Expected http://b.tile.openstreetmap.org/3/1/2.png, got: %s", url)
	}

	// check correct c... url is returned
	url, err = TileProvider.GetTileAddress(1, 2, 3)
	if err != nil {
		t.Error(err)
	}
	if url != "http://c.tile.openstreetmap.org/3/1/2.png" {
		t.Errorf("Expected http://c.tile.openstreetmap.org/3/1/2.png, got: %s", url)
	}

	// check correct a... url is returned (it loops back to a after c)
	url, err = TileProvider.GetTileAddress(1, 2, 3)
	if err != nil {
		t.Error(err)
	}
	if url != "http://a.tile.openstreetmap.org/3/1/2.png" {
		t.Errorf("Expected http://a.tile.openstreetmap.org/3/1/2.png, got: %s", url)
	}

}
