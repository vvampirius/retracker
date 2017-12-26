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

## Behind NGINX
Configure nginx like:
```
# cat /etc/nginx/sites-enabled/retracker.local
server {
        listen 80;

        server_name retracker.local;

        access_log /var/log/nginx/retracker.local-access.log;

        proxy_set_header X-Real-IP $remote_addr;

        location / {
                proxy_pass http://localhost:8080;
        }
}
```

Start tracker on port 8080 with getting remote address from X-Real-IP header.
```
retracker -l :8080 -x
```

Add retracker.local to your local DNS or /etc/hosts.

Add http://retracker.local/announce to your torrent.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
