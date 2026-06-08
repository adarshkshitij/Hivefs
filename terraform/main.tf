terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# 1. VPC Configuration
resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}

# 2. Security Group (Firewall Rules)
resource "aws_security_group" "hivefs_sg" {
  name        = "hivefs_security_group"
  description = "Allow inbound traffic for HiveFS, SSH, Grafana, and Prometheus"
  vpc_id      = aws_default_vpc.default.id

  # SSH
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Grafana Dashboard
  ingress {
    description = "Grafana"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HiveFS Node Ports (3001-3003)
  ingress {
    description = "HiveFS Nodes"
    from_port   = 3001
    to_port     = 3003
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Prometheus Metrics
  ingress {
    description = "Prometheus"
    from_port   = 9090
    to_port     = 9090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "HiveFS SG"
  }
}

# 3. EC2 Instance
resource "aws_instance" "hivefs_node" {
  ami           = var.ami_id
  instance_type = var.instance_type
  
  vpc_security_group_ids = [aws_security_group.hivefs_sg.id]

  # User data script to bootstrap the server
  user_data = file("${path.module}/user_data.sh")
  
  # Ensure the instance replaces user_data when changed
  user_data_replace_on_change = true

  tags = {
    Name = "HiveFS-Server"
    Project = "HiveFS"
    Environment = "Dev"
  }
}
