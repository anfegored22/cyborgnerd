import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as ecrassets from 'aws-cdk-lib/aws-ecr-assets';

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
  }
}
