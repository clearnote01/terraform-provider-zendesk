# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/oauth/oauth_clients/

resource "zendesk_oauth_client" "example" {
  name       = "Example Client"
  identifier = "example-client"

  # Only a confidential client may use the client_credentials grant, which is what a server-side
  # service wants: it exchanges the secret below for an access token with no user interaction.
  kind = "confidential"

  company     = "Example Company"
  description = "An example OAuth client"
  scopes      = ["tickets:read"]
}

# Zendesk returns the secret in full only when the client is created, and answers every later read with a
# truncated stub. Take it from here rather than from the API, and treat the state file as sensitive.
output "example_client_secret" {
  value     = zendesk_oauth_client.example.secret
  sensitive = true
}
