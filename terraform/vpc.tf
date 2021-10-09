resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true

  tags = {
    Name    = "${var.project}-vpc"
    Project = var.project
  }
}

resource "aws_subnet" "public_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.subnet_cidr.public.a
  availability_zone       = var.az.a
  map_public_ip_on_launch = true

  tags = {
    Name    = "${var.project}-subnet-public-a"
    Project = var.project
  }
}

resource "aws_subnet" "public_c" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.subnet_cidr.public.c
  availability_zone       = var.az.c
  map_public_ip_on_launch = true

  tags = {
    Name    = "${var.project}-subnet-public-c"
    Project = var.project
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name    = "${var.project}-igw"
    Project = var.project
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = var.default_gateway_cidr.ipv4
    gateway_id = aws_internet_gateway.main.id
  }

  route {
    ipv6_cidr_block = var.default_gateway_cidr.ipv6
    gateway_id      = aws_internet_gateway.main.id
  }

  tags = {
    Name    = "${var.project}-route-table-public"
    Project = var.project
  }
}

resource "aws_route_table_association" "public_a" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public_a.id
}

resource "aws_route_table_association" "public_c" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public_c.id
}
