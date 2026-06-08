#!/bin/bash
set -e

# Update packages and install prerequisites
apt-get update -y
apt-get install -y ca-certificates curl gnupg git

# Install Docker
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Ensure Ubuntu user can run docker without sudo
usermod -aG docker ubuntu

# Clone the HiveFS repository
# Note: In a real company, we'd pull a specific release tag or use private deployment keys.
cd /home/ubuntu
git clone https://github.com/adarshkshitij/Hivefs.git
cd Hivefs

# Fix permissions
chown -R ubuntu:ubuntu /home/ubuntu/Hivefs

# Start the cluster using Docker Compose
docker compose up -d --build
