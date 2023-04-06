package network

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Network struct {
	scope          constructs.Construct
	vpcName        string
	cidr           string
	hasSSMEndpoint bool
}

func NewNetwork(scope constructs.Construct, vpcName string, cidr string, hasSSMEndpoint bool) Network {
	return Network{
		scope:          scope,
		vpcName:        vpcName,
		cidr:           cidr,
		hasSSMEndpoint: hasSSMEndpoint,
	}
}

type VpcResult struct {
	Vpc       awsec2.Vpc
	Endpoints map[EndPointkey]awsec2.InterfaceVpcEndpoint
}

type EndPointkey string

const ssmEndpoint = "ssm.ap-northeast-1.amazonaws.com"
const ssmMessageEndpoint = "ssmmessages.ap-northeast-1.amazonaws.com"
const ec2MeessageEndpoint = "ec2messages.ap-northeast-1.amazonaws.com"

const (
	SsmEndpoint         EndPointkey = ssmEndpoint
	SsmMessageEndpoint  EndPointkey = ssmMessageEndpoint
	Ec2MeessageEndpoint EndPointkey = ec2MeessageEndpoint
)

func (nr Network) CreateNetworkResources() VpcResult {
	// VPC
	vpc := awsec2.NewVpc(nr.scope, &nr.vpcName, &awsec2.VpcProps{
		IpAddresses:        awsec2.IpAddresses_Cidr(jsii.String(nr.cidr)),
		MaxAzs:             jsii.Number(2),
		EnableDnsSupport:   jsii.Bool(true),
		EnableDnsHostnames: jsii.Bool(true),
		VpcName:            jsii.String(nr.vpcName),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("TransitGateway"),
				SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				CidrMask:   jsii.Number(28),
			},
			{
				Name:       jsii.String("Private"),
				SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				CidrMask:   jsii.Number(24),
			},
		},
	})
	// 指定した時のみVPCエンドポイントを追加
	endpoints := map[EndPointkey]awsec2.InterfaceVpcEndpoint{}
	if nr.hasSSMEndpoint {
		sg := awsec2.NewSecurityGroup(nr.scope, jsii.String(nr.vpcName+"SSMSecurityGroup"), &awsec2.SecurityGroupProps{
			Vpc: vpc,
		})
		sg.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("10.0.0.0/8")), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("allow inbound from"), jsii.Bool(false))
		ssmEndpoint := vpc.AddInterfaceEndpoint(jsii.String("SSM"), &awsec2.InterfaceVpcEndpointOptions{
			Service: awsec2.InterfaceVpcEndpointAwsService_SSM(),
			SecurityGroups: &[]awsec2.ISecurityGroup{
				sg,
			},
			PrivateDnsEnabled: jsii.Bool(false),
			Subnets: &awsec2.SubnetSelection{
				SubnetGroupName: jsii.String("Private"),
			},
		})
		endpoints[SsmEndpoint] = ssmEndpoint
		ssmMessageEndpoint := vpc.AddInterfaceEndpoint(jsii.String(nr.vpcName+"SSMMessage"), &awsec2.InterfaceVpcEndpointOptions{

			Service: awsec2.InterfaceVpcEndpointAwsService_SSM_MESSAGES(),
			SecurityGroups: &[]awsec2.ISecurityGroup{
				sg,
			},
			PrivateDnsEnabled: jsii.Bool(false),
			Subnets: &awsec2.SubnetSelection{
				SubnetGroupName: jsii.String("Private"),
			},
		})
		endpoints[SsmMessageEndpoint] = ssmMessageEndpoint
		ec2MEndpoint := vpc.AddInterfaceEndpoint(jsii.String(nr.vpcName+"EC2Messag"), &awsec2.InterfaceVpcEndpointOptions{
			Service: awsec2.InterfaceVpcEndpointAwsService_EC2_MESSAGES(),
			SecurityGroups: &[]awsec2.ISecurityGroup{
				sg,
			},
			PrivateDnsEnabled: jsii.Bool(false),
			Subnets: &awsec2.SubnetSelection{
				SubnetGroupName: jsii.String("Private"),
			},
		})
		endpoints[ec2MeessageEndpoint] = ec2MEndpoint
	}
	return VpcResult{
		Vpc:       vpc,
		Endpoints: endpoints,
	}
}
