# hatlas

This is a tool to help script Hashicorp Atlas API

## Install

### Using go

```
$ go get -u github.com/BSick7/hatlas
$ hatlas terra list
```

### Downloading on Linux

```
$ HATLAS_VERSION=0.1.0
$ curl -sSL https://github.com/BSick7/hatlas/releases/download/${HATLAS_VERSION}/hatlas_linux_amd64 > hatlas \
    && chmod 755 hatlas \
    && mv hatlas /usr/local/bin/hatlas
$ hatlas terra list
```

### Downloading on Mac

```
$ HATLAS_VERSION=0.1.0
$ curl -sSL https://github.com/BSick7/hatlas/releases/download/${HATLAS_VERSION}/hatlas_darwin_amd64 > hatlas \
    && chmod 755 hatlas \
    && mv hatlas /usr/local/bin/hatlas
$ hatlas terra list
```

#### Downloading on Windows 

```
$ HATLAS_VERSION=0.1.0
$ curl -sSL https://github.com/BSick7/hatlas/releases/download/${HATLAS_VERSION}/hatlas_windows_amd64.exe > hatlas.exe
# Move somewhere on PATH
$ hatlas terra list
```

## Commands

```
hatlas terra
  list      List environments
  state     Introspect state file
  outputs   Introspect outputs
  config    Introspect configuration
```
