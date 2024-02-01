// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bedrock_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfbedrock "github.com/hashicorp/terraform-provider-aws/internal/service/bedrock"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccBedrockCustomModel_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_bedrock_custom_model.test"
	var v bedrock.GetModelCustomizationJobOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.BedrockEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BedrockEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "base_model_identifier"),
					resource.TestCheckNoResourceAttr(resourceName, "custom_model_arn"),
					resource.TestCheckNoResourceAttr(resourceName, "custom_model_kms_key_id"),
					resource.TestCheckResourceAttr(resourceName, "custom_model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "customization_type", "FINE_TUNING"),
					resource.TestCheckResourceAttr(resourceName, "hyperparameters.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "hyperparameters.batchSize", "1"),
					resource.TestCheckResourceAttr(resourceName, "hyperparameters.epochCount", "1"),
					resource.TestCheckResourceAttr(resourceName, "hyperparameters.learningRate", "0.005"),
					resource.TestCheckResourceAttr(resourceName, "hyperparameters.learningRateWarmupSteps", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "job_arn"),
					resource.TestCheckResourceAttr(resourceName, "job_name", rName),
					resource.TestCheckResourceAttr(resourceName, "job_status", "InProgress"),
					resource.TestCheckResourceAttr(resourceName, "output_data_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "output_data_config.0.s3_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "role_arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "training_data_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "training_data_config.0.s3_uri"),
					resource.TestCheckNoResourceAttr(resourceName, "training_metrics"),
					resource.TestCheckResourceAttr(resourceName, "validation_data_config.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "validation_metrics"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base_model_identifier"},
			},
		},
	})
}

func TestAccBedrockCustomModel_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_bedrock_custom_model.test"
	var v bedrock.GetModelCustomizationJobOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.BedrockEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BedrockEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfbedrock.ResourceCustomModel, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBedrockCustomModel_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_bedrock_custom_model.test"
	var v bedrock.GetModelCustomizationJobOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.BedrockEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BedrockEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base_model_identifier"},
			},
			{
				Config: testAccCustomModelConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccCustomModelConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccBedrockCustomModel_validationDataConfig(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_bedrock_custom_model.test"
	var v bedrock.GetModelCustomizationJobOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.BedrockEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BedrockEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_validationDataConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "validation_data_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "validation_data_config.0.validator.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "validation_data_config.0.validator.0.s3_uri"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base_model_identifier"},
			},
		},
	})
}

func TestAccBedrockCustomModel_validationDataConfigWaitForCompletion(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_bedrock_custom_model.test"
	var v bedrock.GetModelCustomizationJobOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.BedrockEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BedrockEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_validationDataConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "validation_data_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "validation_data_config.0.validator.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "validation_data_config.0.validator.0.s3_uri"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base_model_identifier"},
			},
			{
				// TODO Wait for completion.
				PreConfig:    func() {},
				Config:       testAccCustomModelConfig_validationDataConfig(rName),
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "job_status", "Completed"),
				),
			},
		},
	})
}

func TestAccBedrockCustomModel_vpcConfig(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_bedrock_custom_model.test"
	var v bedrock.GetModelCustomizationJobOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.BedrockEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.BedrockEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCustomModelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomModelConfig_vpcConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomModelExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.0.subnet_ids.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base_model_identifier"},
			},
		},
	})
}

func testAccCheckCustomModelDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).BedrockClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_bedrock_custom_model" {
				continue
			}

			_, err := tfbedrock.FindModelCustomizationJobByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Bedrock Custom Model %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckCustomModelExists(ctx context.Context, n string, v *bedrock.GetModelCustomizationJobOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).BedrockClient(ctx)

		output, err := tfbedrock.FindModelCustomizationJobByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCustomModelConfig_base(rName string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_partition" "current" {}

resource "aws_s3_bucket" "training" {
  bucket = "%[1]s-training"
}

resource "aws_s3_bucket" "validation" {
  bucket = "%[1]s-validation"
}

resource "aws_s3_bucket" "output" {
  bucket        = "%[1]s-output"
  force_destroy = true
}

resource "aws_s3_object" "training" {
  bucket = aws_s3_bucket.training.id
  key    = "data/train.jsonl"
  source = "test-fixtures/train.jsonl"
}

resource "aws_s3_object" "validation" {
  bucket = aws_s3_bucket.validation.id
  key    = "data/validate.jsonl"
  source = "test-fixtures/validate.jsonl"
}

resource "aws_iam_role" "test" {
  name = %[1]q

  # See https://docs.aws.amazon.com/bedrock/latest/userguide/model-customization-iam-role.html#model-customization-iam-role-trust.
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Service": "bedrock.amazonaws.com"
    },
    "Action": "sts:AssumeRole",
    "Condition": {
      "StringEquals": {
        "aws:SourceAccount": "${data.aws_caller_identity.current.account_id}"
      },
      "ArnEquals": {
        "aws:SourceArn": "arn:${data.aws_partition.current.partition}:bedrock:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:model-customization-job/*"
      }
    }
  }]
}
EOF
}

# See https://docs.aws.amazon.com/bedrock/latest/userguide/model-customization-iam-role.html#model-customization-iam-role-s3.
resource "aws_iam_policy" "training" {
  name = "%[1]s-training"

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [{
      "Effect" : "Allow",
      "Action" : [
        "s3:GetObject",
        "s3:ListBucket"
      ],
      "Resource" : [
        "${aws_s3_bucket.training.arn}",
        "${aws_s3_bucket.training.arn}/*",
        "${aws_s3_bucket.validation.arn}",
        "${aws_s3_bucket.validation.arn}/*"
      ]
    }]
  })
}

resource "aws_iam_policy" "output" {
  name = "%[1]s-output"

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [{
      "Effect" : "Allow",
      "Action" : [
        "s3:GetObject",
        "s3:PutObject",
        "s3:ListBucket"
      ],
      "Resource" : [
        "${aws_s3_bucket.output.arn}",
        "${aws_s3_bucket.output.arn}/*"
      ]
    }]
  })
}

resource "aws_iam_role_policy_attachment" "training" {
  role       = aws_iam_role.test.name
  policy_arn = aws_iam_policy.training.arn
}

resource "aws_iam_role_policy_attachment" "output" {
  role       = aws_iam_role.test.name
  policy_arn = aws_iam_policy.output.arn
}

data "aws_bedrock_foundation_model" "test" {
  model_id = "amazon.titan-text-express-v1"
}
`, rName)
}

func testAccCustomModelConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccCustomModelConfig_base(rName), fmt.Sprintf(`
resource "aws_bedrock_custom_model" "test" {
  custom_model_name     = %[1]q
  job_name              = %[1]q
  base_model_identifier = data.aws_bedrock_foundation_model.test.model_arn
  role_arn              = aws_iam_role.test.arn

  hyperparameters = {
    "epochCount"              = "1"
    "batchSize"               = "1"
    "learningRate"            = "0.005"
    "learningRateWarmupSteps" = "0"
  }

  output_data_config {
    s3_uri = "s3://${aws_s3_bucket.output.id}/data/"
  }

  training_data_config {
    s3_uri = "s3://${aws_s3_bucket.training.id}/data/train.jsonl"
  }
}
`, rName))
}

func testAccCustomModelConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccCustomModelConfig_base(rName), fmt.Sprintf(`
resource "aws_bedrock_custom_model" "test" {
  custom_model_name     = %[1]q
  job_name              = %[1]q
  base_model_identifier = data.aws_bedrock_foundation_model.test.model_arn
  role_arn              = aws_iam_role.test.arn

  hyperparameters = {
    "epochCount"              = "1"
    "batchSize"               = "1"
    "learningRate"            = "0.005"
    "learningRateWarmupSteps" = "0"
  }

  output_data_config {
    s3_uri = "s3://${aws_s3_bucket.output.id}/data/"
  }

  training_data_config {
    s3_uri = "s3://${aws_s3_bucket.training.id}/data/train.jsonl"
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccCustomModelConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccCustomModelConfig_base(rName), fmt.Sprintf(`
resource "aws_bedrock_custom_model" "test" {
  custom_model_name     = %[1]q
  job_name              = %[1]q
  base_model_identifier = data.aws_bedrock_foundation_model.test.model_arn
  role_arn              = aws_iam_role.test.arn

  hyperparameters = {
    "epochCount"              = "1"
    "batchSize"               = "1"
    "learningRate"            = "0.005"
    "learningRateWarmupSteps" = "0"
  }

  output_data_config {
    s3_uri = "s3://${aws_s3_bucket.output.id}/data/"
  }

  training_data_config {
    s3_uri = "s3://${aws_s3_bucket.training.id}/data/train.jsonl"
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccCustomModelConfig_validationDataConfig(rName string) string {
	return acctest.ConfigCompose(testAccCustomModelConfig_base(rName), fmt.Sprintf(`
resource "aws_bedrock_custom_model" "test" {
  custom_model_name     = %[1]q
  job_name              = %[1]q
  base_model_identifier = data.aws_bedrock_foundation_model.test.model_arn
  role_arn              = aws_iam_role.test.arn

  hyperparameters = {
    "epochCount"              = "1"
    "batchSize"               = "1"
    "learningRate"            = "0.005"
    "learningRateWarmupSteps" = "0"
  }

  output_data_config {
    s3_uri = "s3://${aws_s3_bucket.output.id}/data/"
  }

  training_data_config {
    s3_uri = "s3://${aws_s3_bucket.training.id}/data/train.jsonl"
  }

  validation_data_config {
    validator {
      s3_uri = "s3://${aws_s3_bucket.validation.id}/data/validate.jsonl"
    }
  }
}
`, rName))
}

func testAccCustomModelConfig_vpcConfig(rName string) string {
	return acctest.ConfigCompose(testAccCustomModelConfig_base(rName), acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_iam_policy" "vpc" {
  name = "%[1]s-vpc"

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [{
      "Effect" : "Allow",
      "Action" : [
        "ec2:DescribeNetworkInterfaces",
        "ec2:DescribeVpcs",
        "ec2:DescribeDhcpOptions",
        "ec2:DescribeSubnets",
        "ec2:DescribeSecurityGroups",
        "ec2:CreateNetworkInterface",
        "ec2:CreateNetworkInterfacePermission",
        "ec2:DeleteNetworkInterface",
        "ec2:DeleteNetworkInterfacePermission",
        "ec2:CreateTags"
      ],
      "Resource" : "*"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "vpc" {
  role       = aws_iam_role.test.name
  policy_arn = aws_iam_policy.vpc.arn
}

resource "aws_bedrock_custom_model" "test" {
  custom_model_name     = %[1]q
  job_name              = %[1]q
  base_model_identifier = data.aws_bedrock_foundation_model.test.model_arn
  role_arn              = aws_iam_role.test.arn

  hyperparameters = {
    "epochCount"              = "1"
    "batchSize"               = "1"
    "learningRate"            = "0.005"
    "learningRateWarmupSteps" = "0"
  }

  output_data_config {
    s3_uri = "s3://${aws_s3_bucket.output.id}/data/"
  }

  training_data_config {
    s3_uri = "s3://${aws_s3_bucket.training.id}/data/train.jsonl"
  }

  vpc_config {
    security_group_ids = [aws_security_group.test.id]
    subnet_ids         = aws_subnet.test[*].id
  }
}
`, rName))
}
