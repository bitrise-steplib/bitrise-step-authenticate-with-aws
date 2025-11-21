# Authenticate with Amazon Web Services (AWS)

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-authenticate-with-aws?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-authenticate-with-aws/releases)

The step authenticates with AWS using an OIDC token.

<details>
<summary>Description</summary>

This step authenticates with Amazon Web Services (AWS) using an OpenID Connect (OIDC) token or an account key.

For OIDC based authentication it retrieves an identity token from Bitrise, assumes the specified AWS role using the token, and generates temporary AWS credentials.

The access key details can be created on the AWS Management Console under IAM roles.

The generated AWS credentials are then set as environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN`) for use in subsequent steps.
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `region` | The AWS region to use. |  | `us-east-1` |
| `access_key_id` | The AWS Access Key ID.  You can get this from the AWS Console under IAM users. | sensitive |  |
| `secret_access_key` | The AWS Secret Access Key.  You can get this from the AWS Console under IAM users. | sensitive |  |
| `audience` | The audience for the identity token.  This could be the URL of the service you want to access with the token or a specific identifier provided by the service. |  |  |
| `role_arn` | The ARN of the AWS role to assume.  You can find the ARN in the AWS Management Console under IAM roles. |  |  |
| `session_name` | The session name for the assumed role.  If not provided, a default name will be generated with the format `bitrise-<build-number>`. |  | `bitrise-$BITRISE_BUILD_NUMBER` |
| `docker_login` | Performs Docker login with an auth token.  It is supported only on the Linux stacks. | required | `false` |
| `build_url` | Unique build URL of this build on Bitrise.io.  By default the step will use the Bitrise API. | required | `$BITRISE_BUILD_URL` |
| `build_api_token` | The build's API Token for the build on Bitrise.io  This will be used to communicate with the Bitrise API | required, sensitive | `$BITRISE_BUILD_API_TOKEN` |
| `verbose` | Enable logging additional information for debugging. | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `AWS_ACCESS_KEY_ID` | The newly generated AWS access key ID. |
| `AWS_SECRET_ACCESS_KEY` | The newly generated AWS secret access key. |
| `AWS_SESSION_TOKEN` | The newly generated AWS session token. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-authenticate-with-aws/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-authenticate-with-aws/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
