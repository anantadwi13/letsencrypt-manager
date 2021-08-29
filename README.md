# Let's Encrypt Manager

Manage Let's Encrypt certificates via REST API

## Requirements

- Have a public IP (Port 80 is accessible from the internet)
- Domains that want to obtain a certificate need to be pointed to your public IP (DNS)
- Run the container to your public server

## Installation

```shell
docker run -it --name letsencrypt-manager \
      -p 80:80 \        # Let's Encrypt challenge port
      -p 5555:5555 \    # API port
      -v $(pwd)/temp/letsencrypt:/etc/letsencrypt \
      anantadwi13/letsencrypt-manager
```

## Usage

After running container, open API Specification on `http://{host}:5555/docs`

## Known Issue

- Unable to create wildcard certificates since Certbot needs DNS challenge to verify ownership.