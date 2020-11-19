package newrelic

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsMonitorScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicSyntheticsMonitorScriptCreate,
		Read:   resourceNewRelicSyntheticsMonitorScriptRead,
		Update: resourceNewRelicSyntheticsMonitorScriptUpdate,
		Delete: resourceNewRelicSyntheticsMonitorScriptDelete,
		Importer: &schema.ResourceImporter{
			State: importSyntheticsMonitorScript,
		},
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the monitor to attach the script to.",
			},
			"text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plaintext representing the monitor script.",
			},
			"locations": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the monitor script location of execution.",
						},
						"hmac": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A cryptographic hash.",
						},
					},
				},
			},
		},
	}
}

func importSyntheticsMonitorScript(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("monitor_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

func buildSyntheticsMonitorScriptStruct(d *schema.ResourceData) *synthetics.MonitorScript {
	txt := d.Get("text").(string)
	locations := d.Get("locations").([]interface{})

	scriptTextEncoded := base64.StdEncoding.EncodeToString([]byte(txt))

	script := synthetics.MonitorScript{
		Text:      scriptTextEncoded,
		Locations: expandLocations(locations, scriptTextEncoded),
	}

	// fmt.Print("\n\n **************************** \n")
	// fmt.Printf("\n buildSyntheticsMonitorScriptStruct - Locations:  %+v \n", script)
	// fmt.Print("\n **************************** \n\n")
	// time.Sleep(3 * time.Second)

	return &script
}

func expandLocations(locations []interface{}, scriptText string) []synthetics.MonitorScriptLocation {
	out := make([]synthetics.MonitorScriptLocation, len(locations))

	for i, l := range locations {
		loc := l.(map[string]interface{})

		name := loc["name"].(string)
		// hmac := loc["hmac"].(string)

		// Using hardcoded values to prove out the concept.
		// The user should provide a hashed string so they aren't not storing
		// sensitive values in their HCL and state.
		secret := "password123"
		data := scriptText
		fmt.Printf("Secret: %s \nData: %s\n", secret, data)

		// Create a new HMAC by defining the hash type and the key (as byte array)
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(data))

		out[i] = synthetics.MonitorScriptLocation{
			Name: name,
			// HMAC: base64.StdEncoding.EncodeToString([]byte(scriptText)),
			HMAC: hex.EncodeToString(h.Sum(nil)),
		}
	}

	return out
}

func flattenLocations(locations []synthetics.MonitorScriptLocation) []map[string]string {
	out := make([]map[string]string, len(locations))

	fmt.Print("\n\n **************************** \n")
	fmt.Printf("\n flattenLocations - IN:  %+v \n", locations)

	for i, l := range locations {
		out[i] = map[string]string{
			"name": l.Name,
			"hmac": l.HMAC,
		}
	}

	fmt.Printf("\n flattenLocations - OUT:  %+v \n", out)
	fmt.Print("\n **************************** \n\n")
	time.Sleep(7 * time.Second)

	return out
}

func resourceNewRelicSyntheticsMonitorScriptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id := d.Get("monitor_id").(string)
	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", id)

	resp, err := client.Synthetics.UpdateMonitorScript(id, *buildSyntheticsMonitorScriptStruct(d))
	if err != nil {
		return err
	}

	fmt.Printf("\n PUT - resp:  %+v \n", resp)
	fmt.Print("\n **************************** \n\n")
	time.Sleep(7 * time.Second)

	d.SetId(id)
	return resourceNewRelicSyntheticsMonitorScriptRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics script %s", d.Id())

	script, err := client.Synthetics.GetMonitorScript(d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	fmt.Printf("\n GET - resp:  %+v \n", script)
	fmt.Print("\n **************************** \n\n")
	time.Sleep(7 * time.Second)

	d.Set("text", script.Text)

	locations := flattenLocations(script.Locations)

	if err := d.Set("locations", locations); err != nil {
		return err
	}

	return nil
}

func resourceNewRelicSyntheticsMonitorScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", d.Id())

	_, err := client.Synthetics.UpdateMonitorScript(d.Id(), *buildSyntheticsMonitorScriptStruct(d))
	if err != nil {
		log.Printf("[ERROR] updating monitor script failed: %v", err)
		return err
	}

	d.SetId(d.Id())
	return resourceNewRelicSyntheticsMonitorScriptRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics monitor script %s", d.Id())

	script := synthetics.MonitorScript{
		Text: " ",
	}

	if _, err := client.Synthetics.UpdateMonitorScript(d.Id(), script); err != nil {
		return err
	}

	return nil
}
