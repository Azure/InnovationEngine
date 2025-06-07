#!/bin/bash
# Script to remove sensitive .env files from git history

# Make sure we're in the right directory
cd /home/rgardler/projects/InnovationEngine

# Create a backup of the .env files in case something goes wrong
echo "Creating backup of .env files..."
mkdir -p /tmp/env_backup
cp -f sandbox/innovation-engine-headlamp/.env /tmp/env_backup/.env 2>/dev/null
cp -f sandbox/innovation-engine-headlamp/.env.test /tmp/env_backup/.env.test 2>/dev/null
echo "Backup created in /tmp/env_backup/"

# Use git filter-branch to remove .env files from history
echo "Removing .env files from git history..."
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch sandbox/innovation-engine-headlamp/.env sandbox/innovation-engine-headlamp/.env.test" \
  --prune-empty --tag-name-filter cat -- --all

# Restore the local copies of .env files
echo "Restoring local copies of .env files..."
cp -f /tmp/env_backup/.env sandbox/innovation-engine-headlamp/.env 2>/dev/null
cp -f /tmp/env_backup/.env.test sandbox/innovation-engine-headlamp/.env.test 2>/dev/null

# Force garbage collection and remove old references
echo "Cleaning up git repository..."
git for-each-ref --format="delete %(refname)" refs/original | git update-ref --stdin
git reflog expire --expire=now --all
git gc --aggressive --prune=now

echo "Done. The .env files have been removed from git history but preserved in your working directory."
echo "They are now ignored by git as specified in .gitignore."
echo ""
echo "IMPORTANT: If this is a shared repository, you will need to force push these changes:"
echo "git push origin --force --all"
echo ""
echo "Make sure all collaborators pull the latest changes after this operation."
