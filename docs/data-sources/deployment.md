---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qovery_deployment Data Source - terraform-provider-qovery"
subcategory: ""
description: |-
  Use this data source to retrieve information about an existing deployment.
---

# qovery_deployment (Data Source)

Use this data source to retrieve information about an existing deployment.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) Id of the environment to deploy.

### Optional

- `desired_state` (String) Desired state of the deployment.
- `id` (String) Id of the deployment
- `version` (String) Version to force trigger a deployment when desired_state doesn't change (e.g redeploy a deployment having the 'RUNNING' state)

