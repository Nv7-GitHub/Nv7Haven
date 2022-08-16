echo "Removing..."
rm -rf "$(go env GOPATH)/src/github.com/Nv7-Github/Nv7haven/data"
echo "Downloading..."
curl "https://nv7haven.nv7haven.com/sync_db/$PASSWORD" > sync.zip
echo "Unzipping..."
unzip sync.zip -d "$(go env GOPATH)/src/github.com/Nv7-Github/Nv7haven"
echo "Cleaning up..."
rm sync.zip
echo "Done!"
