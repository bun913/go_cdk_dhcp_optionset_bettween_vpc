package hub_spoke

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

type RouteTable struct {
	name string
	tgw  awsec2.CfnTransitGateway
}

func NewRouteTable(name string, tgw awsec2.CfnTransitGateway) RouteTable {
	return RouteTable{
		name: name,
		tgw:  tgw,
	}
}

func (ra RouteTable) Create() awsec2.CfnTransitGatewayRouteTable {
	return awsec2.NewCfnTransitGatewayRouteTable(ra.tgw, jsii.String("RouteTable"), &awsec2.CfnTransitGatewayRouteTableProps{
		TransitGatewayId: ra.tgw.Ref(),
		Tags: &[]*awscdk.CfnTag{
			{
				Key:   jsii.String("Name"),
				Value: jsii.String(ra.name),
			},
		},
	})
}
