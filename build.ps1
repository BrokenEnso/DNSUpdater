$env:GOOS = "linux" 
$env:GOARCH = "amd64" 

go build -o out/dnsupdater main.go