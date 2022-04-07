package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkspaceEmailGateway() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceWorkspaceEmailGateway().Schema)
	addRequiredFieldsToSchema(dsSchema, "domain_name")

	return &schema.Resource{
		Description: "The outbound email gateway provides outbound routing of mail from users in your domain.",

		ReadContext: dataSourceWorkspaceRead,

		Schema: dsSchema,
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	domainName := d.Get("domain_name").(string)
	d.SetId(domainName)
	return resourceWorkspaceRead(ctx, d, meta)
}

/*
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

	if err := xml.Unmarshal(b_res, &res); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("smart_host", res.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("smtp_mode", res.Value); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
*/
