{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  env.CGO_ENABLED = 0;
  dotenv.enable = true;
  packages = [
    pkgs.lefthook
    pkgs.commitlint-rs
    pkgs.hadolint
    pkgs.shellcheck
    pkgs.container-structure-test
    pkgs.docker
    pkgs.go-task
    pkgs.go
    pkgs.golangci-lint
    pkgs.gotestsum
    pkgs.gofumpt
    pkgs.delve
    pkgs.gomarkdoc
    pkgs.goreleaser
  ];

  enterShell = ''
    lefthook install
  '';
}
