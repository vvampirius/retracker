# retracker

Simple HTTP torrent tracker.

* Keep all in memory (no persistent; doesn't require a database).
* Single binary executable (doesn't require a web-backend [apache, php-fpm, uwsgi, etc.])

## Installing

```
export GOPATH=$HOME/retracker
export PATH="$GOPATH/bin:$PATH"
go get github.com/vvampirius/retracker/...
```

## Usage
### Standalone

Start tracker on port 8080 with debug mode.
```
retracker -l :8080 -d
```
Add http://\<your ip>:8080/announce to your torrent.


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
