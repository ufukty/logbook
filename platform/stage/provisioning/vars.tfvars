digitalocean = {
  activated_regions = {
    vpn = [
      "sfo2",
      # "sfo3",
      # "tor1",
      # "nyc1",
      "nyc3",
      # "lon1",
      # "ams3",
      # "fra1",
      # "blr1",
      # "sgp1"
    ]
  }
  config = {
    vpn = { // has to be different than VPC addresses. they will be masked with /24
      sfo2 = { subnet_address = "10.150.0.0" }
      sfo3 = { subnet_address = "10.150.1.0" }
      tor1 = { subnet_address = "10.150.2.0" }
      nyc1 = { subnet_address = "10.150.3.0" }
      nyc3 = { subnet_address = "10.150.4.0" }
      lon1 = { subnet_address = "10.150.5.0" }
      ams3 = { subnet_address = "10.150.6.0" }
      fra1 = { subnet_address = "10.150.7.0" }
      blr1 = { subnet_address = "10.150.8.0" }
      sgp1 = { subnet_address = "10.150.9.0" }
    }
  }
}
