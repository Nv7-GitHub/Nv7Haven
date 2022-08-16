echo "Removing..."
rm -rf "$(go env GOPATH)/src/github.com/Nv7-Github/Nv7haven/data"
echo "Copying..."
scp -P 119 -r pi@98.247.143.47:/home/pi/go/src/github.com/Nv7-Github/Nv7haven/data "$(go env GOPATH)/src/github.com/Nv7-Github/Nv7haven/data" 
