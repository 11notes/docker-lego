${{ title_caution }}
${{ github:> [!CAUTION] }}
${{ github:> }}v5 of LEGO will use the new yml config and will not work with your existing pre v5 config, please convert your exisiting config into the new format. The format can be found [here](https://github.com/go-acme/lego/blob/main/cmd/internal/configuration/testdata/reference.yml).

${{ content_synopsis }} This image will give you a [rootless](https://github.com/11notes/RTFM/blob/main/linux/container/image/rootless.md) and [distroless](https://github.com/11notes/RTFM/blob/main/linux/container/image/distroless.md) LEGO installation to automate your cert creation.

${{ content_uvp }} Good question! Because ...

${{ github:> [!IMPORTANT] }}
${{ github:> }}* ... this image runs [rootless](https://github.com/11notes/RTFM/blob/main/linux/container/image/rootless.md) as 1000:1000
${{ github:> }}* ... this image has no shell since it is [distroless](https://github.com/11notes/RTFM/blob/main/linux/container/image/distroless.md)
${{ github:> }}* ... this image is auto updated to the latest version via CI/CD
${{ github:> }}* ... this image has a health check
${{ github:> }}* ... this image runs read-only
${{ github:> }}* ... this image is automatically scanned for CVEs before and after publishing
${{ github:> }}* ... this image is created via a secure and pinned CI/CD process
${{ github:> }}* ... this image is very small
${{ github:> }}* ... this image support [inline configs](https://github.com/11notes/RTFM/blob/master/linux/container/image/11notes/inline-config.md)

If you value security, simplicity and optimizations to the extreme, then this image might be for you.

${{ title_volumes }}
* **${{ json_root }}/etc** - Directory of your Let's Encrypt config
* **${{ json_root }}/var** - Directory of your Let's Encrypt certificates and accounts

${{ content_compose }}

${{ content_defaults }}

${{ content_environment }}
| `LEGO_CONFIG` | Will overwrite the default config with the value of this variable if set ([inline config](https://github.com/11notes/RTFM/blob/master/linux/container/image/11notes/inline-config.md)) | |

${{ content_source }}

${{ content_parent }}

${{ content_built }}

${{ content_tips }}