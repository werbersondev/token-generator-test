#!/bin/bash

set -e

# SonarQube credentials
SONARQUBE_URL="http://localhost:9000"
SONARQUBE_USER="admin"
SONARQUBE_PASSWORD="admin"
NEW_SONARQUBE_PASSWORD="new_admin_password"

# Wait for SonarQube to be up
echo "Waiting for SonarQube to be up..."
until $(curl --silent --fail http://localhost:9000/api/system/status | grep '"status":"UP"' > /dev/null); do
  printf '.'
  sleep 5
done
echo "SonarQube is up."

# Change admin password
echo "Changing admin password..."
curl -s -u $SONARQUBE_USER:$SONARQUBE_PASSWORD -X POST \
  "$SONARQUBE_URL/api/users/change_password" \
  -d "login=admin&previousPassword=$SONARQUBE_PASSWORD&password=$NEW_SONARQUBE_PASSWORD"

# Update credentials with new password
SONARQUBE_PASSWORD=$NEW_SONARQUBE_PASSWORD

echo "Admin password changed."

# Generate a user token
TOKEN_NAME="my_user_token"
echo "Generating user token..."
user_token=$(curl -s -u $SONARQUBE_USER:$SONARQUBE_PASSWORD -X POST \
  "$SONARQUBE_URL/api/user_tokens/generate" \
  -d "name=$TOKEN_NAME" | jq -r '.token')

if [ -z "$user_token" ]; then
  echo "Failed to generate user token."
  exit 1
fi

echo "Generated user token: $user_token"

# Project settings
PROJECT_KEYS=("project_key_1" "project_key_2" "project_key_3" "project_key_4")
PROJECT_NAMES=("Project 1" "Project 2" "Project 3" "Project 4")

for i in "${!PROJECT_KEYS[@]}"; do
  PROJECT_KEY=${PROJECT_KEYS[$i]}
  PROJECT_NAME=${PROJECT_NAMES[$i]}

  # Create a new project
  curl -s -u $SONARQUBE_USER:$SONARQUBE_PASSWORD -X POST \
    "$SONARQUBE_URL/api/projects/create" \
    -d "name=$PROJECT_NAME&project=$PROJECT_KEY"
done

echo "SonarQube setup completed."
echo "Project Identifiers:"
for PROJECT_KEY in "${PROJECT_KEYS[@]}"; do
  echo "$PROJECT_KEY"
done
