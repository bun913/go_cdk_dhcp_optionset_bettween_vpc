package resolver

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53resolver"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// this file has a struct for creating dhcp option set

type DHCPOptionSet struct {
	scope           constructs.Construct
	inboundEndpoint awsroute53resolver.CfnResolverEndpoint
	name            string
	ipaddresses     []*string
	vpc             awsec2.Vpc
}

func NewDhcpOptionSet(scope constructs.Construct, name string, inboundEndpoint awsroute53resolver.CfnResolverEndpoint, ipAddresses []*string, vpc awsec2.Vpc) *DHCPOptionSet {
	return &DHCPOptionSet{
		scope:           scope,
		inboundEndpoint: inboundEndpoint,
		name:            name,
		ipaddresses:     ipAddresses,
		vpc:             vpc,
	}
}

func (d DHCPOptionSet) CreateDHCPOptionSet() awsec2.CfnDHCPOptions {
	options := awsec2.NewCfnDHCPOptions(d.scope, jsii.String(d.name), &awsec2.CfnDHCPOptionsProps{
		DomainNameServers: &d.ipaddresses,
	})
	awsec2.NewCfnVPCDHCPOptionsAssociation(d.scope, jsii.String(d.name+"Association"), &awsec2.CfnVPCDHCPOptionsAssociationProps{
		VpcId:         d.vpc.VpcId(),
		DhcpOptionsId: options.Ref(),
	})
	return options
}
