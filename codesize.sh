echo "Code size:"
find . -type f -name '*.go' -exec du -ch {} + | grep total
echo "Number of lines"
find . -type f -name '*.go' | xargs wc -l | grep total
