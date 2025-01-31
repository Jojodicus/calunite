[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![MIT](https://img.shields.io/badge/MIT-green?style=for-the-badge)](https://opensource.org/license/mit)

# 🗓️ CalUnite

Sharing calendars has historically always been a pain, at least for me.
Especially when your goal is to combine a plethora of different calendars into one simple domain you can share with your significant other, relatives, or friend groups.

For the last few years, I've used [davorg/mergecal](https://github.com/davorg/mergecal) to combine calendars.
By now, there are also other alternatives available, like the misfortunately named [mergecal.org](https://mergecal.org/).
My problem with these solutions are, that they either don't provide *enough* features, or are heavily overengineered.

This is where **CalUnite** comes in:
- Combine [RFC 5545](https://datatracker.ietf.org/doc/html/rfc5545) calendars (.ics) using URLs or local paths
- Supports multiple calendars with a single config file
- Recursive merging
- Combined calendars are immediately served with a webserver
- Fast deployment using [Docker](https://www.docker.com/)
- Written in memory-safe [Go](https://go.dev/)

## ✏️ Configuration

A sample configuration file can be found at [config.yml](config.yml).
It's recommended to read the documentation in there.
Additional settings can be tweaked using the container's environment variables, see the following [Deployment](#️-deployment) section for more details.

## 🏎️ Deployment

CalUnite uses [Docker](https://www.docker.com/) for deployment. Preferably, this is done via [Compose](https://docs.docker.com/compose/):

`compose.yml`
```yml
services:
  calunite:
    image: jojodicus/calunite
    container_name: calunite
    restart: unless-stopped
    volumes:
      - ./calunite:/config/
    ports:
      - 8080:8080
    environment: # OPTIONAL - the values below are the defaults
      - # todo
```

---

Files are served via HTTP. If you want HTTP**S**, this should be done via a reverse proxy like [Caddy](https://caddyserver.com/). A minimal example would look something like this:

`compose.yml`
```yml
services:
  calunite:
    image: jojodicus/calunite
    container_name: calunite
    restart: unless-stopped
    volumes:
      - ./calunite:/config/
    # no exposed ports needed

  caddy:
    image: caddy
    container_name: caddy
    restart: unless-stopped
    volumes:
      - ./caddy/:/etc/caddy/
    ports:
      - 80:80
      - 443:443
```

`caddy/Caddyfile`:
```
example.com {
    reverse_proxy calunite:8080
}
```

If you're looking to host this without having a public IP (typical for most home internet connections), you could use something like a [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/).
A curated guide on how to set this up can be found [here](https://dittrich.pro/cloudflare-tunnel-homelab/).
