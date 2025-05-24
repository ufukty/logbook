digitalocean = {
  activated_regions = {
    vpc = [
      "sfo2",
      "sfo3",
      "tor1",
      "nyc1",
      "nyc3",
      "lon1",
      "ams3",
      "fra1",
      "blr1",
      "sgp1"
    ]
    vpn = [
      # "sfo2",
      # "sfo3",
      # "tor1",
      # "nyc1",
      # "nyc3",
      # "lon1",
      # "ams3",
      "fra1",
      # "blr1",
      # "sgp1"
    ]
  }
  config = {
    vpc = {
      sfo2 = { range = "10.140.0.0/16" }
      sfo3 = { range = "10.141.0.0/16" }
      tor1 = { range = "10.142.0.0/16" }
      nyc1 = { range = "10.143.0.0/16" }
      nyc3 = { range = "10.144.0.0/16" }
      lon1 = { range = "10.145.0.0/16" }
      ams3 = { range = "10.146.0.0/16" }
      fra1 = { range = "10.147.0.0/16" }
      blr1 = { range = "10.148.0.0/16" }
      sgp1 = { range = "10.149.0.0/16" }
    }
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
