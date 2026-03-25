import * as pulumi from "@pulumi/pulumi";
import * as awsNative from "@pulumi/aws-native";
import * as aws from "@pulumi/aws";

const config = new pulumi.Config();

// ── Parameters (were CloudFormation Parameters) ───────────────────────────────
// Secret values are injected by CI via `pulumi config set --secret`.
const openAIApiKey    = config.getSecret("openAIApiKey")        ?? "";
const geminiApiKey    = config.getSecret("geminiApiKey")        ?? "";
const vertexApiKey    = config.getSecret("vertexApiKey")        ?? "";
const cfAccountId     = config.getSecret("cloudFlareAccountId") ?? "";
const cfApiKey        = config.getSecret("cloudFlareAPIKey")    ?? "";
const anthropicApiKey = config.getSecret("anthropicAPIKey")     ?? "";
const runtime         = config.get("runtime")                   ?? "provided.al2023";
const architecture    = config.get("architecture")              ?? "arm64";
const handler         = config.get("handler")                   ?? "bootstrap";

// Matches the CloudFormation stack name used by sam deploy in CI ("chat-gpt")
// so physical resource names stay identical to the existing stack.
const cfnStackName = "chat-gpt";

// ── S3 Bucket ─────────────────────────────────────────────────────────────────
// CloudFormation Logical ID: Bucket  |  Type: AWS::S3::Bucket
const Bucket = new awsNative.s3.Bucket("Bucket", {
    bucketName: `${cfnStackName}-contents`,
    corsConfiguration: {
        corsRules: [{
            allowedHeaders: ["*"],
            allowedMethods: ["GET"],
            allowedOrigins: [
                "http://ask-ifr-download.s3.amazonaws.com",
                "https://ask-ifr-download.s3.amazonaws.com",
            ],
        }],
    },
    tags: [{ key: "Project", value: "Alexa-ChatGPT" }],
});

// ── SQS Queues ────────────────────────────────────────────────────────────────
// CloudFormation Logical ID: RequestsDLQ  |  Type: AWS::SQS::Queue
const RequestsDLQ = new awsNative.sqs.Queue("RequestsDLQ", {
    queueName: `${cfnStackName}-Requests-DLQ`,
    tags: [{ key: "Project", value: "Alexa-ChatGPT" }],
});

// CloudFormation Logical ID: RequestsQueue  |  Type: AWS::SQS::Queue
const RequestsQueue = new awsNative.sqs.Queue("RequestsQueue", {
    queueName: `${cfnStackName}-Requests`,
    visibilityTimeout: 301,
    redrivePolicy: {
        deadLetterTargetArn: RequestsDLQ.arn,
        maxReceiveCount: 5,
    },
    tags: [{ key: "Project", value: "Alexa-ChatGPT" }],
});

// CloudFormation Logical ID: ResponsesQueue  |  Type: AWS::SQS::Queue
const ResponsesQueue = new awsNative.sqs.Queue("ResponsesQueue", {
    queueName: `${cfnStackName}-Responses`,
    visibilityTimeout: 301,
    tags: [{ key: "Project", value: "Alexa-ChatGPT" }],
});

// ── OTEL Lambda Layer (us-east-1) ─────────────────────────────────────────────
// From SAM Globals.Function.Layers
const otelLayerArn = "arn:aws:lambda:us-east-1:901920570463:layer:aws-otel-collector-arm64-ver-0-115-0:3";

// ── Shared Lambda environment variables ──────────────────────────────────────
// From SAM Globals.Function.Environment.Variables
const commonEnvVars: { [key: string]: pulumi.Input<string> } = {
    RESPONSES_QUEUE_URI:   ResponsesQueue.queueUrl,
    REQUESTS_QUEUE_URI:    RequestsQueue.queueUrl,
    OPENAI_API_KEY:        openAIApiKey,
    GEMINI_API_KEY:        geminiApiKey,
    VERTEX_API_KEY:        vertexApiKey,
    POLL_DELAY:            "7",
    CLOUDFLARE_ACCOUNT_ID: cfAccountId,
    CLOUDFLARE_API_KEY:    cfApiKey,
    S3_BUCKET:             Bucket.bucketName.apply(n => n ?? ""),
    OPENAI_BASE_URL:       "https://api.openai.com/v1",
    ANTHROPIC_API_KEY:     anthropicApiKey,
};

// ── IAM Execution Role: ChatGPTFunction ───────────────────────────────────────
// SAM auto-generates this from the Policies block. Expanded SAM policy templates:
//   AWSXrayWriteOnlyAccess  → managed policy
//   SQSPollerPolicy         → receive/delete/get/change-visibility on ResponsesQueue
//   SQSSendMessagePolicy    → send on RequestsQueue
//   (inline)                → purge on both queues
const ChatGPTFunctionRole = new awsNative.iam.Role("ChatGPTFunctionRole", {
    roleName: `${cfnStackName}-ChatGPTFunctionRole`,
    assumeRolePolicyDocument: {
        Version: "2012-10-17",
        Statement: [{
            Effect:    "Allow",
            Principal: { Service: "lambda.amazonaws.com" },
            Action:    "sts:AssumeRole",
        }],
    },
    managedPolicyArns: [
        "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
        "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
    ],
    policies: [{
        policyName: "ChatGPTFunctionPolicy",
        policyDocument: pulumi.all([RequestsQueue.arn, ResponsesQueue.arn]).apply(
            ([reqArn, resArn]) => ({
                Version: "2012-10-17",
                Statement: [
                    {
                        Effect: "Allow",
                        Action: [
                            "sqs:ReceiveMessage",
                            "sqs:DeleteMessage",
                            "sqs:GetQueueAttributes",
                            "sqs:ChangeMessageVisibility",
                            "sqs:GetQueueUrl",
                        ],
                        Resource: resArn,
                    },
                    {
                        Effect:   "Allow",
                        Action:   ["sqs:SendMessage"],
                        Resource: reqArn,
                    },
                    {
                        Effect:   "Allow",
                        Action:   ["sqs:PurgeQueue"],
                        Resource: [reqArn, resArn],
                    },
                ],
            }),
        ),
    }],
    tags: [{ key: "Project", value: "Alexa-ChatGPT" }],
});

// ── Lambda Function: ChatGPTFunction ─────────────────────────────────────────
// CloudFormation Logical ID: ChatGPTFunction  |  Type: AWS::Serverless::Function
//
// Uses @pulumi/aws (classic) so FileArchive handles zip+upload automatically,
// matching the sam build workflow. The binary must be pre-built before pulumi up:
//   GOOS=linux GOARCH=arm64 go build -o dist/alexa/bootstrap ./cmd/alexa
const ChatGPTFunction = new aws.lambda.Function("ChatGPTFunction", {
    functionName:  "chatGPT",
    runtime:       runtime as aws.lambda.Runtime,
    handler:       handler,
    architectures: [architecture],
    role:          ChatGPTFunctionRole.arn,
    timeout:       300,
    code:          new pulumi.asset.FileArchive("../dist/alexa"),
    layers:        [otelLayerArn],
    tracingConfig: { mode: "Active" },
    environment:   { variables: commonEnvVars },
    tags:          { Project: "Alexa-ChatGPT" },
});

// ── Lambda Permission: Alexa Skill → ChatGPTFunction ─────────────────────────
// SAM AlexaSkill event type auto-generates this AWS::Lambda::Permission.
const AlexaSkillEventPermission = new awsNative.lambda.Permission("AlexaSkillEventPermission", {
    action:       "lambda:InvokeFunction",
    functionName: ChatGPTFunction.functionName,
    principal:    "alexa-appkit.amazon.com",
});

// ── IAM Execution Role: ChatGPTRequests ───────────────────────────────────────
// SAM auto-generates this from the Policies block. Expanded SAM policy templates:
//   AWSXrayWriteOnlyAccess         → managed policy
//   AWSLambdaSQSQueueExecutionRole → managed policy (required for SQS event source mapping)
//   SQSSendMessagePolicy           → send on ResponsesQueue
//   S3CrudPolicy                   → CRUD on Bucket
const ChatGPTRequestsRole = new awsNative.iam.Role("ChatGPTRequestsRole", {
    roleName: `${cfnStackName}-ChatGPTRequestsRole`,
    assumeRolePolicyDocument: {
        Version: "2012-10-17",
        Statement: [{
            Effect:    "Allow",
            Principal: { Service: "lambda.amazonaws.com" },
            Action:    "sts:AssumeRole",
        }],
    },
    managedPolicyArns: [
        "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
        "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole",
        "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
    ],
    policies: [{
        policyName: "ChatGPTRequestsPolicy",
        policyDocument: pulumi.all([ResponsesQueue.arn, Bucket.arn]).apply(
            ([resArn, bucketArn]) => ({
                Version: "2012-10-17",
                Statement: [
                    {
                        Effect:   "Allow",
                        Action:   ["sqs:SendMessage"],
                        Resource: resArn,
                    },
                    {
                        Effect: "Allow",
                        Action: [
                            "s3:GetObject",
                            "s3:ListBucket",
                            "s3:GetBucketLocation",
                            "s3:GetObjectVersion",
                            "s3:PutObject",
                            "s3:PutObjectAcl",
                            "s3:GetLifecycleConfiguration",
                            "s3:DeleteObject",
                            "s3:GetBucketPolicy",
                        ],
                        Resource: [bucketArn, `${bucketArn}/*`],
                    },
                ],
            }),
        ),
    }],
    tags: [{ key: "Project", value: "Alexa-ChatGPT" }],
});

// ── Lambda Function: ChatGPTRequests ─────────────────────────────────────────
// CloudFormation Logical ID: ChatGPTRequests  |  Type: AWS::Serverless::Function
//
//   GOOS=linux GOARCH=arm64 go build -o dist/sqs/bootstrap ./cmd/sqs
const ChatGPTRequests = new aws.lambda.Function("ChatGPTRequests", {
    functionName:                 "chatGPTRequests",
    runtime:                      runtime as aws.lambda.Runtime,
    handler:                      handler,
    architectures:                [architecture],
    role:                         ChatGPTRequestsRole.arn,
    timeout:                      300,
    reservedConcurrentExecutions: 1,
    code:          new pulumi.asset.FileArchive("../dist/sqs"),
    layers:        [otelLayerArn],
    tracingConfig: { mode: "Active" },
    environment:   { variables: commonEnvVars },
    tags:          { Project: "Alexa-ChatGPT" },
});

// ── Lambda Event Source Mapping: RequestsQueue → ChatGPTRequests ──────────────
// SAM SQS event type auto-generates this AWS::Lambda::EventSourceMapping.
const MySQSEventSourceMapping = new awsNative.lambda.EventSourceMapping("MySQSEventSourceMapping", {
    functionName:   ChatGPTRequests.functionName,
    eventSourceArn: RequestsQueue.arn,
    batchSize:      1,
});

// ── Outputs (match CloudFormation Outputs section) ────────────────────────────
export const chatGPTLambdaArn   = ChatGPTFunction.arn;
export const chatGPTRequestsArn = ChatGPTRequests.arn;
export const requestsQueueArn   = RequestsQueue.arn;
export const requestsDLQArn     = RequestsDLQ.arn;
export const responsesQueueArn  = ResponsesQueue.arn;