terraform {
  backend "local" {
    path = ".terraform/terraform.tfstate"
  }
}

provider "aws" {
  region = var.region
}
