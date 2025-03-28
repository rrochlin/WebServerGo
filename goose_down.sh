export base=$(pwd)
cd /home/rrochlin/boot_dev/WebServerGo/sql/schema/
goose postgres "postgres://postgres:postgres@localhost:5432/chirpygo" down
cd $base
