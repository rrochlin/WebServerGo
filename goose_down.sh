#!/bin/bash

base=$(pwd)
script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")
target_dir="$script_dir/sql/schema/"
set -o allexport
source .env
set +o allexport

cd "$target_dir" || {
    cd "$base"
    exit 1
}

goose postgres "$DB_URL" down

cd "$base"
