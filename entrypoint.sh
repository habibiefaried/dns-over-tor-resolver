#!/bin/sh
cat /app/hosts >> /etc/hosts
echo "Running the DNS server..."
/app/main