package main

import (
	"multi_vpc_resolver/cmd/hub_spoke"
	"multi_vpc_resolver/cmd/network"
	"multi_vpc_resolver/cmd/resolver"
	"multi_vpc_resolver/cmd/server"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VpcAndEc2StackProps struct {
	awscdk.StackProps
}

func NewVpcAndEc2Stack(scope constructs.Construct, id string, props *VpcAndEc2StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	// 共通タグを設定
	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String("Sample"), nil)
	awscdk.Tags_Of(app).Add(jsii.String("Env"), jsii.String("Dev"), nil)

	stack := NewVpcAndEc2Stack(app, "Sample", &VpcAndEc2StackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	// Hub(Shared)VPC関連リソース
	sharedNetworkResource := network.NewNetwork(stack, "SharedVpc", "10.10.0.0/16", true)
	sharedVpcResult := sharedNetworkResource.CreateNetworkResources()
	sharedVpc := sharedVpcResult.Vpc
	severResource := server.NewServer(stack, "SharedVPCInstance", sharedVpc)
	severResource.CreateServerResources()

	// WorkloadVPC関連リソース
	workloadNetwork := network.NewNetwork(stack, "WorkloadVpc", "10.20.0.0/16", false)
	workloadVpcResult := workloadNetwork.CreateNetworkResources()
	workloadVpc := workloadVpcResult.Vpc
	workloadServer := server.NewServer(stack, "WorkloadVPCInstance", workloadVpc)
	workloadServer.CreateServerResources()

	// SharedVpcにRoute53 Private Hosted Zoneを作成
	endpoints := sharedVpcResult.Endpoints
	for key, endpoint := range endpoints {
		name := string(key)
		hostedZoneResource := network.NewHostedZone(stack, name, endpoint, sharedVpc)
		zone := hostedZoneResource.CreateHostedZone()
		hostedZoneResource.AddAliasRecord(zone)
	}

	// Transit GatewayによるVPC間の双方向通信を実現するためのリソースを作成
	hubParameters := hub_spoke.NewHubParameters(stack, sharedVpc, workloadVpc)
	hubResult := hubParameters.CreateHubResources()

	// EC2が属するサブネットのルートテーブルからTransit Gatewayへのルートを追加
	routeHubSubnetToTransit := network.NewRouteToTransitGateway(stack, "HubSubnetToTransitGW", sharedVpc, hubResult.Tgw, hubResult.HubAttachment)
	routeHubSubnetToTransit.CreateRouteToTransitGateway()
	routeSpokeSubnetToTransit := network.NewRouteToTransitGateway(stack, "SpokeSubnetToTransitGW", workloadVpc, hubResult.Tgw, hubResult.SpokeAttachment)
	routeSpokeSubnetToTransit.CreateRouteToTransitGateway()

	// Route53 Resolver関連のリソースを作成
	endPointAddresses := []*string{jsii.String("10.10.1.100"), jsii.String("10.10.2.100")}
	rv := resolver.NewResolverEndpoint(stack, "ResolverEndpoint", endPointAddresses, sharedVpc)
	endpoint := rv.CreateResolverEndpoint()
	// WorkloadVPCのDHCPオプションセットでRoute53 Resolver Endpointに名前解決を向ける
	do := resolver.NewDhcpOptionSet(stack, "DhcpOptions", endpoint, endPointAddresses, workloadVpc)
	do.CreateDHCPOptionSet()

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String("ap-northeast-1"),
	}
}
