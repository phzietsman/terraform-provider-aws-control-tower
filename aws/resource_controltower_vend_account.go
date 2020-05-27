package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/phzietsman/terraform-provider-aws-control-tower/aws/internal/keyvaluetags"
)

func resourceControlTowerAccountVending() *schema.Resource {
	return &schema.Resource{
		Create: resourceControlTowerAccountVendingCreate,
		Read:   resourceControlTowerAccountVendingRead,
		Update: resourceControlTowerAccountVendingUpdate,
		Delete: resourceControlTowerAccountVendingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"artefact_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceControlTowerAccountVendingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).scconn
	input := servicecatalog.ProvisionProductInput{

		AcceptLanguage: aws.String("en"),
		ProvisionToken: aws.String(resource.UniqueId()),

		NotificationArns: []*string{},

		ProductId:              aws.String(d.Get("product_id").(string)),
		ProvisioningArtifactId: aws.String(d.Get("artefact_id").(string)),
		ProvisionedProductName: aws.String(d.Get("name").(string)),

		Tags: keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().ServicecatalogTags(),
	}

	value := keyvaluetags.New(d.Get("parameters").(map[string]interface{}))
	provisioningParameter := make([]*servicecatalog.ProvisioningParameter, len(value))
	index := 0
	for k, v := range value {
		log.Printf("[DEBUG] ProvisioningParameter[%s] Key:%s Value:%s", string(index), k, *v)
		strv := *v
		strk := k
		provisioningParameter[index] = &servicecatalog.ProvisioningParameter{
			Key:   &strk,
			Value: &strv,
		}

		index++
	}

	input.ProvisioningParameters = provisioningParameter

	log.Printf("[DEBUG] Creating Service Catalog Portfolio: %#v", input)
	resp, err := conn.ProvisionProduct(&input)
	if err != nil {
		return fmt.Errorf("Creating Service Catalog Portfolio failed: %s", err.Error())
	}

	recordDetail := *resp.RecordDetail

	d.SetId(*recordDetail.ProvisionedProductId)
	log.Printf("[INFO] Provisioned Product Id: %s", d.Id())

	d.Set("record_id", *recordDetail.RecordId)

	// Wait for the Provisioned Product to become available
	log.Printf("[DEBUG] Waiting for Provisioned Product (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATED", "IN_PROGRESS"},
		Target:  []string{"SUCCEEDED"},
		Refresh: ProvisionedProductStateRefreshFunc(conn, *recordDetail.RecordId),
		Timeout: 60 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Provisioned Product (%s) to become available: %s", d.Id(), err)
	}

	return resourceControlTowerAccountVendingRead(d, meta)
}

func resourceControlTowerAccountVendingRead(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).scconn
	// input := servicecatalog.DescribePortfolioInput{
	// 	AcceptLanguage: aws.String("en"),
	// }
	// input.Id = aws.String(d.Id())

	// log.Printf("[DEBUG] Reading Service Catalog Portfolio: %#v", input)
	// resp, err := conn.DescribePortfolio(&input)
	// if err != nil {
	// 	if scErr, ok := err.(awserr.Error); ok && scErr.Code() == "ResourceNotFoundException" {
	// 		log.Printf("[WARN] Service Catalog Portfolio %q not found, removing from state", d.Id())
	// 		d.SetId("")
	// 		return nil
	// 	}
	// 	return fmt.Errorf("Reading ServiceCatalog Portfolio '%s' failed: %s", *input.Id, err.Error())
	// }
	// portfolioDetail := resp.PortfolioDetail
	// if err := d.Set("created_time", portfolioDetail.CreatedTime.Format(time.RFC3339)); err != nil {
	// 	log.Printf("[DEBUG] Error setting created_time: %s", err)
	// }
	// d.Set("arn", portfolioDetail.ARN)
	// d.Set("description", portfolioDetail.Description)
	// d.Set("name", portfolioDetail.DisplayName)
	// d.Set("provider_name", portfolioDetail.ProviderName)

	// if err := d.Set("tags", keyvaluetags.ServicecatalogKeyValueTags(resp.Tags).IgnoreAws().Map()); err != nil {
	// 	return fmt.Errorf("error setting tags: %s", err)
	// }

	return nil
}

func resourceControlTowerAccountVendingUpdate(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).scconn
	// input := servicecatalog.UpdatePortfolioInput{
	// 	AcceptLanguage: aws.String("en"),
	// 	Id:             aws.String(d.Id()),
	// }

	// if d.HasChange("name") {
	// 	v, _ := d.GetOk("name")
	// 	input.DisplayName = aws.String(v.(string))
	// }

	// if d.HasChange("accept_language") {
	// 	v, _ := d.GetOk("accept_language")
	// 	input.AcceptLanguage = aws.String(v.(string))
	// }

	// if d.HasChange("description") {
	// 	v, _ := d.GetOk("description")
	// 	input.Description = aws.String(v.(string))
	// }

	// if d.HasChange("provider_name") {
	// 	v, _ := d.GetOk("provider_name")
	// 	input.ProviderName = aws.String(v.(string))
	// }

	// if d.HasChange("tags") {
	// 	o, n := d.GetChange("tags")

	// 	input.AddTags = keyvaluetags.New(n).IgnoreAws().ServicecatalogTags()
	// 	input.RemoveTags = aws.StringSlice(keyvaluetags.New(o).IgnoreAws().Keys())
	// }

	// log.Printf("[DEBUG] Update Service Catalog Portfolio: %#v", input)
	// _, err := conn.UpdatePortfolio(&input)
	// if err != nil {
	// 	return fmt.Errorf("Updating Service Catalog Portfolio '%s' failed: %s", *input.Id, err.Error())
	// }
	// return resourceControlTowerAccountVendingRead(d, meta)

	return nil
}

func resourceControlTowerAccountVendingDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).scconn
	input := servicecatalog.TerminateProvisionedProductInput{
		ProvisionedProductId: aws.String(d.Id()),
		TerminateToken:       aws.String(resource.UniqueId()),
	}

	log.Printf("[DEBUG] Delete Vended Account (Service Catalog Provisioned Product): %#v", input)
	resp, err := conn.TerminateProvisionedProduct(&input)
	if err != nil {
		return fmt.Errorf("Deleting Vended Account (Service Catalog Provisioned Product) '%s' failed: %s", *input.ProvisionedProductId, err.Error())
	}

	// Wait for the Provisioned Product to become available
	log.Printf("[DEBUG] Waiting for Provisioned Product (%s) to become terminated", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATED", "IN_PROGRESS"},
		Target:  []string{"SUCCEEDED"},
		Refresh: ProvisionedProductStateRefreshFunc(conn, *resp.RecordDetail.RecordId),
		Timeout: 60 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Provisioned Product (%s) to become available: %s", d.Id(), err)
	}
	return nil
}

// ProvisionedProductStateRefreshFunc returns a resource.StateRefreshFunc
// that is used to watch a Provisioned Product.
func ProvisionedProductStateRefreshFunc(conn *servicecatalog.ServiceCatalog, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		opts := &servicecatalog.DescribeRecordInput{
			Id: aws.String(id),
		}
		resp, err := conn.DescribeRecord(opts)
		if err != nil {
			if isAWSErr(err, "ResourceNotFoundException", "") {
				resp = nil
			} else {
				log.Printf("Error on ProvisionedProductStateRefreshFunc: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// the resource yet. Return an empty state.
			return nil, "", nil
		}

		rec := resp.RecordDetail
		return rec, *rec.Status, nil
	}
}
