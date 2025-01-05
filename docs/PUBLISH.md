# Publish a module

1. `go mod tidy`
2. `go test ./...`
3. `git add .`
4. `git commit -m "mymodule: changes for v0.1.0"` (select new version from GitHub)
5. `git tag v0.1.0`
6. `git push origin v0.1.0`
7. `git push`
