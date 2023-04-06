package network

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type routeToTransitGateway struct {
	scope         constructs.Construct
	name          string
	vpc           awsec2.Vpc
	tgw           awsec2.CfnTransitGateway
	tgwAttachment awsec2.CfnTransitGatewayAttachment
}

func NewRouteToTransitGateway(scope constructs.Construct, name string, vpc awsec2.Vpc, tgw awsec2.CfnTransitGateway, tgwAttachment awsec2.CfnTransitGatewayAttachment) routeToTransitGateway {
	return routeToTransitGateway{
		scope:         scope,
		name:          name,
		vpc:           vpc,
		tgw:           tgw,
		tgwAttachment: tgwAttachment,
	}
}

func (rttg routeToTransitGateway) CreateRouteToTransitGateway() {
	subnets := rttg.vpc.SelectSubnets(&awsec2.SubnetSelection{
		SubnetGroupName: jsii.String("Private"),
	}).Subnets
	for i, subnet := range *subnets {
		routeName := fmt.Sprintf("%s%d", rttg.name, i)
		awsec2.NewCfnRoute(rttg.scope, jsii.String(routeName), &awsec2.CfnRouteProps{
			RouteTableId:         subnet.RouteTable().RouteTableId(),
			DestinationCidrBlock: jsii.String("0.0.0.0/0"),
			TransitGatewayId:     rttg.tgw.Ref(),
		}).AddDependency(rttg.tgwAttachment)
	}
}
