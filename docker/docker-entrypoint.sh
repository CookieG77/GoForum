#!/bin/sh
set -e

# Define the list of required files
REQUIRED_FILES="key.pem cert.pem .db .env"

# Check each file in the current working directory
for file in $REQUIRED_FILES; do
    if [ ! -f "$file" ]; then
        echo "Error: Required file '$file' not found!"
        exit 1
    fi
done

# If all files exist, execute the main binary
echo "All required files are present. Starting the app..."
exec ./main
