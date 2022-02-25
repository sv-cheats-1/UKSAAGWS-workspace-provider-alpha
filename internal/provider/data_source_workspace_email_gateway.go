package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"net/http"
	"encoding/xml"
)

// Google Workspace Admin SDK -> Admin Settings API
const apiEndPoint string = "https://apps-apis.google.com/a/feeds/domain/2.0/%s/email/gateway"

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

type Response struct {
	Property struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"property"`
} 


func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	domainName := d.Get("domain_name").(string)
	d.SetId(domainName)

	c := meta.(*apiClient)

	// Send an HTTP request to this legacy API and process the response
	resp, err := c.client.Get(fmt.Sprintf(apiEndPoint, domainName))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error in HTTP response. Status code: %v, Status: %v. Full response: %#v", resp.StatusCode, resp.Status, resp)
	}

	
	// TODO: trim the XML parser struct to only have relevant fields and remove the unnecessary declarations 
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	 
	var r Response
	b_res := []byte(body)

	xml.Unmarshal(b_res, &r)

	if err := d.Set("smart_host", r.Property.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("smtp_mode", r.Property.Value); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
