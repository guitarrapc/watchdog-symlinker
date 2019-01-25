workflow "Golang workflow" {
  on = "release"
  resolves = ["RELEASE-Windows"]
}
 
action "GolangCI-Lint" {
  uses = "./.github/actions/golang"
  args = "lint"
}
 
action "Build-Windows" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
  env = {
    GOOS = "windows"
    GOARCH = "amd64"
  }
}

action "Build-Linux" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
  env = {
    GOOS = "linux"
    GOARCH = "amd64"
  }
}

action "RELEASE-Windows" {
  needs = ["Build-Windows"]
  uses = "./.github/actions/release"
  args = "watchdog-symlinker.exe"
}
