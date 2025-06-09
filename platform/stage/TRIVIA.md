# Trivia

Trivia is for the generic information mostly for tool use and always not specific to the project. Safe to skip for experts unless forgotten.

**Connecting to a VPN server on macOS**

```sh
sudo openvpn path/to/client/profile.ovpn
sudo killall mDNSResponder{,Helper}
```

**Running a Terraform config**

```sh
(cd provision/? && terraform apply --auto-approve)
(cd provision/? && terraform destroy --auto-approve)
```

**Running deployment playbooks**

```sh
ansible-playbook --forks="20" playbook.yml
ansible-playbook --forks="20" --limit="?" --tags="redeploy" playbook.yml
```