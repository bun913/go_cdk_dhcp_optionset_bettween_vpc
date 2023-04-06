package resolver

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53resolver"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ResolverEndpoint struct {
	scope       constructs.Construct
	name        string
	ipAddresses []*string
	vpc         awsec2.Vpc
}

func NewResolverEndpoint(scope constructs.Construct, name string, ipAddresses []*string, vpc awsec2.Vpc) *ResolverEndpoint {
	return &ResolverEndpoint{
		scope:       scope,
		name:        name,
		ipAddresses: ipAddresses,
		vpc:         vpc,
	}
}

// CreateResolverEndpoint creates a resolver endpoint in the VPC
func (re ResolverEndpoint) CreateResolverEndpoint() awsroute53resolver.CfnResolverEndpoint {
	return awsroute53resolver.NewCfnResolverEndpoint(re.scope, jsii.String(re.name), &awsroute53resolver.CfnResolverEndpointProps{
		Direction:        jsii.String("INBOUND"),
		IpAddresses:      re.getInboundSubnets(),
		Name:             jsii.String(re.name),
		SecurityGroupIds: &[]*string{re.createSecurityGroup().SecurityGroupId()},
	})
}

// create security group for the resolver endpoint
func (re ResolverEndpoint) createSecurityGroup() awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(re.scope, jsii.String(re.name+"SecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc: re.vpc,
	})
	sg.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")), awsec2.Port_Tcp(jsii.Number(53)), jsii.String("Allow inbound DNS"), jsii.Bool(false))
	sg.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")), awsec2.Port_Udp(jsii.Number(53)), jsii.String("Allow inbound DNS"), jsii.Bool(false))
	return sg
}

func (re ResolverEndpoint) getInboundSubnets() []*awsroute53resolver.CfnResolverEndpoint_IpAddressRequestProperty {
	subnets := re.vpc.SelectSubnets(&awsec2.SubnetSelection{
		SubnetGroupName: jsii.String("Private"),
	}).Subnets
	ipAddresses := make([]*awsroute53resolver.CfnResolverEndpoint_IpAddressRequestProperty, len(*subnets))
	for i, subnet := range *subnets {
		ipAddresses[i] = &awsroute53resolver.CfnResolverEndpoint_IpAddressRequestProperty{
			SubnetId: subnet.SubnetId(),
			Ip:       re.ipAddresses[i],
		}
	}
	return ipAddresses
}
