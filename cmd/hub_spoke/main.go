package hub_spoke

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
)

type HubParameters struct {
	scope     constructs.Construct
	sharedVpc awsec2.Vpc
	hubVpc    awsec2.Vpc
}

func NewHubParameters(scope constructs.Construct, sharedVpc awsec2.Vpc, hubVpc awsec2.Vpc) HubParameters {
	return HubParameters{
		scope:     scope,
		sharedVpc: sharedVpc,
		hubVpc:    hubVpc,
	}
}

type HubResult struct {
	Tgw             awsec2.CfnTransitGateway
	HubAttachment   awsec2.CfnTransitGatewayAttachment
	SpokeAttachment awsec2.CfnTransitGatewayAttachment
}

func (hp HubParameters) CreateHubResources() HubResult {
	hub := NewHub(hp.scope)
	// Transit Gatewayを作成
	tgw := hub.CreateTransitGateway()
	// Transit GatewayにVPCをアタッチ
	attachmentShared := NewVpcAttachment("HubVpcAttachment", hp.sharedVpc, tgw, "TransitGateway")
	attachmentSharedVpc := attachmentShared.Attach()
	attchmentWorkload := NewVpcAttachment("SpokeVpcAttachment", hp.hubVpc, tgw, "TransitGateway")
	attachmentWorkloadVpc := attchmentWorkload.Attach()
	// Transit Gatewayのルートテーブル作成
	rt := NewRouteTable("RouteTable", tgw)
	routeTable := rt.Create()
	// RouteTableにルート追加・VPCアタッチメントのアソシエーション
	hubVpcRouteAssociation := NewVpcRouteAssociation("HubVpcAssocation", attachmentSharedVpc, routeTable)
	hubVpcRouteAssociation.Create()
	spokeVpcRouteAssociation := NewVpcRouteAssociation("SpokeVpcAssociation", attachmentWorkloadVpc, routeTable)
	spokeVpcRouteAssociation.Create()
	// hubとspokeVPCの双方向のルートをルートテーブルに追加
	connectVpcs := NewVpcsConnection(hp.sharedVpc, attachmentSharedVpc, hp.hubVpc, attachmentWorkloadVpc, routeTable)
	connectVpcs.Create()
	return HubResult{
		Tgw:             tgw,
		HubAttachment:   attachmentSharedVpc,
		SpokeAttachment: attachmentWorkloadVpc,
	}
}
