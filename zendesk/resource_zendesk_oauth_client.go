package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/oauth/oauth_clients/
func resourceZendeskOAuthClient() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an OAuth client resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createOAuthClient(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readOAuthClient(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateOAuthClient(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteOAuthClient(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this OAuth client.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The id of the user the OAuth client belongs to.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The name of the OAuth client, shown to users asked to grant it access.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"identifier": {
				Description: "The unique identifier of the OAuth client.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"kind": {
				Description: "Whether the client can keep its secret. Only a confidential client may use the " +
					"client_credentials grant.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "confidential"}, false),
			},
			"company": {
				Description: "The company name shown to users asked to grant the client access.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": {
				Description: "The description shown to users asked to grant the client access.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"scopes": {
				Description: "The scopes the OAuth client may request, e.g. `tickets:read`.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"redirect_uri": {
				Description: "The redirect URIs accepted for the authorization code grant. Unused by the " +
					"client_credentials grant, where an empty list is accepted.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enforce_token_expiration": {
				Description: "Whether tokens issued for this client expire.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"logo_url": {
				Description: "The url of the logo shown to users asked to grant the client access.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"global": {
				Description: "Whether the OAuth client is available to every Zendesk account.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"secret": {
				Description: "The client secret. Zendesk returns it in full only when the client is created and " +
					"never again, so a client brought under Terraform by import has no secret in state.",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"created_at": {
				Description: "The time the OAuth client was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the OAuth client was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// marshalOAuthClient deliberately leaves `secret` alone. Zendesk answers every read after the create with a
// truncated stub, so writing it back would replace the real secret in state with a value that cannot
// authenticate, and leave a diff that no apply can settle. Only createOAuthClient sets it.
func marshalOAuthClient(oauthClient newClient.OAuthClient, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":          oauthClient.URL,
		"user_id":      oauthClient.UserID,
		"name":         oauthClient.Name,
		"identifier":   oauthClient.Identifier,
		"kind":         oauthClient.Kind,
		"company":      oauthClient.Company,
		"description":  oauthClient.Description,
		"scopes":       oauthClient.Scopes,
		"redirect_uri": oauthClient.RedirectURI,
		"logo_url":     oauthClient.LogoURL,
		"global":       oauthClient.Global,
		"created_at":   oauthClient.CreatedAt,
		"updated_at":   oauthClient.UpdatedAt,
	}

	if oauthClient.EnforceTokenExpiration != nil {
		fields["enforce_token_expiration"] = *oauthClient.EnforceTokenExpiration
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalOAuthClient(d identifiableGetterSetter) (newClient.OAuthClient, error) {
	oauthClient := newClient.OAuthClient{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return oauthClient, fmt.Errorf("could not parse oauth client id %s: %v", v, err)
		}
		oauthClient.ID = id
	}

	if v, ok := d.GetOk("name"); ok {
		oauthClient.Name = v.(string)
	}

	if v, ok := d.GetOk("identifier"); ok {
		oauthClient.Identifier = v.(string)
	}

	if v, ok := d.GetOk("kind"); ok {
		oauthClient.Kind = v.(string)
	}

	if v, ok := d.GetOk("company"); ok {
		oauthClient.Company = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		oauthClient.Description = v.(string)
	}

	if v, ok := d.GetOk("scopes"); ok {
		oauthClient.Scopes = expandStringList(v)
	}

	if v, ok := d.GetOk("redirect_uri"); ok {
		oauthClient.RedirectURI = expandStringList(v)
	}

	enforceTokenExpiration := d.Get("enforce_token_expiration").(bool)
	oauthClient.EnforceTokenExpiration = &enforceTokenExpiration

	return oauthClient, nil
}

func createOAuthClient(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	oauthClient, err := unmarshalOAuthClient(d)
	if err != nil {
		return diag.FromErr(err)
	}

	oauthClient, err = zd.CreateOAuthClient(ctx, oauthClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", oauthClient.ID))

	err = marshalOAuthClient(oauthClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// The only response carrying the full secret. Zendesk cannot be asked for it again; the sole recovery is
	// generate_secret, which invalidates whatever the client is already authenticating with.
	err = d.Set("secret", oauthClient.Secret)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readOAuthClient(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	oauthClient, err := zd.GetOAuthClient(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOAuthClient(oauthClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateOAuthClient(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	oauthClient, err := unmarshalOAuthClient(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	oauthClient, err = zd.UpdateOAuthClient(ctx, id, oauthClient)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOAuthClient(oauthClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteOAuthClient(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteOAuthClient(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
