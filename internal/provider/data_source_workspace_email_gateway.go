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
		},
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId(d.Get("domain_name").(string))

	// TODO request data from REST API and process the result

	if err := d.Set("smart_host", "TODO smart_host"); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("smtp_mode", "TODO smtp_mode"); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
