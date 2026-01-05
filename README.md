[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![MIT](https://img.shields.io/badge/MIT-green?style=for-the-badge)](https://opensource.org/license/mit)

[![Go Report Card](https://goreportcard.com/badge/github.com/Jojodicus/calunite)](https://goreportcard.com/report/github.com/Jojodicus/calunite)
[![Docker Image](https://github.com/jojodicus/calunite/actions/workflows/docker-image.yml/badge.svg)](https://hub.docker.com/r/jojodicus/calunite/tags)

# 🗓️ [CalUnite](https://hub.docker.com/r/jojodicus/calunite)

Sharing calendars has historically always been a pain, at least for me.
Especially when your goal is to combine a plethora of different calendars into one simple domain you can share with your significant other, relatives, or friend groups.

For the last few years, I've used [davorg/mergecal](https://github.com/davorg/mergecal) to combine calendars.
By now, there are also other alternatives available, like the misfortunately named [mergecal.org](https://mergecal.org/).
My problem with these solutions are, that they either don't provide *enough* features, or are heavily overengineered.

This is where **CalUnite** comes in:
- Combine [RFC 5545](https://datatracker.ietf.org/doc/html/rfc5545) calendars (.ics) using URLs or local paths
- Supports multiple calendars with a single config file
- Custom formatting of event titles for different calendars
- Calendar aliases (private or publicly served)
- Recursive merging (with cycle detection)
- Combined calendars are immediately served with a webserver
- Hot-reload configuration, no need to stop the container
- Fast deployment using [Docker](https://www.docker.com/)
- Written in memory-safe [Go](https://go.dev/)

## ✏️ Configuration

A sample configuration file can be found at [config.yml](config.yml).
It's recommended to read the documentation in there.
Additional settings can be tweaked using the container's environment variables, see the following [Deployment](#️-deployment) section for more details.

If you want to use a hot-reloadable configuration, make sure to use a bind mount (directory path in volume specifier, not a regular file), as is also shown in the examples. This ensures that the file is synced between host and container.

During a hot-reload, the **old content directory will be cleared entirely**, keep that in mind when mounting your own volumes (`/wwwdata` by default).
If you want to include files other than ones specified in the 'config.yml', it's advised to do that via a [reverse proxy](#-reverse-proxy) or as a sub-path from your existing webserver.

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
      CFG_PATH: /config/config.yml # path of config file within the container
      CRON: "@every 15m"           # how often the merger should run, format: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
      PROD_ID: CalUnite            # RFC 5545 PRODID, who created the calendar
      CONTENT_DIR: /wwwdata        # directory from which the files are served, no change needed in most cases
      FILE_NAVIGATION: false       # generate an index.html for directories to allow for navigation
      DOT_PRIVATE: true            # if files prefixed with a dot (.) should not be served, only applicable when FILE_NAVIGATION is false
      ADDR: 0.0.0.0                # address to bind to
      PORT: 8080                   # port to expose
```

### 🖧 Reverse Proxy

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

A Caddy configuration with files other than calendars managed by CalUnite may look like this.
Adjust to your own needs:

`caddy/Caddyfile`:
```
example.com {
    handle_path /*.ics {
        reverse_proxy calunite:8080
    }

    # or use your own webserver
    handle {
        root * /path/to/mounted/website
        file_server
    }
}
```

## 🔗 How to get iCal links

### Google Calendar

For Google, this is only possible on the desktop webiste.
If you are on mobile, you can change the view at the bottom of the page.

Go to your [Google Calendar Settings](https://calendar.google.com/calendar/u/0/r/settings) and click on your desired calendar.
Here, you can copy the "address in iCal format".

For non-public calendars, you'll need the secret address.

### Apple Calendar

Go to your [Apple Calendar](https://www.icloud.com/calendar/).
Then, click on the person icon next to a calendar (will appear when hovering over the name).

On mobile, press the calendar icon at the bottem of the page, then press the "i" icon next to a calendar.

Make the calendar public and copy the link.

## 📡 Subscribing to calendars

### Google Calendar

For Google, this is only possible on the desktop website.
If you are on mobile, you can change the view at the bottom of the page.
Go to [Google Calendar - Add from URL](https://calendar.google.com/calendar/u/0/r/settings/addbyurl) and paste the URL.

### Apple Calendar

On your Apple device, click the calendar icon.
At the bottom under "Add Calendar", add a subscription calendar and paste the URL.

## ⌨️ Development

For local development, create a `testdata/config.yml`, then run `./run.sh` to start CalUnite via the commandline (may need elevated permissions depending on your Docker setup).
