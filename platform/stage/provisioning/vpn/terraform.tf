terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.27.1"
    }
  }
}

variable "digitalocean" {
  type = object({
    activated_regions = object({
      vpc = set(string)
      vpn = set(string)
    })
    config = object({
      vpc = object({
        sfo2 = object({ range = string })
        sfo3 = object({ range = string })
        tor1 = object({ range = string })
        nyc1 = object({ range = string })
        nyc3 = object({ range = string })
        lon1 = object({ range = string })
        ams3 = object({ range = string })
        fra1 = object({ range = string })
        blr1 = object({ range = string })
        sgp1 = object({ range = string })
      })
      vpn = object({
        sfo2 = object({ subnet_address = string })
        sfo3 = object({ subnet_address = string })
        tor1 = object({ subnet_address = string })
        nyc1 = object({ subnet_address = string })
        nyc3 = object({ subnet_address = string })
        lon1 = object({ subnet_address = string })
        ams3 = object({ subnet_address = string })
        fra1 = object({ subnet_address = string })
        blr1 = object({ subnet_address = string })
        sgp1 = object({ subnet_address = string })
      })
    })
  })
}
variable "DO_SSH_FINGERPRINT" { type = string }

locals {
  sudo_user           = "olwgtzjzhnvexhpr"
  openvpn_client_name = "provisioner"
}

# MARK: Data gathering

data "digitalocean_droplet_snapshot" "vpn" {
  for_each = var.digitalocean.activated_regions.vpn

  name_regex  = "^logbook_builder_vpn_.*"
  region      = each.value
  most_recent = true
}

data "digitalocean_vpc" "vpc" {
  for_each = var.digitalocean.activated_regions.vpn

  name = "logbook-${each.value}"
}

# MARK: Resource creation

resource "digitalocean_droplet" "vpn-server" {
  for_each = var.digitalocean.activated_regions.vpn

  ipv6        = true
  name        = "${each.value}-vpn"
  size        = "s-1vcpu-1gb"
  image       = data.digitalocean_droplet_snapshot.vpn[each.value].id
  region      = each.value
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.DO_SSH_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.vpc[each.value].id
  tags        = ["vpn"]

  connection {
    host    = self.ipv4_address
    user    = local.sudo_user
    type    = "ssh"
    agent   = true
    timeout = "2m"
  }

  provisioner "file" {
    source      = "${path.module}/provisioner-files"
    destination = "/home/${local.sudo_user}"
  }

  provisioner "remote-exec" {
    inline = [
      <<EOF
        cd ~/provisioner-files && \
                     USER_ACCOUNT_NAME='${local.sudo_user}' \
                           SERVER_NAME='logbook-do-${each.value}-vpn' \
                             PUBLIC_IP='${self.ipv4_address}' \
                            PRIVATE_IP='${self.ipv4_address_private}' \
                OPENVPN_SUBNET_ADDRESS='${var.digitalocean.config.vpn[each.value].subnet_address}' \
                   OPENVPN_SUBNET_MASK='255.255.255.0' \
             PUBLIC_ETHERNET_INTERFACE='eth0' \
            PRIVATE_ETHERNET_INTERFACE='eth1' \
            sudo --preserve-env bash deployment.sh 
      EOF   
    ]
  }

  // Clean up in server & last changes
  provisioner "remote-exec" {
    on_failure = continue
    inline = [<<EOF
        rm -rf ~/artifacts ~/provisioner-files
        sudo bash -c "\
            systemctl restart systemd-journald;\
            systemctl restart iptables-activation;\
            sed -E -in-place \"s;${local.sudo_user}(.*)NOPASSWD:(.*);${local.sudo_user} \1 \2;\" /etc/sudoers;\
        "
    EOF
    ]
  }
}

resource "local_file" "ssh-config" {
  content = templatefile(
    "${path.module}/templates/ssh.conf.tftpl",
    {
      providers = {
        digitalocean = digitalocean_droplet.vpn-server
      }
    }
  )
  filename = abspath("${path.module}/../../artifacts/ssh.conf.d/0.vpn.conf")
}
