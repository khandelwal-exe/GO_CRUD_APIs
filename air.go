# air.toml

# The entry point of your application
[build]
  cmd = "main.go"

# Folders or files to watch for changes
[build.exec]
  command = "go"
  args = ["run", "main.go"]
  include = [".go"]
  exclude = ["vendor"]

# Environment variables to set
[build.env]
  PORT = "8080"
  ENV = "development"

# Configuration for hot reloading
[build.log]
  color = "auto"
