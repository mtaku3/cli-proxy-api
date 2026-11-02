# cli-proxy-api

[CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) built from an upstream release tag with the patches in `patches/` applied, published as `ghcr.io/mtaku3/cli-proxy-api`.

The upstream version is the `CLIPROXYAPI_VERSION` build arg in `Dockerfile`; Renovate bumps it. The patch currently carried adds `REQUIRE_API_KEY`, which rejects requests when no api-keys are configured.
