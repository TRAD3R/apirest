#!/bin/bash
set -e

# Создание базы данных
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE "$POSTGRES_DB";
EOSQL