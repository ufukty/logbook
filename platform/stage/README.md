# Platform

## Order of steps

Follow the order:

- build
  - time consuming, low-frequency tasks
- provision
  - vpc
  - vpn
  - databases
  - services
- artifacts
  - ssh config
  - instance list
- deploy
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
├── VPN Users CA
│   └── ...
└── VPN CA
    ├── DO-NYC2 (server)
    ├── DO-SFO3 (server)
    └── ...
```

### Building images

Run `bash build.sh` in the folder of each image under `build` folder. Building images based on `internal` requires active VPN connection into the build region.

**Image hierarchy**

```
base                # user, utilities, fail2ban, basic security
├── vpn             # openvpn, easy-rsa, ovpn-auth
└── internal        # firewall, accessible with internal network
    ├── gateway     # allows :8080 on firewall
    ├── database    # postgres, tunnel with application
    └── application # systemd service, logging, certs, database tunnel
```

### Provisioning

Run terraform apply on each folder in `provisioning`.

### Artifacts

Run `bash scripts/...` to build ssh config and instance list. They rely on information generated during provisioning.

### Deployment

Run the shell scripts under `deployment` folder, from the `$STAGE` directory.

| Script      | What it does                                   |
| ----------- | ---------------------------------------------- |
| local.sh    | Adds Root CA. Refreshes `known_hosts` entries. |
| all.sh      | Uploads Root CA.                               |
| finalize.sh | Removes passwordless sudo. Reloads journald.   |

Deployment scripts are generally safe to run repeatedly until finalization.
