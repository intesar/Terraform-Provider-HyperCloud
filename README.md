terraform-provider-hypercloud
=========================

Terraform provider plugin to access [HyperCloud](https://hypercloud.hypergrid.com)

Usage
-------------------------

```
provider "hypercloud" {
    url       = "https://<ip>:<port>/",
    accessKey = "<username>",
    secretKey = "<password>",
}

resource "hypercloud_compute" "compute-001" {
    Blueprint-ID = "<402881905c0e96c2015c0ea31b370019>"
}

```

Installation
-------------------------

```
go get -u github.com/intesar/Terraform-Provider-HyperCloud
```
