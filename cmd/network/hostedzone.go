package network

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	targets "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HostedZone struct {
	scope       constructs.Construct
	name        string
	vpcEndpoint awsec2.InterfaceVpcEndpoint
	vpc         awsec2.Vpc
}

func NewHostedZone(scope constructs.Construct, name string, vpcEndpoint awsec2.InterfaceVpcEndpoint, vpc awsec2.Vpc) HostedZone {
	return HostedZone{
		scope:       scope,
		name:        name,
		vpcEndpoint: vpcEndpoint,
		vpc:         vpc,
	}
}

func (hz HostedZone) CreateHostedZone() awsroute53.PrivateHostedZone {
	zone := awsroute53.NewPrivateHostedZone(hz.scope, jsii.String(hz.name), &awsroute53.PrivateHostedZoneProps{
		Vpc:      hz.vpc,
		ZoneName: jsii.String(hz.name),
	})
	return zone
}

func (hz HostedZone) AddAliasRecord(zone awsroute53.PrivateHostedZone) awsroute53.ARecord {
	return awsroute53.NewARecord(hz.scope, jsii.String(hz.name+"alias"), &awsroute53.ARecordProps{
		Zone:   zone,
		Target: awsroute53.RecordTarget_FromAlias(targets.NewInterfaceVpcEndpointTarget(hz.vpcEndpoint)),
	})
}
