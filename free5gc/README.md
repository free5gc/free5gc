# free5gc

## Description
free5gc

## Usage

### Fetch the package
`kpt pkg get REPO_URI[.git]/PKG_PATH[@VERSION] free5gc`
Details: https://kpt.dev/reference/cli/pkg/get/

### View package content
`kpt pkg tree free5gc`
Details: https://kpt.dev/reference/cli/pkg/tree/

### Apply the package
```
kpt live init free5gc
kpt live apply free5gc --reconcile-timeout=2m --output=table
```
Details: https://kpt.dev/reference/cli/live/
