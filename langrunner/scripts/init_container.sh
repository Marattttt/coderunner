#!/bin/bash

# Create system user 'gouser' with home directory but no login shell
useradd -r \
    -m \
    -d /home/gouser \
    -s /usr/sbin/nologin \
    gouser

# # Verify the user was created
# id gouser || {echo "Failed to create user" && exit 1}
#
# # Display user information
# grep gouser /etc/passwd
