workflow "Golang workflow" {
  on = "push"
  resolves = ["Build windows", "Build Linux"]
}
 
action "GolangCI-Lint" {
  uses = "./.github/actions/golang"
  args = "lint"
}
 
action "Build windows" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
  env = {
    GOOS = "windows"
    GOARCH = "amd64"
  }
}
action "Build linux" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
  env = {
    GOOS = "linux"
    GOARCH = "amd64"
  }
}