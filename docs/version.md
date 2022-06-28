A newer version `y` being compatible with an older version `x` means `y` can
load data serialized by `x`. But `x` should never try to load data serialized by
a newer version `y`.

- `v0.5.*` is compatible with `0.2.*`, `0.3.*`, `0.4.*`, `0.5.*`.
- `v0.4.*` is compatible with `0.2.*`, `0.3.*`, `0.4.*`.
- `v0.3.*` is compatible with `0.2.*`, `0.3.*`.
- `v0.2.*` is compatible with `0.2.*`.
