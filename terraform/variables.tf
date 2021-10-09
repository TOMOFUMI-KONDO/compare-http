variable "region" {
  type = string
}

variable "project" {
  type = string
}

variable "my_global_ip" {
  type = string
}

variable "az" {
  type = map(string)
  default = {
    a = "ap-northeast-1a"
    c = "ap-northeast-1c"
    d = "ap-northeast-1d"
  }
}

variable "default_gateway_cidr" {
  type = map(string)
  default = {
    ipv4 = "0.0.0.0/0"
    ipv6 = "::/0"
  }
}

variable "vpc_cidr" {
  type = string
}

variable "subnet_cidr" {
  type = map(map(string))
  default = {
    public = {
      a = ""
      c = ""
      d = ""
    }
    private = {
      a = ""
      c = ""
      d = ""
    }
  }
}
