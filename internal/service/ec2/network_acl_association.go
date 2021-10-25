package ec2

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceNetworkAclAssociation() *schema.Resource {
	return &schema.Resource{
		Create: ResourceNetworkAclAssociationCreate,
		Read:   ResourceNetworkAclAssociationRead,
		Delete: ResourceNetworkAclAssociationDelete,

		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceNetworkAclAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	naclId := d.Get("network_acl_id").(string)
	subnetId := d.Get("subnet_id").(string)

	association, errAssociation := findNetworkAclAssociation(subnetId, conn)
	if errAssociation != nil {
		return fmt.Errorf("Failed to find association for subnet %s: %s", subnetId, errAssociation)
	}

	associationOpts := &ec2.ReplaceNetworkAclAssociationInput{
		AssociationId: association.NetworkAclAssociationId,
		NetworkAclId:  aws.String(naclId),
	}
	log.Printf("[DEBUG] Creating Network ACL association: %#v", associationOpts)

	resp, err := conn.ReplaceNetworkAclAssociation(associationOpts)
	if err != nil {
		return fmt.Errorf("Error creating network acl association: %w", err)
	}

	associationId := resp.NewAssociationId
	d.SetId(aws.StringValue(associationId))
	log.Printf("[INFO] New Association ID: %s", d.Id())

	return ResourceNetworkAclAssociationRead(d, meta)
}

func ResourceNetworkAclAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	// Inspect that the association exists
	subnetId := d.Get("subnet_id").(string)
	association, err := findNetworkAclAssociation(subnetId, conn)
	if err != nil {
		if tfresource.NotFound(err) {
			log.Printf("[WARN] Unable to find association for subnet %s", subnetId)
			d.SetId("")
			return nil
		}
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr != nil {
				log.Printf("[WARN] Unable to find association for subnet %s", subnetId)
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("network_acl_id", aws.StringValue(association.NetworkAclId))

	return nil
}

func ResourceNetworkAclAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	subnetId := d.Get("subnet_id").(string)

	req := &ec2.DescribeNetworkAclsInput{}
	req.Filters = BuildAttributeFilterList(
		map[string]string{
			"association.subnet-id": subnetId,
		},
	)

	resp, err := conn.DescribeNetworkAcls(req)
	if err != nil {
		return err
	}

	if len(resp.NetworkAcls) == 0 {
		return fmt.Errorf("Unable to find Network ACL for subnet: %s", subnetId)
	}

	nacl := resp.NetworkAcls[0]
	var association *ec2.NetworkAclAssociation
	if len(resp.NetworkAcls) > 0 {
		for _, assoc := range nacl.Associations {
			if aws.StringValue(assoc.SubnetId) == subnetId {
				association = assoc
			}
		}
	}

	defaultAcl, err := GetDefaultNetworkACL(*nacl.VpcId, conn)

	if err != nil {
		return fmt.Errorf("Failed to get default Network Acl : %s", err)
	}

	associationOpts := ec2.ReplaceNetworkAclAssociationInput{
		AssociationId: association.NetworkAclAssociationId,
		NetworkAclId:  defaultAcl.NetworkAclId,
	}

	log.Printf("[DEBUG] Replacing Network ACL association: %#v", associationOpts)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.ReplaceNetworkAclAssociation(&associationOpts)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr != nil {
					return resource.RetryableError(awsErr)
				}
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.ReplaceNetworkAclAssociation(&associationOpts)
	}
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
