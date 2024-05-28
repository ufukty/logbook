terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.27.1"
    }
  }
}

variable "project_prefix" { type = string }
variable "SSH_KEY_FINGERPRINT" { type = string }

data "digitalocean_droplet_snapshot" "application" {
  name_regex  = "^build_application_.*"
  region      = "fra1"
  most_recent = true
}

data "digitalocean_droplet_snapshot" "gateway" {
  name_regex  = "^build_gateway_.*"
  region      = "fra1"
  most_recent = true
}

data "digitalocean_vpc" "fra1" {
  name = "${var.project_prefix}-fra1"
}

resource "digitalocean_droplet" "account" {
  count = 1

  image  = data.digitalocean_droplet_snapshot.application.id
  name   = "fra1-account-${count.index}"
  region = "fra1"
  size   = "s-1vcpu-1gb"
  tags   = ["account"]

  ipv6        = true
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.SSH_KEY_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.fra1.id
}

resource "digitalocean_droplet" "objectives" {
  count = 1

  image  = data.digitalocean_droplet_snapshot.application.id
  name   = "fra1-objectives-${count.index}"
  region = "fra1"
  size   = "s-1vcpu-1gb"
  tags   = ["objectives"]

  ipv6        = true
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.SSH_KEY_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.fra1.id
}

resource "digitalocean_droplet" "gateway" {
  count = 1

  image  = data.digitalocean_droplet_snapshot.gateway.id
  name   = "fra1-gateway-${count.index}"
  region = "fra1"
  size   = "s-1vcpu-1gb"
  tags   = ["gateway"]

  ipv6        = true
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.SSH_KEY_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.fra1.id
}

resource "local_file" "inventory" {
  content = templatefile(
    "${path.module}/templates/inventory.cfg.tftpl",
    {
      providers = {
        digitalocean = {
          fra1 = {
            vpc = data.digitalocean_vpc.fra1
            services = {
              gateway    = digitalocean_droplet.gateway
              objectives = digitalocean_droplet.objectives
              account    = digitalocean_droplet.account
            }
          }
        }
      }
    }
  )
  filename = abspath("${path.module}/../../artifacts/deployment/inventory.cfg")
}

resource "local_file" "ssh-config" {
  content = templatefile(
    "${path.module}/templates/ssh.conf.tftpl",
    {
      providers = {
        digitalocean = {
          fra1 = {
            vpc = data.digitalocean_vpc.fra1
            services = {
              gateway    = digitalocean_droplet.gateway
              objectives = digitalocean_droplet.objectives
              account    = digitalocean_droplet.account
            }
          }
        }
      }
    }
  )
  filename = abspath("${path.module}/../../artifacts/ssh.conf.d/0.application.conf")
}

resource "local_file" "service_discovery" {
  content = templatefile(
    "${path.module}/templates/service_discovery.json.tftpl",
    {
      content = jsonencode({
        digitalocean = {
          fra1 = {
            vpc = data.digitalocean_vpc.fra1
            services = {
              gateway    = digitalocean_droplet.gateway
              objectives = digitalocean_droplet.objectives
              account    = digitalocean_droplet.account
            }
          }
        }
      })
    }
  )
  filename = abspath("${path.module}/../../artifacts/deployment/service_discovery.json")
}
