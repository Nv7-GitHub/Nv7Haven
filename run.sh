git pull

export PORT=8080
export PHP_PORT=3000
export MYSQL_HOST=127.0.0.1

echo "Building..."
go build -o main -tags="arm_logs"

echo "Running..."
until ./main; do
  echo "Server 'nv7haven' crashed with exit code $?.  Respawning.." >&2
  sleep 1
done