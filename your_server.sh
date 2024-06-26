set -e
tmpFile=$(mktemp)
go build -o "$tmpFile" cmd/main.go
exec "$tmpFile" "$@"
