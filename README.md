# retracker

Simple HTTP torrent tracker.

* Keep all in memory (no persistent; doesn't require a database).
* Single binary executable (doesn't require a web-backend [apache, php-fpm, uwsgi, etc.])
* Can collect peers from external trackers (HTTP only)
* Expose some metrics for Prometheus monitoring

## Installing

```
go install 'github.com/vvampirius/retracker@latest'
```
> Executables are installed in the directory named by the GOBIN environment variable, which defaults to $GOPATH/bin or $HOME/go/bin if the GOPATH environment variable is not set. Executables in $GOROOT are installed in $GOROOT/bin or $GOTOOLDIR instead of $GOBIN.

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

        location /metrics {
                allow 10.0.0.0/8;
                deny  all;
                proxy_pass http://localhost:8080;
        }

        location / {
                proxy_pass http://localhost:8080;
        }
}
```

Start tracker on port 8080 with getting remote address from X-Real-IP header.
```
retracker -l :8080 -x -p
```

Add retracker.local to your local DNS or /etc/hosts.

Add http://retracker.local/announce to your torrent.

### Standalone with announce forwarding

You can forward announce request to some external HTTP trackers and append peers from them to response to your torrent client.
```
retracker -l :8080 -d -f forwarders.yml
```
forwarders.yml:
```
- uri: http://1.2.3.4:8080/announce
- uri: http://5.6.7.8:8080/announce
- uri: http://5.6.7.8:8080/announce
  ip: 192.168.1.15 # announce different torrent client IP to this forwarder
```
Add http://\<your ip>:8080/announce to your torrent.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
