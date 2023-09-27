<div align="center">
<h1 style="text-decoration: underline">Swagger2YAML Github Action</h1>
<p>
This GitHub Action converts Swagger JSON to YAML for use with AWS API Gateway configurations.The conversion is mainly based around swagger files generated from <a href="https://github.com/grpc-ecosystem/grpc-gateway">grpc-Gateway</a>.
</p>
<a href="https://github.com/ImtiyaazL/swagger2yaml-github-action/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ImtiyaazL/swagger2yaml-github-action?color=379c9c"/></a>
<a href="https://github.com/ImtiyaazL/swagger2yaml-github-action/releases"><img src="https://img.shields.io/github/v/release/ImtiyaazL/swagger2yaml-github-action?color=379c9c&logoColor=ffffff"/></a>
</div>

## Usage

To use this action in your GitHub workflows, you can include it as a step in your workflow YAML file. Here's an example of how to use it:

```yaml
name: Convert Swagger to YAML

on:
  push:
    branches:
      - main

jobs:
  convert:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Convert Swagger to YAML
        uses: ImtiyaazL/-github-action@v1.0.0
        with:
          account: ${{ secrets.AWS_ACCOUNT }}
          host: ${{ secrets.CLOUD_HOST }}
          input: path/to/swagger.json
          output: path/to/swagger.yaml
          region: ${{ secrets.AWS_REGION }}
          vpc: ${{ secrets.AWS_VPC }}

```

### Inputs

- `account` (required): AWS account ID.
- `host` (required): Cloud host URL.
- `input` (required): Path to the input Swagger JSON file.
- `output` (required): Path to the output YAML file.
- `region` (required): AWS region.
- `vpc` (required): AWS VPC ID.

## Feedback and Contributions

We value your feedback and contributions. If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the [GitHub repository](https://github.com/ImtiyaazL/swagger2yaml-github-action).

Thank you for using the Swagger2YAML GitHub Action. We look forward to enhancing your API development process!

## License

This project is licensed under the [MIT License](LICENSE).
