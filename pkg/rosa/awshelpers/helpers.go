/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awshelpers

import (
	"fmt"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/aws"
)

// GetOCMRoleName generates the OCM role name from prefix, role type, and postfix (external ID).
func GetOCMRoleName(prefix, role, postfix string) string {
	return fmt.Sprintf("%s-%s-Role-%s", prefix, role, postfix)
}

// GetPolicyName returns the standard policy name for a role.
func GetPolicyName(name string) string {
	return fmt.Sprintf("%s-Policy", name)
}

// GetAdminPolicyName returns the admin policy name for a role.
func GetAdminPolicyName(name string) string {
	return fmt.Sprintf("%s-Admin-Policy", name)
}

// GetNoConsolePolicyName returns the no-console policy name for a role.
func GetNoConsolePolicyName(name string) string {
	return fmt.Sprintf("%s-NoConsole-Policy", name)
}

// GetPolicyArn returns the ARN for a policy with a given name.
func GetPolicyArn(partition, accountID, name, path string) string {
	return getPolicyARN(partition, accountID, name, path)
}

// GetPolicyArnWithSuffix returns the ARN for a customer-managed policy.
func GetPolicyArnWithSuffix(partition, accountID, name, path string) string {
	return getPolicyARN(partition, accountID, GetPolicyName(name), path)
}

// GetAdminPolicyARN returns the ARN for the admin policy.
func GetAdminPolicyARN(partition, accountID, name, path string) string {
	return getPolicyARN(partition, accountID, GetAdminPolicyName(name), path)
}

// GetNoConsolePolicyARN returns the ARN for the no-console policy.
func GetNoConsolePolicyARN(partition, accountID, name, path string) string {
	return getPolicyARN(partition, accountID, GetNoConsolePolicyName(name), path)
}

// getPolicyARN constructs an IAM policy ARN.
func getPolicyARN(partition, accountID, name, path string) string {
	str := fmt.Sprintf("arn:%s:iam::%s:policy", partition, accountID)
	if path != "" {
		str = fmt.Sprintf("%s%s", str, path)
		return fmt.Sprintf("%s%s", str, name)
	}
	return fmt.Sprintf("%s/%s", str, name)
}

// GetPolicyDetails retrieves from the map the policy details for customer-managed policies.
func GetPolicyDetails(policies map[string]*cmv1.AWSSTSPolicy, key string) string {
	policy, ok := policies[key]
	if ok {
		return policy.Details()
	}

	return ""
}

// InterpolatePolicyDocument replaces template variables in a policy document.
func InterpolatePolicyDocument(partition, doc string, replacements map[string]string) string {
	for key, val := range replacements {
		doc = strings.ReplaceAll(doc, fmt.Sprintf("%%{%s}", key), val)
	}

	// TODO Remove once MCC policies are all updated
	doc = strings.ReplaceAll(doc, "arn:aws:", fmt.Sprintf("arn:%s:", partition))

	return doc
}

// GetJumpAccount returns the OCM jump account ID for the given environment.
func GetJumpAccount(env string) string {
	return aws.JumpAccounts[env]
}
