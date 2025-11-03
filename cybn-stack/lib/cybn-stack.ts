import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as ecrassets from 'aws-cdk-lib/aws-ecr-assets';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ecs_patterns from 'aws-cdk-lib/aws-ecs-patterns';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';

export class CybnStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Build and publish a Docker image asset to ECR (CDK-managed assets repo)
    const imageAsset = new ecrassets.DockerImageAsset(this, 'CybnAppImage', {
      directory: path.join(__dirname, '../../src'),
      file: 'Dockerfile',
    });

    // Optionally, provision a named ECR repository for CI/CD pushes
    const appRepo = new ecr.Repository(this, 'CybnAppRepository', {
      repositoryName: 'cybn-app',
      imageScanOnPush: true,
      imageTagMutability: ecr.TagMutability.IMMUTABLE,
      lifecycleRules: [
        { maxImageCount: 20 },
      ],
    });

    // Useful outputs
    new cdk.CfnOutput(this, 'AssetImageUri', {
      value: imageAsset.imageUri,
      description: 'URI of the Docker image uploaded by CDK assets.',
    });

    new cdk.CfnOutput(this, 'EcrRepositoryUri', {
      value: appRepo.repositoryUri,
      description: 'URI of the named ECR repository for CI/CD pushes.',
    });

    // --- ECS on Fargate with a public ALB ---
    // Network: small VPC across up to 2 AZs (defaults are fine for a simple service)
    const vpc = new ec2.Vpc(this, 'CybnVpc', {
      maxAzs: 2,
    });

    // ECS Cluster
    const cluster = new ecs.Cluster(this, 'CybnCluster', {
      vpc,
      containerInsights: true,
    });

    // Fargate Service behind an Application Load Balancer
    // Create a Route 53 public hosted zone for cyborgnerd.com
    // Note: you still need to set these NS at your registrar (Namecheap)
    const hostedZone = new route53.PublicHostedZone(this, 'CybnHostedZone', {
      zoneName: 'cyborgnerd.com',
      comment: 'Managed by CDK for cyborgnerd.com',
    });

    // Create an ACM certificate for the apex domain, validated via DNS in the hosted zone
    const certificate = new acm.Certificate(this, 'CybnApexCertificate', {
      domainName: 'cyborgnerd.com',
      validation: acm.CertificateValidation.fromDns(hostedZone),
    });

    const fargateService = new ecs_patterns.ApplicationLoadBalancedFargateService(
      this,
      'CybnFargateService',
      {
        cluster,
        publicLoadBalancer: true,
        desiredCount: 1,
        cpu: 256,
        memoryLimitMiB: 512,
        healthCheckGracePeriod: cdk.Duration.seconds(30),
        domainName: 'cyborgnerd.com',
        domainZone: hostedZone,
        certificate,
        redirectHTTP: true,
        taskImageOptions: {
          // Use the image we just built as a CDK asset
          image: ecs.ContainerImage.fromDockerImageAsset(imageAsset),
          containerPort: 8080,
          enableLogging: true,
        },
      }
    );

    // Make sure health check path is reasonable for the app
    fargateService.targetGroup.configureHealthCheck({
      path: '/',
      healthyHttpCodes: '200-399',
    });

    // Output the public URL of the load balancer
    new cdk.CfnOutput(this, 'ServiceUrl', {
      value: `https://cyborgnerd.com`,
      description: 'Public URL for the Fargate service via ALB and custom domain.',
    });

    // Output NS records to configure at Namecheap
    new cdk.CfnOutput(this, 'Route53NameServers', {
      value: cdk.Fn.join(',', hostedZone.hostedZoneNameServers ?? []),
      description: 'Set these nameservers at your domain registrar (Namecheap).',
    });
  }
}
