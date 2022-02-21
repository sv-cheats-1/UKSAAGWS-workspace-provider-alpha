package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkspaceEmailGateway() *schema.Resource {
	return &schema.Resource{
		Description: "The outbound email gateway provides outbound routing of mail from users in your domain.",

		ReadContext: dataSourceWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Description: "Your Google Workspace domain name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"smart_host": {
				Description: "Either the IP address or hostname of your SMTP server. Google Workspace routes outgoing mail to this server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"smtp_mode": {
				Description: "The default value is SMTP. Another value, SMTP_TLS, secures a connection with TLS when delivering the message.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"access_token": {
				Description: "A temporary [OAuth 2.0 access token] obtained from " +
					"the Google Authorization server, i.e. the `Authorization: Bearer` token used to " +
					"authenticate HTTP requests to Google Admin SDK APIs. This is an alternative to `credentials`, " +
					"and ignores the `scopes` field. If both are specified, `access_token` will be " +
					"used over the `credentials` field.",
				Type:     schema.TypeString,
				Optional: true,
			},

			"credentials": {
				Description: "Either the path to or the contents of a service account key file in JSON format " +
					"you can manage key files using the Cloud Console).  If not provided, the application default " +
					"credentials will be used.",
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLEWORKSPACE_CREDENTIALS",
					"GOOGLEWORKSPACE_CLOUD_KEYFILE_JSON",
					"GOOGLE_CREDENTIALS",
				}, nil),
				// ValidateDiagFunc: validateCredentials,
			},

			"customer_id": {
				Description: "The customer id provided with your Google Workspace subscription. It is found " +
					"in the admin console under Account Settings.",
				Type: schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLEWORKSPACE_CUSTOMER_ID",
				}, nil),
				Optional: true,
			},

			"impersonated_user_email": {
				Description: "The impersonated user's email with access to the Admin APIs can access the Admin SDK Directory API. " +
					"`impersonated_user_email` is required for all services except group and user management.",
				Type: schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL",
				}, nil),
				Optional: true,
			},

			"oauth_scopes": {
				Description: "The list of the scopes required for your application (for a list of possible scopes, see " +
					"[Authorize requests](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing))",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"service_account": {
				Description: "The service account used to create the provided `access_token` if authenticating using " +
					"the `access_token` method and needing to impersonate a user. This service account will require the " +
					"GCP role `Service Account Token Creator` if needing to impersonate a user.",
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId(d.Get("domain_name").(string))
	d.SetId(d.Get("customer_id").(string))

	// TODO request data from REST API and process the result

	if err := d.Set("smart_host", "TODO smart_host"); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("smtp_mode", "TODO smtp_mode"); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
