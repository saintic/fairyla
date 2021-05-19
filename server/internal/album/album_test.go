package album

import (
	"fairyla/vars"
	"strings"
	"testing"

	"tcw.im/gtc"
)

func TestAlbum(t *testing.T) {
	ani_s := AlbumName2ID("u", "n")
	ani_d := vars.AlbumPre + "u" + gtc.MD5("un")
	if ani_s != ani_d {
		t.Fatal("AlbumName2ID error")
	}

	u := AlbumID2User(ani_s)
	if u != "u" {
		t.Fatal("AlbumID2User error")
	}

	a, _ := NewAlbum("u", "n")
	if AlbumID2User(a.ID) != a.Owner {
		t.Fatal("NewAlbum error")
	}
	a.AddLabel("x")
	if len(a.Label) != 1 || a.Label[0] != "x" {
		t.Fatal("AddLabel error")
	}
	a.AddLabel("y")
	a.AddLabel("z")
	a.RemoveLabel("y")
	if strings.Join(a.Label, "") != "xz" {
		t.Fatal("RemoveLabel error")
	}
}
