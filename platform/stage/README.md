# Platform

## Order of steps

Follow the order:

- image
  - base, internal
  - vpn, database, application etc.
- provision
  - vpc
  - vpn
  - databases
  - services
- artifacts
  - ssh config
  - instance list
- deployment
  - secrets
  - instance list

## Instructions

### Open shell

Run `autosource` in shell in the most specific folder.

### Rotate secrets

Use `scripts/secrets/*.sh` to init secrets and rotate regularly.

**CA hierarchy**

```
Root CA
├── Web CA
│   ├── API Gateway (server)
│   └── ...
└── VPN CA
    ├── DO-NYC2 (server)
    ├── DO-SFO3 (server)
    └── ...
```

### Building images

Run `bash build.sh` in the folder of each image. Building images based on `internal` requires active VPN connection into the build region.

**Image hierarchy**

```
base                # user, utilities, fail2ban, basic security
├── vpn             # openvpn, easy-rsa, ovpn-auth
└── internal        # firewall, accessible with internal network
    ├── gateway     # allows :8080 on firewall
    ├── database    # postgres, tunnel with application
    └── application # systemd service, logging, certs, database tunnel
```

### Creating cloud resources

Run terraform apply on each folder in `provisioning`.

### Artifacts

Run `bash scripts/...` to build ssh config and instance list.

### Deployment

Run `ansible playbook.yml` in deployment folder.
