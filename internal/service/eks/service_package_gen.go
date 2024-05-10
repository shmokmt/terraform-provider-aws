// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package eks

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	eks_sdkv2 "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newPodIdentityAssociationResource,
			Name:    "Pod Identity Association",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "association_arn",
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceAccessEntry,
			TypeName: "aws_eks_access_entry",
			Name:     "Access Entry",
		},
		{
			Factory:  dataSourceAddon,
			TypeName: "aws_eks_addon",
		},
		{
			Factory:  dataSourceAddonVersion,
			TypeName: "aws_eks_addon_version",
		},
		{
			Factory:  dataSourceCluster,
			TypeName: "aws_eks_cluster",
		},
		{
			Factory:  dataSourceClusterAuth,
			TypeName: "aws_eks_cluster_auth",
		},
		{
			Factory:  dataSourceClusters,
			TypeName: "aws_eks_clusters",
		},
		{
			Factory:  dataSourceNodeGroup,
			TypeName: "aws_eks_node_group",
		},
		{
			Factory:  dataSourceNodeGroups,
			TypeName: "aws_eks_node_groups",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceAccessEntry,
			TypeName: "aws_eks_access_entry",
			Name:     "Access Entry",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "access_entry_arn",
			},
		},
		{
			Factory:  resourceAccessPolicyAssociation,
			TypeName: "aws_eks_access_policy_association",
			Name:     "Access Policy Association",
		},
		{
			Factory:  resourceAddon,
			TypeName: "aws_eks_addon",
			Name:     "Add-On",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceCluster,
			TypeName: "aws_eks_cluster",
			Name:     "Cluster",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceFargateProfile,
			TypeName: "aws_eks_fargate_profile",
			Name:     "Fargate Profile",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceIdentityProviderConfig,
			TypeName: "aws_eks_identity_provider_config",
			Name:     "Identity Provider Config",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceNodeGroup,
			TypeName: "aws_eks_node_group",
			Name:     "Node Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.EKS
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*eks_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return eks_sdkv2.NewFromConfig(cfg, func(o *eks_sdkv2.Options) {
		if endpoint := config[names.AttrEndpoint].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
