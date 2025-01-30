# CalUnite

Sharing calenders has historically always been a pain, at least for me.
Especially when your goal is to combine a plethora of different calendars into one simple domain you can share with your significant other, relatives, or friend groups.

For the last few years, I've used [davorg/mergecal](https://github.com/davorg/mergecal) to combine calendars.
By now, there are also other alternatives available, like the misfortunately named [mergecal.org](https://mergecal.org/).
My problem with these solutions are, that they either don't provide *enough* features, or are heavily overengineered.

This is where **CalUnite** comes in:
- Combine RFC 5545 calendars (.ics) using urls or local paths
- Supports multiple calendars with a single config file
- Combined calendars are immediately served with a webserver
- Fast deployment using Docker
- Written in memory-safe Go

## Deploy

CalUnite uses docker for deployment. Preferrably, this is done via Compose:

`compose.yml`
```yml
# todo
environment: # OPTIONAL - the values below are the defaults
```

---

Files are served via HTTP. If you want HTTP**S**, this should be done via a reverse proxy like [Caddy](https://caddyserver.com/). A minimal example would look something like this:

`compose.yml`
```yml
# compose...
```

`caddy/Caddyfile`:
```
example.com {
    reverse_proxy calunite:8080
}
```
