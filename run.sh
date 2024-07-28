echo "Updating..."
 > logs.txt
git pull --recurse-submodules

export MYSQL_HOST=127.0.0.1

echo "Building..."
go build -o main -tags="arm_logs" -ldflags="-s -w"

echo "Running..."
until ./main > output.log 2>&1; do
  echo "Go crashed with exit code $?.  Respawning.." >&2
  sleep 1
done