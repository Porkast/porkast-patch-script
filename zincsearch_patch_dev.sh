#!/bin/sh

export ZINC_FIRST_ADMIN_USER="admin"
export ZINC_FIRST_ADMIN_PASSWORD="1qaz!QAZ"
export ZINC_BASE_URL="http://localhost:4080"

go run main.go PatchZincsearchData
