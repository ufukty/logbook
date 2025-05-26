client
proto tcp-client
remote {{PUBLIC_IP}} 443
dev tun

resolv-retry infinite
nobind
persist-key
persist-tun

remote-cert-tls server
verify-x509-name {{EASYRSA_SERVER_NAME}} name

auth SHA256
auth-nocache

cipher AES-128-GCM
tls-client
tls-version-min 1.2
tls-cipher TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256

ignore-unknown-option block-outside-dns
setenv opt block-outside-dns # Prevent Windows 10 DNS leak

auth-user-pass

verb 4

<ca>
{{EASYRSA_CA_CERT_CONTENT}}
</ca>

<cert>
{{EASYRSA_CLIENT_CERT_CONTENT}}
</cert>

<key>
{{EASYRSA_CLIENT_KEY_CONTENT}}
</key>

key-direction 1
<tls-crypt>
{{TLS_SIG_KEY_CONTENT}}
</tls-crypt>
