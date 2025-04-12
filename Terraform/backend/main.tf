terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
}

provider "aws" {
  region     = "us-east-2"
  access_key = var.access_key
  secret_key = var.secret_key
}

resource "aws_security_group" "web" {

  name        = "security_group_terraform"
  description = "Security group para permitir SSH y HTTP"

  ingress {
    description = "SSH desde el mismo SG"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    self        = true
  }

  ingress {
    description = "SSH desde mi IP"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["181.51.33.28/32"]
  }

  ingress {
    description = "Puerto 9950 abierto a todo"
    from_port   = 9950
    to_port     = 9950
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Permitir todo de salida"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "Backend-terraform" {
  ami           = "ami-0100e595e1cc1ff7f" # linux aws
  instance_type = "t2.micro"
  subnet_id     = "subnet-0febf33cc7c460afe"
  tags = {
    Name = "Backend-terraform"
  }
}