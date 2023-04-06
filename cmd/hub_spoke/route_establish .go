package hub_spoke

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

// VPCアタッチメントをTransitGatewayのルートテーブルと関連付け
type VpcRouteAssociation struct {
	name          string
	vpcAttachment awsec2.CfnTransitGatewayAttachment
	routeTable    awsec2.CfnTransitGatewayRouteTable
}

func NewVpcRouteAssociation(name string, vpcAttachment awsec2.CfnTransitGatewayAttachment, routeTable awsec2.CfnTransitGatewayRouteTable) VpcRouteAssociation {
	return VpcRouteAssociation{
		name:          name,
		vpcAttachment: vpcAttachment,
		routeTable:    routeTable,
	}
}
func (vra VpcRouteAssociation) Create() { // Association
	awsec2.NewCfnTransitGatewayRouteTableAssociation(vra.vpcAttachment, jsii.String(vra.name+"Association"), &awsec2.CfnTransitGatewayRouteTableAssociationProps{
		TransitGatewayAttachmentId: vra.vpcAttachment.Ref(),
		TransitGatewayRouteTableId: vra.routeTable.Ref(),
	})
	// Propagation
	awsec2.NewCfnTransitGatewayRouteTablePropagation(vra.vpcAttachment, jsii.String(vra.name+"Propagation"), &awsec2.CfnTransitGatewayRouteTablePropagationProps{
		TransitGatewayAttachmentId: vra.vpcAttachment.Ref(),
		TransitGatewayRouteTableId: vra.routeTable.Ref(),
	})
}

// HubVPCとSpokeVPCの双方向のルートをルートテーブルに追加
type VpcsConnection struct {
	hubVpc             awsec2.Vpc
	hubVpcAttachment   awsec2.CfnTransitGatewayAttachment
	spokeVpc           awsec2.Vpc
	spokeVpcAttachment awsec2.CfnTransitGatewayAttachment
	routetable         awsec2.CfnTransitGatewayRouteTable
}

func NewVpcsConnection(hubVpc awsec2.Vpc, hubVpcAttachment awsec2.CfnTransitGatewayAttachment, spokeVpc awsec2.Vpc, spokeVpcAttachment awsec2.CfnTransitGatewayAttachment, routetable awsec2.CfnTransitGatewayRouteTable) VpcsConnection {
	return VpcsConnection{
		hubVpc:             hubVpc,
		hubVpcAttachment:   hubVpcAttachment,
		spokeVpc:           spokeVpc,
		spokeVpcAttachment: spokeVpcAttachment,
		routetable:         routetable,
	}
}

func (cv VpcsConnection) Create() {
	awsec2.NewCfnTransitGatewayRoute(cv.spokeVpc, jsii.String("ToSpokeVpc"), &awsec2.CfnTransitGatewayRouteProps{
		DestinationCidrBlock:       cv.spokeVpc.VpcCidrBlock(),
		TransitGatewayAttachmentId: cv.spokeVpcAttachment.Ref(),
		TransitGatewayRouteTableId: cv.routetable.Ref(),
	})
	awsec2.NewCfnTransitGatewayRoute(cv.hubVpc, jsii.String("ToHubVpc"), &awsec2.CfnTransitGatewayRouteProps{
		DestinationCidrBlock:       cv.hubVpc.VpcCidrBlock(),
		TransitGatewayAttachmentId: cv.hubVpcAttachment.Ref(),
		TransitGatewayRouteTableId: cv.routetable.Ref(),
	})
}
