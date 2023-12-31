package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	xec2 "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	xecs "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		vpc, err := xec2.NewVpc(ctx, "vpc", nil)
		if err != nil {
			return err
		}
		securityGroup, err := ec2.NewSecurityGroup(ctx, "securityGroup", &ec2.SecurityGroupArgs{
			VpcId: vpc.VpcId,
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					Protocol: pulumi.String("-1"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					Ipv6CidrBlocks: pulumi.StringArray{
						pulumi.String("::/0"),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		cluster, err := ecs.NewCluster(ctx, "cluster", nil)
		if err != nil {
			return err
		}
		_, err = xecs.NewFargateService(ctx, "service", &xecs.FargateServiceArgs{
			Cluster: cluster.Arn,
			NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
				Subnets: vpc.PrivateSubnetIds,
				SecurityGroups: pulumi.StringArray{
					securityGroup.ID(),
				},
			},
			DesiredCount: pulumi.Int(2),
			TaskDefinitionArgs: &xecs.FargateServiceTaskDefinitionArgs{
				Container: &xecs.TaskDefinitionContainerDefinitionArgs{
					Image:     pulumi.String("nginx:latest"),
					Cpu:       pulumi.Int(512),
					Memory:    pulumi.Int(128),
					Essential: pulumi.Bool(true),
				},
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
