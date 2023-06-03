# Nginx

**Extra code (eg maint code), goes into _meta folder**

Edit ``_meta.yaml`` with ``nginx_cert_path`` and ``origin_cert_path``

```yaml
nginx_certs: /certs
nginx_definitions: data/nginx
```

## Each definition is formated as:

```yaml
servers:
  - names: 
    - api
    - spider
    comment: "Popplio"
    locations:
      - path: / # Path of the location
        proxy: "http://127.0.0.1:8081"
        opts: # Options for the location, note that this cannot be a simple map as it needs to be ordered
          - "proxy_set_header Host $host" # As an example
```

Cert and key file paths are derived by removing suffix and adding cert- and key- respectively as well as adding a PEM extension, e.g:

E.g.: infinitybots-gg.yaml => cert-infinitybots-gg.pem and key-infinitybots-gg.pem