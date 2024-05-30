2024.05.30

# Images

## Image Hierarchy

```
base                  user, utilities, fail2ban, basic security
├── vpn               openvpn, easy-rsa, ovpn-auth
└── internal          firewall, accessible with internal network
    ├── gateway       allows :8080 on firewall
    ├── database      postgres, tunnel with application
    └── application   systemd service, logging, certs, database tunnel
```
