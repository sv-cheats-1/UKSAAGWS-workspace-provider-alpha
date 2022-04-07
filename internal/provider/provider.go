package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var DefaultClientScopes = []string{
	//"https://www.googleapis.com/auth/gmail.settings.basic",
	//"https://www.googleapis.com/auth/gmail.settings.sharing",
	//"https://www.googleapis.com/auth/chrome.management.policy",
	"https://www.googleapis.com/auth/cloud-platform",
	//"https://www.googleapis.com/auth/admin.directory.customer",
	//"https://www.googleapis.com/auth/admin.directory.domain",
	//"https://www.googleapis.com/auth/admin.directory.group",
	//"https://www.googleapis.com/auth/admin.directory.orgunit",
	//"https://www.googleapis.com/auth/admin.directory.rolemanagement",
	//"https://www.googleapis.com/auth/admin.directory.userschema",
	//"https://www.googleapis.com/auth/admin.directory.user",
	//"https://www.googleapis.com/auth/apps.groups.settings",
	"https://apps-apis.google.com/a/feeds/domain/",
}

func init() {
	// Set descriptions to support Markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example, you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"customer_id": {
					Description: "The customer id provided with your Google Workspace subscription. It is found " +
						"in the admin console under Account Settings.",
					Type: schema.TypeString,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"GOOGLEWORKSPACE_CUSTOMER_ID",
					}, nil),
					Required: true,
				},
				"impersonated_user_email": {
					Description: "The impersonated user's email with access to the Admin APIs can access the Admin SDK Directory API. " +
						"`impersonated_user_email` is required for all services except group and user management.",
					Type: schema.TypeString,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL",
					}, nil),
					Required: true,
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
					ValidateDiagFunc: validateCredentials,
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"wpa_email_gateway": dataSourceWorkspaceEmailGateway(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"wpa_email_gateway": resourceWorkspaceEmailGateway(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

type apiClient struct {
	client *http.Client

	ClientScopes          []string
	Customer              string
	ImpersonatedUserEmail string
	UserAgent             string
	Credentials           string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		config := apiClient{}

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

		// We only support default scopes
		if len(config.ClientScopes) == 0 {
			config.ClientScopes = DefaultClientScopes
		}

		// Get credentials
		if v, ok := d.GetOk("credentials"); ok {
			config.Credentials = v.(string)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "credentials are required",
			})

			return nil, diags
		}

		config.UserAgent = p.UserAgent("terraform-provider-wpa", version)

		// nolint
		newCtx, _ := schema.StopContext(ctx)
		diags = config.loadAndValidate(newCtx)

		return &config, diags
	}
}

func (c *apiClient) loadAndValidate(ctx context.Context) diag.Diagnostics {

	var diags diag.Diagnostics

	if c.Credentials != "" {
		contents, _, err := pathOrContents(c.Credentials)
		if err != nil {
			return diag.FromErr(err)
		}

		credParams := googleoauth.CredentialsParams{
			Scopes:  c.ClientScopes,
			Subject: c.ImpersonatedUserEmail,
		}

		creds, err := googleoauth.CredentialsFromJSONWithParams(ctx, []byte(contents), credParams)
		if err != nil {
			return diag.FromErr(err)
		}

		diags = c.SetupClient(ctx, creds)

		return diags
	}

	// TODO: this does not work...
	credParams := googleoauth.CredentialsParams{
		Scopes:  c.ClientScopes,
		Subject: c.ImpersonatedUserEmail,
	}

	creds, err := googleoauth.FindDefaultCredentialsWithParams(ctx, credParams)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = c.SetupClient(ctx, creds)

	return diags
}

func (c *apiClient) SetupClient(ctx context.Context, creds *googleoauth.Credentials) diag.Diagnostics {

	var diags diag.Diagnostics

	cleanCtx := context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())

	// 1. mTLS TRANSPORT/CLIENT - sets up proper auth headers
	client, _, err := transport.NewHTTPClient(cleanCtx, option.WithTokenSource(creds.TokenSource))
	if err != nil {
		return diag.FromErr(err)
	}

	// Starting with the Workspace provider, we removed the advanced logging and retry middleware.

	c.client = client
	return diags
}

func pathOrContents(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, true, err
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}

func validateCredentials(v interface{}, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if v == nil || v.(string) == "" {
		return diags
	}
	creds := v.(string)
	path, err := homedir.Expand(creds)
	if err != nil {
		return diag.FromErr(err)
	}
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(path); err == nil {
		return diags
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds)); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("JSON credentials in %q are not valid: %s", creds, err),
			AttributePath: p,
		})
	}

	return diags
}
