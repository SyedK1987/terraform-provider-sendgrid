# terraform-provider-sendgrid

Sendgrid terraform provider.

## Authentication

The Sendgrid provider offers a flexible means of providing credentials for authentication. The following methods are supported, in this order, and explained below:

 - Environment variables
 - Configuration file

### Environment variables

You can provide your Sendgrid API key via the `SENDGRID_API_KEY` environment variable, representing your Sendgrid API key.

```terraform
provider "sendgrid" {}
```

Usage:

```shell
$ export SENDGRID_API_KEY="SG.*************************************"
$ terraform plan
```

### Configuration file

You can provide your Sendgrid API key via the `apikey` attribute in the `provider` block in your configuration, representing your Sendgrid API key.

```terraform
provider "sendgrid" {
  apikey = "SG.*************************************"
}
```

## Buil Provider

Run the following command to build the provider

```shell
$ make
```

## Install Provider

Run the following command to install the provider

```shell
$ make install
```

## Test Provider

Navigate to the `examples` directory.

```shell
$ cd examples/resources/sendgrid_api_key
```

Run the following command to initialize the provider

```shell
$ terraform init && terraform apply
```

## Debugging

Run the following command to enable debugging

```shell
$ export TF_LOG="TRACE"
```
