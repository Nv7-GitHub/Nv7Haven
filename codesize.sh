SIZE="$(find . -type f -name '*.go' -exec du -ch {} + | grep total)"
LINES="$(find . -type f -name '*.go' | xargs wc -l | grep total)"

SIZE="${SIZE//total/}"
LINES="${LINES//total/}"

echo "Size: ${SIZE// /}"
echo "Lines: ${LINES// /}"
