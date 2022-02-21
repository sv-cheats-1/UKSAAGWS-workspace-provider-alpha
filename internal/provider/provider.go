package provider

import (
	"context"
	"net/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DefaultClientScopes = []string{
	"https://www.googleapis.com/auth/gmail.settings.basic",
	"https://www.googleapis.com/auth/gmail.settings.sharing",
	"https://www.googleapis.com/auth/chrome.management.policy",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/admin.directory.customer",
	"https://www.googleapis.com/auth/admin.directory.domain",
	"https://www.googleapis.com/auth/admin.directory.group",
	"https://www.googleapis.com/auth/admin.directory.orgunit",
	"https://www.googleapis.com/auth/admin.directory.rolemanagement",
	"https://www.googleapis.com/auth/admin.directory.userschema",
	"https://www.googleapis.com/auth/admin.directory.user",
	"https://www.googleapis.com/auth/apps.groups.settings",
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"wpa_email_gateway": dataSourceWorkspaceEmailGateway(),
			},
			/*
				ResourcesMap: map[string]*schema.Resource{
					"scaffolding_resource": resourceScaffolding(),
				},
			*/
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	client *http.Client

	AccessToken           string
	ClientScopes          []string
	Credentials           string
	Customer              string
	ImpersonatedUserEmail string
	ServiceAccount        string
	UserAgent             string
}

// func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
// 	return func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
// 		// Setup a User-Agent for your API client (replace the provider name for yours):
// 		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
// 		// TODO: myClient.UserAgent = userAgent

// 		return &apiClient{}, nil
// 	}
// }
func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		config := apiClient{}

		// Get access token
		if v, ok := d.GetOk("access_token"); ok {
			config.AccessToken = v.(string)
		}

		// Get credentials
		if v, ok := d.GetOk("credentials"); ok {
			config.Credentials = v.(string)
		}

		// Get customer id
		if v, ok := d.GetOk("customer_id"); ok {
			config.Customer = v.(string)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "customer_id is required",
			})

			return nil, diags
		}

		// Get impersonated user email
		if v, ok := d.GetOk("impersonated_user_email"); ok {
			config.ImpersonatedUserEmail = v.(string)
		}

		// Get scopes
		scopes := d.Get("oauth_scopes").([]interface{})
		if len(scopes) > 0 {
			config.ClientScopes = make([]string, len(scopes))
		}
		for i, scope := range scopes {
			config.ClientScopes[i] = scope.(string)
		}

		// Get service account
		if v, ok := d.GetOk("service_account"); ok {
			config.ServiceAccount = v.(string)
		}

		config.UserAgent = p.UserAgent("terraform-provider-wpa", version)

		// nolint
		// newCtx, _ := schema.StopContext(ctx)
		// diags = config.loadAndValidate(newCtx)

		return &config, diags
	}
}
