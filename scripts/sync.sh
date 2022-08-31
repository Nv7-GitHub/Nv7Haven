echo "Removing..."
rm -rf "data"
echo "Downloading..."
curl "https://nv7haven.nv7haven.com/sync_db/$PASSWORD" > sync.zip
echo "Unzipping..."
unzip sync.zip -d .
echo "Cleaning up..."
rm sync.zip
echo "Done!"
