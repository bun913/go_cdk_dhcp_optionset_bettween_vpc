package hub_spoke

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Hub struct {
	scope constructs.Construct
}

func NewHub(scope constructs.Construct) Hub {
	return Hub{
		scope: scope,
	}
}

func (h Hub) CreateTransitGateway() awsec2.CfnTransitGateway {
	return awsec2.NewCfnTransitGateway(h.scope, jsii.String("TransitGateway"), &awsec2.CfnTransitGatewayProps{
		DefaultRouteTableAssociation: jsii.String("disable"),
		DefaultRouteTablePropagation: jsii.String("disable"),
	})
}

type VpcAttachment struct {
	name            string
	vpc             awsec2.Vpc
	tgw             awsec2.CfnTransitGateway
	subnetGroupName string
}

func NewVpcAttachment(name string, vpc awsec2.Vpc, tgw awsec2.CfnTransitGateway, subnetGroupName string) VpcAttachment {
	return VpcAttachment{
		name:            name,
		vpc:             vpc,
		tgw:             tgw,
		subnetGroupName: subnetGroupName,
	}
}

func (va VpcAttachment) Attach() awsec2.CfnTransitGatewayAttachment {
	return awsec2.NewCfnTransitGatewayAttachment(va.vpc, jsii.String("VpcAttachment"), &awsec2.CfnTransitGatewayAttachmentProps{
		SubnetIds:        va.vpc.SelectSubnets(&awsec2.SubnetSelection{SubnetGroupName: &va.subnetGroupName}).SubnetIds,
		TransitGatewayId: va.tgw.Ref(),
		VpcId:            va.vpc.VpcId(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String("Name"),
				Value: jsii.String(va.name),
			},
		},
	})
}
