echo "Updating..."
 > logs.txt
git pull --recurse-submodules

export ROCKET_PORT=49154

echo "Building..."
cargo build --release

echo "Running..."
until ./target/release/nv7haven; do
  echo "Rust crashed with exit code $?.  Respawning.." >&2
  sleep 1
done