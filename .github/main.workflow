workflow "Golang workflow" {
  on = "push"
  resolves = ["Build"]
}
 
action "GolangCI-Lint" {
  uses = "./.github/actions/golang"
  args = "lint"
}
 
action "Build" {
  needs = ["GolangCI-Lint"]
  uses = "./.github/actions/golang"
  args = "build"
}