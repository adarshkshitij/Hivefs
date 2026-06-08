output "hivefs_server_public_ip" {
  description = "The Public IP address of the HiveFS Server"
  value       = aws_instance.hivefs_node.public_ip
}

output "grafana_dashboard_url" {
  description = "URL to access the Grafana Dashboard"
  value       = "http://${aws_instance.hivefs_node.public_ip}:3000"
}

output "ssh_command" {
  description = "Command to SSH into the server"
  value       = "ssh -i <your-key.pem> ubuntu@${aws_instance.hivefs_node.public_ip}"
}
