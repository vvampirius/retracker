module github.com/vvampirius/retracker/bittorrent/response

go 1.21.0

replace github.com/vvampirius/retracker/bittorrent/common => ../common

require (
	github.com/vvampirius/retracker/bittorrent/common v0.0.0
	github.com/zeebo/bencode v1.0.0
)
