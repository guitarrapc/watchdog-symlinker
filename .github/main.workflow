workflow "Golang workflow" {
  on = "push"
  resolves = ["Build-Windows", "Build-Linux"]
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
