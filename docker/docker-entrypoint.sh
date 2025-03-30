#!/bin/sh

# Allow changing UID/GID at runtime
if [ "$(id -u)" -ne "$PUID" ] || [ "$(id -g)" -ne "$PGID" ]; then
  echo "Updating UID:GID to $PUID:$PGID"

  # Update group ID
  groupmod -g "$PGID" app

  # Update user ID
  usermod -u "$PUID" -g "$PGID" app

  # Update ownership of app files
  chown -R app:app /app
fi

# Switch to the app user
exec su-exec app "$@"
