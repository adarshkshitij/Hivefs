variable "aws_region" {
  description = "The AWS region to deploy the infrastructure in"
  default     = "us-east-1"
}

variable "instance_type" {
  description = "The EC2 instance type"
  default     = "t2.micro" # Free tier eligible
}

variable "ami_id" {
  description = "The AMI ID for the EC2 instance (Ubuntu 22.04 LTS by default in us-east-1)"
  default     = "ami-0fc5d935ebf8bc3bc" 
}
