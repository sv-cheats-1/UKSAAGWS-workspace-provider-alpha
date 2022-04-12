package provider

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"log"
	"net/http"
	"strings"
)

// Google Workspace Admin SDK -> Admin Settings API
const apiEndPoint string = "https://apps-apis.google.com/a/feeds/domain/2.0/%s/email/gateway"

func resourceWorkspaceEmailGateway() *schema.Resource {
	return &schema.Resource{
		Description: "The outbound email gateway provides outbound routing of mail from users in your domain.",

		CreateContext: resourceGatewayCreate,
		ReadContext:   resourceWorkspaceRead,
		DeleteContext: resourceGatewayDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Description: "Your Google Workspace domain name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"smart_host": {
				Description: "Either the IP address or hostname of your SMTP server. Google Workspace routes outgoing mail to this server.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"smtp_mode": {
				Description: "The default value is SMTP. Another value, SMTP_TLS, secures a connection with TLS when delivering the message.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"response": {
				Description: "The default value is SMTP. Another value, SMTP_TLS, secures a connection with TLS when delivering the message.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

type Property struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Response struct {
	XMLName  xml.Name `xml:"entry"`
	Text     string   `xml:",chardata"`
	Xmlns    string   `xml:"xmlns,attr"`
	Apps     string   `xml:"apps,attr"`
	Property []struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"property"`
}

type Request struct {
	XMLName  xml.Name   `xml:"atom:entry"`
	Text     string     `xml:",chardata"`
	Atom     string     `xml:"xmlns:atom,attr"`
	Apps     string     `xml:"xmlns:apps,attr"`
	Property []Property //`xml:"property"`
}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	domainName := d.Get("domain_name").(string)
	d.SetId(domainName)

	c := meta.(*apiClient)

	// TODO create XML and output it as a string

	p1 := Property{
		Name:  "smartHost",
		Value: d.Get("smart_host").(string),
	}

	p2 := Property{
		Name:  "smtpMode",
		Value: d.Get("smtp_mode").(string),
	}

	p := []Property{p1, p2}

	r := Request{
		Atom:     "http://www.w3.org/2005/Atom",
		Apps:     "http://schemas.google.com/apps/2006",
		Property: p,
	}

	bXml, err := xml.Marshal(r)
	if err != nil {
		fmt.Printf("Marshaling error: %v\n", err)
		return nil
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(apiEndPoint, domainName), strings.NewReader(string(bXml)))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error in HTTP response. Status code: %v, Status: %v. Full response: %#v", resp.StatusCode, resp.Status, resp)
	}

	// TODO no error, process the response

	return diag.Diagnostics{}
}
func resourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	respRead := Response{}

	if err := xml.Unmarshal([]byte(body), &respRead); err != nil {
		fmt.Printf("Unmarshaling error: %v\n", err)
		return diag.FromErr(err)
	}

	if err := d.Set("smart_host", respRead.Property[0].Value); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("smtp_mode", respRead.Property[1].Value); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	domainName := d.Get("domain_name").(string)
	d.SetId(domainName)

	log.Printf("[DEBUG] Deleting Domain %q: %#v", d.Id(), domainName)

	c := meta.(*apiClient)

	// TODO create XML and output it as a string

	p1 := Property{
		Name:  "smartHost",
		Value: "",
	}

	p2 := Property{
		Name:  "smtpMode",
		Value: "",
	}

	p := []Property{p1, p2}

	r := Request{
		Atom:     "http://www.w3.org/2005/Atom",
		Apps:     "http://schemas.google.com/apps/2006",
		Property: p,
	}

	bXml, err := xml.Marshal(r)
	if err != nil {
		fmt.Printf("Marshaling error: %v\n", err)
		return nil
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(apiEndPoint, domainName), strings.NewReader(string(bXml)))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error in HTTP response. Status code: %v, Status: %v. Full response: %#v", resp.StatusCode, resp.Status, resp)
	}

	log.Printf("[DEBUG] Finished deleting Domain %q: %#v", d.Id(), domainName)

	return diag.Diagnostics{}
}
