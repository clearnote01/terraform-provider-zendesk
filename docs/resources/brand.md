---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zendesk_brand Resource - terraform-provider-zendesk"
subcategory: ""
description: |-
  Provides a brand resource.
---

# zendesk_brand (Resource)

Provides a brand resource.

## Example Usage

```terraform
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/brands

resource "zendesk_brand" "T-800" {
  name            = "T-800"
  active          = true
  subdomain       = "d3v-terraform-provider-t800"
  # TODO: logo
}

resource "zendesk_brand" "T-1000" {
  name            = "T-1000"
  active          = false
  subdomain       = "d3v-terraform-provider-t1000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the brand.
- `subdomain` (String) The subdomain of the brand.

### Optional

- `active` (Boolean) If the brand is set as active.
- `default` (Boolean) Is the brand the default brand for this account.
- `host_mapping` (String) The hostmapping to this brand, if any. Only admins view this property.
- `id` (String) The ID of this resource.
- `logo_attachment_id` (Number) Logo attachment id for the brand.
- `signature_template` (String) The signature template for a brand.

### Read-Only

- `brand_url` (String) The url of the brand.
- `has_help_center` (Boolean) If the brand has a Help Center.
- `help_center_state` (String) The state of the Help Center. Allowed values are "enabled", "disabled", or "restricted".
- `ticket_form_ids` (Set of Number) The ids of ticket forms that are available for use by a brand.
- `url` (String) The API url of this brand.

