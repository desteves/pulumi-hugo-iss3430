package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecr"
	xecs "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		repository, err := ecr.NewRepository(ctx, "repository", nil)
		if err != nil {
			return err
		}

		image, err := ecr.NewImage(ctx, "image", &ecr.ImageArgs{
			RepositoryUrl: repository.Url,
			Path:          pulumi.String("./app"),
		})
		if err != nil {
			return err
		}

		cluster, err := ecs.NewCluster(ctx, "cluster", nil)
		if err != nil {
			return err
		}

		lb, err := lb.NewApplicationLoadBalancer(ctx, "lb", nil)
		if err != nil {
			return err
		}

		tdpma := xecs.TaskDefinitionPortMappingArray{
			xecs.TaskDefinitionPortMappingArgs{
				TargetGroup: lb.DefaultTargetGroup,
			},
		}
		_, err = xecs.NewFargateService(ctx, "service", &xecs.FargateServiceArgs{
			Cluster:        cluster.Arn,
			AssignPublicIp: pulumi.Bool(true),
			TaskDefinitionArgs: &xecs.FargateServiceTaskDefinitionArgs{
				Container: &xecs.TaskDefinitionContainerDefinitionArgs{
					Image:        image.ImageUri,
					Cpu:          pulumi.Int(512),
					Memory:       pulumi.Int(128),
					Essential:    pulumi.Bool(true),
					PortMappings: tdpma,
				},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("url", lb.LoadBalancer.DnsName())
		return nil
	})
}
