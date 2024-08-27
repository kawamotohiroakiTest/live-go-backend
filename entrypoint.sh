#!/bin/sh

# Run migrations
echo "Running database migrations..."
go run db/migration.go -exec up

# Start the application
echo "Starting application..."
air
