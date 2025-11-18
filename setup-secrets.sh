#!/bin/bash
# Setup script for Docker secrets in CodexGigantus
# This script helps create secure passwords for database access

set -e

SECRETS_DIR=".secrets"
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}CodexGigantus Security Setup${NC}"
echo "================================"
echo ""

# Create secrets directory
if [ ! -d "$SECRETS_DIR" ]; then
    mkdir -p "$SECRETS_DIR"
    chmod 700 "$SECRETS_DIR"
    echo -e "${GREEN}✓${NC} Created secrets directory"
else
    echo -e "${YELLOW}!${NC} Secrets directory already exists"
fi

# Function to generate a secure random password
generate_password() {
    openssl rand -base64 32 | tr -d "=+/" | cut -c1-25
}

# Function to create or update a secret
create_secret() {
    local secret_name=$1
    local secret_file="${SECRETS_DIR}/${secret_name}"

    if [ -f "$secret_file" ]; then
        echo -e "${YELLOW}!${NC} $secret_name already exists"
        read -p "Do you want to regenerate it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return
        fi
    fi

    # Ask user if they want to provide their own password or generate one
    echo ""
    echo "Creating: $secret_name"
    read -p "Generate a secure password automatically? (Y/n): " -n 1 -r
    echo

    if [[ $REPLY =~ ^[Nn]$ ]]; then
        # User wants to provide their own password
        read -sp "Enter password: " password
        echo
        read -sp "Confirm password: " password_confirm
        echo

        if [ "$password" != "$password_confirm" ]; then
            echo -e "${RED}✗${NC} Passwords do not match!"
            return 1
        fi

        if [ ${#password} -lt 12 ]; then
            echo -e "${RED}✗${NC} Password must be at least 12 characters!"
            return 1
        fi
    else
        # Generate a secure password
        password=$(generate_password)
    fi

    # Write password to file
    echo -n "$password" > "$secret_file"
    chmod 600 "$secret_file"

    echo -e "${GREEN}✓${NC} Created $secret_name"

    if [[ $REPLY =~ ^[Nn]$ ]]; then
        echo "  (using your provided password)"
    else
        echo "  Password: $password"
        echo -e "  ${YELLOW}IMPORTANT: Save this password securely!${NC}"
    fi
}

# Create secrets
echo ""
echo "Setting up database secrets..."
echo "-------------------------------"
create_secret "db_admin_password"
echo ""
create_secret "db_password"

# Create .gitignore if it doesn't exist
if [ ! -f "${SECRETS_DIR}/.gitignore" ]; then
    echo "*" > "${SECRETS_DIR}/.gitignore"
    echo "!.gitignore" >> "${SECRETS_DIR}/.gitignore"
    echo -e "${GREEN}✓${NC} Created .gitignore in secrets directory"
fi

# Update .env.example
echo ""
echo "Updating configuration files..."
echo "-------------------------------"

if [ ! -f ".env" ] && [ -f ".env.example" ]; then
    echo -e "${YELLOW}!${NC} .env file not found, copying from .env.example"
    cp .env.example .env
    echo -e "${GREEN}✓${NC} Created .env from .env.example"
fi

echo ""
echo -e "${GREEN}✓ Setup complete!${NC}"
echo ""
echo "Next steps:"
echo "1. Review the generated passwords in $SECRETS_DIR/"
echo "2. Update .env file if needed for local development"
echo "3. Run: docker-compose up -d"
echo ""
echo -e "${YELLOW}SECURITY NOTES:${NC}"
echo "• Never commit files in $SECRETS_DIR/ to version control"
echo "• For production, use a proper secrets management system (e.g., Vault, AWS Secrets Manager)"
echo "• Change the db_admin_password regularly"
echo "• The application uses codex_readonly with minimal (SELECT-only) privileges"
echo "• Adminer is disabled by default (use --profile debug to enable)"
echo ""
