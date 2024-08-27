#!/bin/sh

# Maximum number of attempts to connect to the database
MAX_RETRIES=10
RETRY_INTERVAL=5

# Check if MySQL is ready
echo "Checking MySQL connection..."
for i in $(seq 1 $MAX_RETRIES); do
    if mysqladmin ping -h"$MYSQL_HOST" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --silent; then
        echo "MySQL is up and running!"
        break
    fi

    if [ "$i" -eq "$MAX_RETRIES" ]; then
        echo "MySQL did not become available after $MAX_RETRIES attempts. Exiting."
        exit 1
    fi

    echo "MySQL is not available yet. Retrying in $RETRY_INTERVAL seconds... (Attempt $i/$MAX_RETRIES)"
    sleep $RETRY_INTERVAL
done

# Run migrations
echo "Running database migrations..."
go run db/migration.go -exec up

# Start the application
echo "Starting application..."
air
