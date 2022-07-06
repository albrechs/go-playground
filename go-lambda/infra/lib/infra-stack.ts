import * as cdk from "@aws-cdk/core";
import * as path from "path";
import * as lambda from "@aws-cdk/aws-lambda";
import * as apigwv2 from "@aws-cdk/aws-apigatewayv2-alpha";

export class GoLambdaStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
    const lambdaFunction = this.buildAndInstallGoApp(
      "http-test",
      path.join(__dirname, "../../app"),
      "main"
    );
    const httpGateway = this.createHttpGatewayForLambda(
      "http-test",
      lambdaFunction
    );
    new cdk.CfnOutput(this, "lambda-url", { value: apiGtw.url! });
  }
  // build and install go function
  buildAndInstallGoApp(
    id: string,
    lambdaPath: string,
    handler: string
  ): lambda.Function {
    const environment = {
      CGO_ENABLED: "0",
      GOOS: "linux",
      GOARCH: "amd64",
    };
    return new lambda.Function(this, id, {
      code: lambda.Code.fromAsset(lambdaPath, {
        bundling: {
          image: lambda.Runtime.GO_1_X.bundlingDockerImage,
          user: "root",
          environment,
          command: [
            "bash",
            "-c",
            ["make vendor", "make lambda-build"].join(" && "),
          ],
        },
      }),
      handler,
      runtime: lambda.Runtime.GO_1_X,
    });
  }

  createHttpGatewayForLambda(
    id: string,
    handler: lambda.Function
  ): apigwv2.HttpApi {
    // apigw domain
    domain = apigwv2.DomainName.fromDomainNameAttributes(this, id, {
      name: "spotdash.albrechs.dev",
      regionalDomainName: "d-qw3b6mdexh.execute-api.us-east-1.amazonaws.com",
      regionalHostedZoneId: "Z1UJRXOUMOOFQ8",
    });
    // httpgw
    gateway = new apigwv2.HttpApi(this, id, {
      defaultDomainMapping: {
        domainName: domain,
        mappingKey: id,
      },
    });
    // htppgw route
    gateway.addRoutes({
      path: "/",
      methods: [apigwv2.HttpMethod.GET],
      integration: handler,
    });
  }
}
