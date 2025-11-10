${{ content_synopsis }} Run [11notes/distroless:lego](https://github.com/11notes/docker-distroless/blob/master/lego.dockerfile) [rootless](https://github.com/11notes/RTFM/blob/main/linux/container/image/rootless.md) and [distroless](https://github.com/11notes/RTFM/blob/main/linux/container/image/distroless.md) on a schedule to automatically renew all your certificates from a single yml config.

${{ title_volumes }}
* **${{ json_root }}/etc** - Directory of your Let's Encrypt accounts
* **${{ json_root }}/var** - Directory of your Let's Encrypt certificates

${{ content_compose }}

${{ content_defaults }}

${{ content_environment }}
| `LEGO_CONFIG` | Your config for all your domains, either as a file or inline | |

${{ content_source }}

${{ content_parent }}

${{ content_built }}

${{ content_tips }}