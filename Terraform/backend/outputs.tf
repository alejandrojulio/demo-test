output "instance_id" {
  description = "ID de instancia EC2"
  value       = aws_instance.Backend-terraform.id
}

output "instance_public_ip" {
  description = "IP publica de instancia EC2"
  value       = aws_instance.Backend-terraform.public_ip
}