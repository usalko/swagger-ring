# swagger-merge-docs-traefik-plugin

A Middleware plugin for Traefik allow merge multiply swagger doc endpoints to a single one.
Perhaps you'll find it usable for multiply microservices which served to one traefik balancer.

## Use case

swagger-merge-docs-config.yaml

```yaml
http:
  routers:
    docs-router:
      rule: PathPrefix(`/api/v1/docs`)
      service: docs-service
      middlewares:
        - swagger

  services:
    docs-service:
      loadBalancer:
        servers:
          - url: http://whoami
  
  middlewares:
    swagger:
      plugin:
        swagger-merge-docs:
          path: /api/v1/docs
          docs:
            - path: http://service1:3000/swagger.yaml
            - path: http://service2:3000/swagger.yaml
```

docker-compose.yaml

```yaml
services:
  traefik:
    image: traefik:latest
    restart: unless-stopped
    command:
      # configuration folder
      - --providers.file.directory=/config
      - --providers.file.watch=true
      # plugin
      - --experimental.plugins.swagger-merge-docs.modulename=github.com/usalko/swagger-merge-docs
      - --experimental.plugins.swagger-merge-docs.version=v0.1.1
    volumes:
      - ./swagger-merge-docs-config.yaml:/config/swagger-merge-docs-config.yaml

  whoami:
    restart: unless-stopped
    image: traefik/whoami
```
