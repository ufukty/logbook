terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.27.1"
    }
  }
}

variable "DO_SSH_FINGERPRINT" { type = string }

data "digitalocean_vpc" "sfo2" { name = "logbook-sfo2" }
data "digitalocean_vpc" "sfo3" { name = "logbook-sfo3" }
data "digitalocean_vpc" "tor1" { name = "logbook-tor1" }
data "digitalocean_vpc" "nyc1" { name = "logbook-nyc1" }
data "digitalocean_vpc" "nyc3" { name = "logbook-nyc3" }
data "digitalocean_vpc" "lon1" { name = "logbook-lon1" }
data "digitalocean_vpc" "ams3" { name = "logbook-ams3" }
data "digitalocean_vpc" "fra1" { name = "logbook-fra1" }
data "digitalocean_vpc" "blr1" { name = "logbook-blr1" }
data "digitalocean_vpc" "sgp1" { name = "logbook-sgp1" }

data "digitalocean_droplet_snapshot" "sfo2" {
  name_regex  = "^logbook_builder_vpn_.*"
  region      = "sfo2"
  most_recent = true
}

data "digitalocean_droplet_snapshot" "sfo3" {
  name_regex  = "^logbook_builder_vpn_.*"
  region      = "sfo3"
  most_recent = true
}

# data "digitalocean_droplet_snapshot" "tor1" {
#   name_regex  = "^logbook_builder_vpn_.*"
#   region      = "tor1"
#   most_recent = true
# }

data "digitalocean_droplet_snapshot" "nyc1" {
  name_regex  = "^logbook_builder_vpn_.*"
  region      = "nyc1"
  most_recent = true
}

data "digitalocean_droplet_snapshot" "nyc3" {
  name_regex  = "^logbook_builder_vpn_.*"
  region      = "nyc3"
  most_recent = true
}

# data "digitalocean_droplet_snapshot" "lon1" {
#   name_regex  = "^logbook_builder_vpn_.*"
#   region      = "lon1"
#   most_recent = true
# }

# data "digitalocean_droplet_snapshot" "ams3" {
#   name_regex  = "^logbook_builder_vpn_.*"
#   region      = "ams3"
#   most_recent = true
# }

# data "digitalocean_droplet_snapshot" "fra1" {
#   name_regex  = "^logbook_builder_vpn_.*"
#   region      = "fra1"
#   most_recent = true
# }

# data "digitalocean_droplet_snapshot" "blr1" {
#   name_regex  = "^logbook_builder_vpn_.*"
#   region      = "blr1"
#   most_recent = true
# }

# data "digitalocean_droplet_snapshot" "sgp1" {
#   name_regex  = "^logbook_builder_vpn_.*"
#   region      = "sgp1"
#   most_recent = true
# }

resource "digitalocean_droplet" "sfo2" {
  ipv6        = true
  name        = "logbook-vpn-sfo2"
  size        = "s-1vcpu-1gb"
  image       = data.digitalocean_droplet_snapshot.sfo2.id
  region      = "sfo2"
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.DO_SSH_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.sfo2.id
  tags        = ["vpn"]
}

resource "digitalocean_droplet" "sfo3" {
  ipv6        = true
  name        = "logbook-vpn-sfo3"
  size        = "s-1vcpu-1gb"
  image       = data.digitalocean_droplet_snapshot.sfo3.id
  region      = "sfo3"
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.DO_SSH_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.sfo3.id
  tags        = ["vpn"]
}

# resource "digitalocean_droplet" "tor1" {
#   ipv6        = true
#   name        = "logbook-vpn-tor1"
#   size        = "s-1vcpu-1gb"
#   image       = data.digitalocean_droplet_snapshot.tor1.id
#   region      = "tor1"
#   backups     = false
#   monitoring  = true
#   resize_disk = false
#   ssh_keys    = [var.DO_SSH_FINGERPRINT]
#   vpc_uuid    = data.digitalocean_vpc.tor1.id
#   tags        = ["vpn"]
# }

resource "digitalocean_droplet" "nyc1" {
  ipv6        = true
  name        = "logbook-vpn-nyc1"
  size        = "s-1vcpu-1gb"
  image       = data.digitalocean_droplet_snapshot.nyc1.id
  region      = "nyc1"
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.DO_SSH_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.nyc1.id
  tags        = ["vpn"]
}

resource "digitalocean_droplet" "nyc3" {
  ipv6        = true
  name        = "logbook-vpn-nyc3"
  size        = "s-1vcpu-1gb"
  image       = data.digitalocean_droplet_snapshot.nyc3.id
  region      = "nyc3"
  backups     = false
  monitoring  = true
  resize_disk = false
  ssh_keys    = [var.DO_SSH_FINGERPRINT]
  vpc_uuid    = data.digitalocean_vpc.nyc3.id
  tags        = ["vpn"]
}

# resource "digitalocean_droplet" "lon1" {
#   ipv6        = true
#   name        = "logbook-vpn-lon1"
#   size        = "s-1vcpu-1gb"
#   image       = data.digitalocean_droplet_snapshot.lon1.id
#   region      = "lon1"
#   backups     = false
#   monitoring  = true
#   resize_disk = false
#   ssh_keys    = [var.DO_SSH_FINGERPRINT]
#   vpc_uuid    = data.digitalocean_vpc.lon1.id
#   tags        = ["vpn"]
# }

# resource "digitalocean_droplet" "ams3" {
#   ipv6        = true
#   name        = "logbook-vpn-ams3"
#   size        = "s-1vcpu-1gb"
#   image       = data.digitalocean_droplet_snapshot.ams3.id
#   region      = "ams3"
#   backups     = false
#   monitoring  = true
#   resize_disk = false
#   ssh_keys    = [var.DO_SSH_FINGERPRINT]
#   vpc_uuid    = data.digitalocean_vpc.ams3.id
#   tags        = ["vpn"]
# }

# resource "digitalocean_droplet" "fra1" {
#   ipv6        = true
#   name        = "logbook-vpn-fra1"
#   size        = "s-1vcpu-1gb"
#   image       = data.digitalocean_droplet_snapshot.fra1.id
#   region      = "fra1"
#   backups     = false
#   monitoring  = true
#   resize_disk = false
#   ssh_keys    = [var.DO_SSH_FINGERPRINT]
#   vpc_uuid    = data.digitalocean_vpc.fra1.id
#   tags        = ["vpn"]
# }

# resource "digitalocean_droplet" "blr1" {
#   ipv6        = true
#   name        = "logbook-vpn-blr1"
#   size        = "s-1vcpu-1gb"
#   image       = data.digitalocean_droplet_snapshot.blr1.id
#   region      = "blr1"
#   backups     = false
#   monitoring  = true
#   resize_disk = false
#   ssh_keys    = [var.DO_SSH_FINGERPRINT]
#   vpc_uuid    = data.digitalocean_vpc.blr1.id
#   tags        = ["vpn"]
# }

# resource "digitalocean_droplet" "sgp1" {
#   ipv6        = true
#   name        = "logbook-vpn-sgp1"
#   size        = "s-1vcpu-1gb"
#   image       = data.digitalocean_droplet_snapshot.sgp1.id
#   region      = "sgp1"
#   backups     = false
#   monitoring  = true
#   resize_disk = false
#   ssh_keys    = [var.DO_SSH_FINGERPRINT]
#   vpc_uuid    = data.digitalocean_vpc.sgp1.id
#   tags        = ["vpn"]
# }
