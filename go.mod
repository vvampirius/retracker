module github.com/vvampirius/retracker

go 1.21.0

replace (
	github.com/vvampirius/retracker/bittorrent/common => ./bittorrent/common
	github.com/vvampirius/retracker/bittorrent/response => ./bittorrent/response
	github.com/vvampirius/retracker/bittorrent/tracker => ./bittorrent/tracker
	github.com/vvampirius/retracker/common => ./common
)

require (
	github.com/vvampirius/retracker/bittorrent/common v0.0.0
	github.com/vvampirius/retracker/bittorrent/response v0.0.0
	github.com/vvampirius/retracker/bittorrent/tracker v0.0.0
	github.com/vvampirius/retracker/common v0.0.0
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/zeebo/bencode v1.0.0 // indirect
