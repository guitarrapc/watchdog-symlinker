workflow "Golang workflow" {
  on = "release"
  resolves = ["Release-Windows", "Release-Linux"]
}
 
action "GolangCI-Lint" {
  uses = "./.github/actions/golang"
  args = "lint"
}

# windows
action "Build-Windows" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
  env = {
    GOOS = "windows"
    GOARCH = "amd64"
  }
}


action "Release-Windows" {
  needs = ["Build-Windows"]
  uses = "./.github/actions/release"
  env = {
    FILE_NAME = "watchdog-symlinker.exe"
    GOOS = "windows"
    GOARCH = "amd64"
  }
}

# linux
action "Build-Linux" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
  env = {
    GOOS = "linux"
    GOARCH = "amd64"
  }
}

action "Release-Linux" {
  needs = ["Build-Linux"]
  uses = "./.github/actions/release"
  env = {
    FILE_NAME = "watchdog-symlinker"
    GOOS = "linux"
    GOARCH = "amd64"
  }
}
