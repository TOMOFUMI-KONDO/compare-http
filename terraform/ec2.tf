resource "aws_instance" "server" {
  ami                    = data.aws_ami.amazon-linux-2.id
  instance_type          = "t3.nano"
  subnet_id              = aws_subnet.public_a.id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  key_name               = var.project

  tags = {
    Name    = "${var.project}-ec2-a"
    Project = var.project
  }
}

resource "aws_instance" "client" {
  ami                    = data.aws_ami.amazon-linux-2.id
  instance_type          = "t3.nano"
  subnet_id              = aws_subnet.public_c.id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  key_name               = var.project

  tags = {
    Name    = "${var.project}-ec2-c"
    Project = var.project
  }
}

output "server_public_ip" {
  value = aws_instance.server.public_ip
}

output "server_private_ip" {
  value = aws_instance.server.private_ip
}

output "client_public_ip" {
  value = aws_instance.client.public_ip
}

output "client_private_ip" {
  value = aws_instance.client.private_ip
}

resource "aws_security_group" "ec2" {
  description = "For ec2 instance."
  vpc_id      = aws_vpc.main.id

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = [var.default_gateway_cidr.ipv4] #tfsec:ignore:AWS009
    ipv6_cidr_blocks = [var.default_gateway_cidr.ipv6] #tfsec:ignore:AWS009
  }

  tags = {
    Name    = "${var.project}-sg"
    Project = var.project
  }
}

resource "aws_security_group_rule" "allow_ssh" {
  description       = "Allow ssh from admin."
  security_group_id = aws_security_group.ec2.id
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = [var.my_global_ip]
}

resource "aws_security_group_rule" "allow_https_tcp" {
  description       = "Allow internal https."
  security_group_id = aws_security_group.ec2.id
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = [var.vpc_cidr]
}

resource "aws_security_group_rule" "allow_https_udp" {
  description       = "Allow internal https."
  security_group_id = aws_security_group.ec2.id
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "udp"
  cidr_blocks       = [var.vpc_cidr]
}

resource "aws_security_group_rule" "allow_ping" {
  description       = "Allow internal ping."
  security_group_id = aws_security_group.ec2.id
  type              = "ingress"
  from_port         = 8
  to_port           = 0
  protocol          = "icmp"
  cidr_blocks       = [var.vpc_cidr]
}

data "aws_ami" "amazon-linux-2" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  owners = ["amazon"]
}
