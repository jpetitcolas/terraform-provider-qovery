---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qovery_container_registry Resource - terraform-provider-qovery"
subcategory: ""
description: |-
  Provides a Qovery container registry resource. This can be used to create and manage Qovery container registry.
---

# qovery_container_registry (Resource)

Provides a Qovery container registry resource. This can be used to create and manage Qovery container registry.

## Example Usage

```terraform
resource "qovery_container_registry" "my_container_registry" {
  # Required
  organization_id = qovery_organization.my_organization.id
  name            = "my_aws_creds"
  kind            = "DOCKER_HUB"
  url             = "https://docker.io"
  config = {
    username = "<my_username>"
    password = "<my_password>"
  }

  # Optional
  description = "My Docker Hub Registry"

  depends_on = [
    qovery_organization.my_organization
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kind` (String) Kind of the container registry.
	- Can be: `DOCKER_HUB`, `DOCR`, `ECR`, `PUBLIC_ECR`, `SCALEWAY_CR`.
- `name` (String) Name of the container registry.
- `organization_id` (String) Id of the organization.
- `url` (String) URL of the container registry.

### Optional

- `config` (Map of String) Configuration needed to authenticate the container registry.
- `description` (String) Description of the container registry.

### Read-Only

- `id` (String) Id of the container registry.

## Import

Import is supported using the following syntax:

```shell
terraform import qovery_container_registry.my_container_registry "<organization_id>,<container_registry_id>"
```