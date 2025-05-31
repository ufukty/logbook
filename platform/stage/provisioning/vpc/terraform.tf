terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.27.1"
    }
  }
}

resource "digitalocean_vpc" "sfo2" {
  name     = "logbook-sfo2"
  region   = "sfo2"
  ip_range = "10.140.0.0/16"
}

resource "digitalocean_vpc" "sfo3" {
  name     = "logbook-sfo3"
  region   = "sfo3"
  ip_range = "10.141.0.0/16"
}

resource "digitalocean_vpc" "tor1" {
  name     = "logbook-tor1"
  region   = "tor1"
  ip_range = "10.142.0.0/16"
}

resource "digitalocean_vpc" "nyc1" {
  name     = "logbook-nyc1"
  region   = "nyc1"
  ip_range = "10.143.0.0/16"
}

resource "digitalocean_vpc" "nyc3" {
  name     = "logbook-nyc3"
  region   = "nyc3"
  ip_range = "10.144.0.0/16"
}

resource "digitalocean_vpc" "lon1" {
  name     = "logbook-lon1"
  region   = "lon1"
  ip_range = "10.145.0.0/16"
}

resource "digitalocean_vpc" "ams3" {
  name     = "logbook-ams3"
  region   = "ams3"
  ip_range = "10.146.0.0/16"
}

resource "digitalocean_vpc" "fra1" {
  name     = "logbook-fra1"
  region   = "fra1"
  ip_range = "10.147.0.0/16"
}

resource "digitalocean_vpc" "blr1" {
  name     = "logbook-blr1"
  region   = "blr1"
  ip_range = "10.148.0.0/16"
}

resource "digitalocean_vpc" "sgp1" {
  name     = "logbook-sgp1"
  region   = "sgp1"
  ip_range = "10.149.0.0/16"
}
