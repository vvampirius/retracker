module github.com/vvampirius/common

go 1.21.0

replace (
	github.com/vvampirius/retracker/bittorrent/common => ../bittorrent/common
)

require (
	github.com/vvampirius/retracker/bittorrent/common v0.0.0
)
