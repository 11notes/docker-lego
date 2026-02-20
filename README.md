![banner](https://raw.githubusercontent.com/11notes/static/refs/heads/main/img/banner/README.png)

# LEGO
![size](https://img.shields.io/badge/image_size-${{ image_size }}-green?color=%2338ad2d)![5px](https://raw.githubusercontent.com/11notes/static/refs/heads/main/img/markdown/transparent5x2px.png)![pulls](https://img.shields.io/docker/pulls/11notes/lego?color=2b75d6)![5px](https://raw.githubusercontent.com/11notes/static/refs/heads/main/img/markdown/transparent5x2px.png)[<img src="https://img.shields.io/github/issues/11notes/docker-lego?color=7842f5">](https://github.com/11notes/docker-lego/issues)![5px](https://raw.githubusercontent.com/11notes/static/refs/heads/main/img/markdown/transparent5x2px.png)![swiss_made](https://img.shields.io/badge/Swiss_Made-FFFFFF?labelColor=FF0000&logo=data:image/svg%2bxml;base64,PHN2ZyB2ZXJzaW9uPSIxIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgdmlld0JveD0iMCAwIDMyIDMyIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgogIDxyZWN0IHdpZHRoPSIzMiIgaGVpZ2h0PSIzMiIgZmlsbD0idHJhbnNwYXJlbnQiLz4KICA8cGF0aCBkPSJtMTMgNmg2djdoN3Y2aC03djdoLTZ2LTdoLTd2LTZoN3oiIGZpbGw9IiNmZmYiLz4KPC9zdmc+)

Run lego on a schedule rootless, distroless and secure by default!

# INTRODUCTION 📢

Let's Encrypt client and ACME library written in Go.

# SYNOPSIS 📖
**What can I do with this?** Run [11notes/distroless:lego](https://github.com/11notes/docker-distroless/blob/master/lego.dockerfile) [rootless](https://github.com/11notes/RTFM/blob/main/linux/container/image/rootless.md) and [distroless](https://github.com/11notes/RTFM/blob/main/linux/container/image/distroless.md) on a schedule to automatically renew all your certificates from a single yml config.

# VOLUMES 📁
* **/lego/etc** - Directory of your Let's Encrypt accounts
* **/lego/var** - Directory of your Let's Encrypt certificates

# COMPOSE ✂️
```yaml
name: "letsencrypt"
services:
  lego:
    image: "11notes/lego:4.32.0"
    dns:
      - "8.8.8.8"
      - "9.9.9.9"
    read_only: true
    environment:
      TZ: "Europe/Zurich"
      LEGO_CONFIG: |-
        domains:
          - name: "domain.com"
            fqdns:
              - "*.domain.com"
              - "domain.com"
            commands:
              - "--dns"
              - "rfc2136"

          - name: "porkbun.com"
            fqdns:
              - "*.porkbun.com"
              - "porkbun.com"
            commands:
              - "--dns"
              - "porkbun"
        global:
          LEGO_EMAIL: "info@domain.com"
          RFC2136_NAMESERVER: "ns.domain.com"
          RFC2136_TSIG_ALGORITHM: "hmac-sha512"
          RFC2136_TSIG_KEY: "lego"
          RFC2136_TSIG_SECRET: ${RFC2136_TSIG_SECRET}
          PORKBUN_SECRET_API_KEY: ${PORKBUN_SECRET_API_KEY}
          PORKBUN_API_KEY: ${PORKBUN_API_KEY}
    volumes:
      - "lego.etc:/lego/etc" 
      - "lego.var:/lego/var"
    networks:
      frontend:
    restart: "always"

volumes:
  lego.etc:
  lego.var:

networks:
  frontend:
```
To find out how you can change the default UID/GID of this container image, consult the [RTFM](https://github.com/11notes/RTFM/blob/main/linux/container/image/11notes/how-to.changeUIDGID.md#change-uidgid-the-correct-way).

# DEFAULT SETTINGS 🗃️
| Parameter | Value | Description |
| --- | --- | --- |
| `user` | docker | user name |
| `uid` | 1000 | [user identifier](https://en.wikipedia.org/wiki/User_identifier) |
| `gid` | 1000 | [group identifier](https://en.wikipedia.org/wiki/Group_identifier) |
| `home` | /lego | home directory of user docker |

# ENVIRONMENT 📝
| Parameter | Value | Default |
| --- | --- | --- |
| `TZ` | [Time Zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) | |
| `DEBUG` | Will activate debug option for container image and app (if available) | |
| `LEGO_CONFIG` | Your config for all your domains, either as a file or inline | |

# MAIN TAGS 🏷️
These are the main tags for the image. There is also a tag for each commit and its shorthand sha256 value.

* [4.32.0](https://hub.docker.com/r/11notes/lego/tags?name=4.32.0)

### There is no latest tag, what am I supposed to do about updates?
It is my opinion that the ```:latest``` tag is a bad habbit and should not be used at all. Many developers introduce **breaking changes** in new releases. This would messed up everything for people who use ```:latest```. If you don’t want to change the tag to the latest [semver](https://semver.org/), simply use the short versions of [semver](https://semver.org/). Instead of using ```:4.32.0``` you can use ```:4``` or ```:4.32```. Since on each new version these tags are updated to the latest version of the software, using them is identical to using ```:latest``` but at least fixed to a major or minor version. Which in theory should not introduce breaking changes.

If you still insist on having the bleeding edge release of this app, simply use the ```:rolling``` tag, but be warned! You will get the latest version of the app instantly, regardless of breaking changes or security issues or what so ever. You do this at your own risk!

# REGISTRIES ☁️
```
docker pull 11notes/lego:4.32.0
docker pull ghcr.io/11notes/lego:4.32.0
docker pull quay.io/11notes/lego:4.32.0
```

# SOURCE 💾
* [11notes/lego](https://github.com/11notes/docker-lego)

# PARENT IMAGE 🏛️
> [!IMPORTANT]
>This image is not based on another image but uses [scratch](https://hub.docker.com/_/scratch) as the starting layer.
>The image consists of the following distroless layers that were added:
>* [11notes/distroless](https://github.com/11notes/docker-distroless/blob/master/arch.dockerfile) - contains users, timezones and Root CA certificates, nothing else
>* 11notes/distroless:lego

# BUILT WITH 🧰
* [go-acme/lego](https://github.com/go-acme/lego)

# GENERAL TIPS 📌
> [!TIP]
>* Use a reverse proxy like Traefik, Nginx, HAproxy to terminate TLS and to protect your endpoints
>* Use Let’s Encrypt DNS-01 challenge to obtain valid SSL certificates for your services

# ElevenNotes™️
This image is provided to you at your own risk. Always make backups before updating an image to a different version. Check the [releases](https://github.com/11notes/docker-lego/releases) for breaking changes. If you have any problems with using this image simply raise an [issue](https://github.com/11notes/docker-lego/issues), thanks. If you have a question or inputs please create a new [discussion](https://github.com/11notes/docker-lego/discussions) instead of an issue. You can find all my other repositories on [github](https://github.com/11notes?tab=repositories).

*created 20.02.2026, 07:01:03 (CET)*