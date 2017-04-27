## 0.2.2 (Unreleased)

FEATURES:

  * Added `hatlas artifacts push` command to push artifact metadata to atlas registry.

## 0.2.1 (Apr 24, 2017)

FEATURES:

  * Added `hatlas artifacts list` command to list artifacts in atlas registry.
  * Added `hatlas artifacts meta` command to dump artifact metadata from atlas registry.

## 0.2.0 (Apr 20, 2017)

FEATURES:

  * Added `hatlas terra push` command with 2 subcommands: `vars` and `state`.

INCOMPATIBILITIES:

  * Replaced `push-vars` command with `push vars`.

## 0.1.4 (Apr 19, 2017)

FEATURES:

  * Added `push-vars` command to update remote vars.
  * Added `-raw` flag to `config` command to dump in raw json format.
  * Added `ATLAS_ADDR` environment variable.

## 0.1.3 (Mar 21, 2017)

BUG FIXES:

  * Fixed `terra list` command truncating to just page 1.

## 0.1.2 (Jan 24, 2017)

BUG FIXES:

  * Fixed format of dumped terraform environment config.

## 0.1.1 (Jan 21, 2017)

FEATURES:

  * Added ability to dump single var from terraform environment config.

## 0.1.0 (Nov 12, 2016)

First build of hatlas