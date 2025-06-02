![banner](https://github.com/11notes/defaults/blob/main/static/img/banner.png?raw=true)

# LEGO
[<img src="https://img.shields.io/badge/github-source-blue?logo=github&color=040308">](https://github.com/11notes/docker-LEGO)![5px](https://github.com/11notes/defaults/blob/main/static/img/transparent5x2px.png?raw=true)![size](https://img.shields.io/docker/image-size/11notes/lego/1.0.0?color=0eb305)![5px](https://github.com/11notes/defaults/blob/main/static/img/transparent5x2px.png?raw=true)![version](https://img.shields.io/docker/v/11notes/lego/1.0.0?color=eb7a09)![5px](https://github.com/11notes/defaults/blob/main/static/img/transparent5x2px.png?raw=true)![pulls](https://img.shields.io/docker/pulls/11notes/lego?color=2b75d6)![5px](https://github.com/11notes/defaults/blob/main/static/img/transparent5x2px.png?raw=true)[<img src="https://img.shields.io/github/issues/11notes/docker-LEGO?color=7842f5">](https://github.com/11notes/docker-LEGO/issues)![5px](https://github.com/11notes/defaults/blob/main/static/img/transparent5x2px.png?raw=true)![swiss_made](https://img.shields.io/badge/Swiss_Made-FFFFFF?labelColor=FF0000&logo=data:image/svg%2bxml;base64,PHN2ZyB2ZXJzaW9uPSIxIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgdmlld0JveD0iMCAwIDMyIDMyIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Im0wIDBoMzJ2MzJoLTMyeiIgZmlsbD0iI2YwMCIvPjxwYXRoIGQ9Im0xMyA2aDZ2N2g3djZoLTd2N2gtNnYtN2gtN3YtNmg3eiIgZmlsbD0iI2ZmZiIvPjwvc3ZnPg==)

Run lego on a schedule rootless, distroless and secure by default!

# SYNOPSIS 📖
**What can I do with this?** Run [11notes/distroless:lego](https://github.com/11notes/docker-distroless/blob/master/lego.dockerfile) on a schedule to automatically renew all your certificates from a single yml config. 

# VOLUMES 📁
* **/lego/etc** - Directory of your Let's Encrypt accounts
* **/lego/var** - Directory of your Let's Encrypt certificates

# COMPOSE ✂️
```yaml
name: "letsencrypt"
services:
  lego:
    image: "11notes/lego:1.0.0"
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
      - "etc:/lego/etc" 
      - "var:/lego/var"
    networks:
      frontend:
    restart: "always"
volumes:
  etc:
  var:
networks:
  frontend:
```

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

* [1.0.0](https://hub.docker.com/r/11notes/lego/tags?name=1.0.0)

### There is no latest tag, what am I supposed to do about updates?
It is of my opinion that the ```:latest``` tag is super dangerous. Many times, I’ve introduced **breaking** changes to my images. This would have messed up everything for some people. If you don’t want to change the tag to the latest [semver](https://semver.org/), simply use the short versions of [semver](https://semver.org/). Instead of using ```:1.0.0``` you can use ```:1``` or ```:1.0```. Since on each new version these tags are updated to the latest version of the software, using them is identical to using ```:latest``` but at least fixed to a major or minor version.

If you still insist on having the bleeding edge release of this app, simply use the ```:rolling``` tag, but be warned! You will get the latest version of the app instantly, regardless of breaking changes or security issues or what so ever. You do this at your own risk!

# REGISTRIES ☁️
```
docker pull 11notes/lego:1.0.0
docker pull ghcr.io/11notes/lego:1.0.0
docker pull quay.io/11notes/lego:1.0.0
```

# SOURCE 💾
* [11notes/lego](https://github.com/11notes/docker-LEGO)

# PARENT IMAGE 🏛️
> [!IMPORTANT]
>This image is not based on another image but uses [scratch](https://hub.docker.com/_/scratch) as the starting layer.
>The image consists of the following distroless layers that were added:
>* [11notes/distroless](https://github.com/11notes/docker-distroless/blob/master/arch.dockerfile) - contains users, timezones and Root CA certificates
>* 11notes/distroless:lego

# BUILT WITH 🧰
* [lego](https://github.com/11notes/docker-lego)

# GENERAL TIPS 📌
> [!TIP]
>* Use a reverse proxy like Traefik, Nginx, HAproxy to terminate TLS and to protect your endpoints
>* Use Let’s Encrypt DNS-01 challenge to obtain valid SSL certificates for your services

# ElevenNotes™️
This image is provided to you at your own risk. Always make backups before updating an image to a different version. Check the [releases](https://github.com/11notes/docker-lego/releases) for breaking changes. If you have any problems with using this image simply raise an [issue](https://github.com/11notes/docker-lego/issues), thanks. If you have a question or inputs please create a new [discussion](https://github.com/11notes/docker-lego/discussions) instead of an issue. You can find all my other repositories on [github](https://github.com/11notes?tab=repositories).

*created 03.06.2025, 01:16:54 (CET)*