package v1alpha1

import (
	"fmt"
	"strings"

	"gopkg.in/robfig/cron.v2"
)

func (r Restic) IsValid() error {
	for i, fg := range r.Spec.FileGroups {
		if fg.RetentionPolicyName == "" {
			continue
		}

		found := false
		for _, policy := range r.Spec.RetentionPolicies {
			if policy.Name == fg.RetentionPolicyName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("spec.fileGroups[%d].retentionPolicyName %s is not found", i, fg.RetentionPolicyName)
		}
	}

	_, err := cron.Parse(r.Spec.Schedule)
	if err != nil {
		return fmt.Errorf("spec.schedule %s is invalid. Reason: %s", r.Spec.Schedule, err)
	}
	if r.Spec.Backend.StorageSecretName == "" {
		return fmt.Errorf("missing repository secret name")
	}
	return nil
}

func (r Recovery) IsValid() error {
	if len(r.Spec.Paths) == 0 {
		return fmt.Errorf("missing filegroup paths")
	}
	if len(r.Spec.RecoveredVolumes) == 0 {
		return fmt.Errorf("missing recovery volume")
	}

	if r.Spec.Repository == "" {
		return fmt.Errorf("missing repository name")
	} else {
		if !(strings.HasPrefix(r.Spec.Repository, "deployment.") ||
			strings.HasPrefix(r.Spec.Repository, "replicationcontroller.") ||
			strings.HasPrefix(r.Spec.Repository, "replicaset.") ||
			strings.HasPrefix(r.Spec.Repository, "statefulset.") ||
			strings.HasPrefix(r.Spec.Repository, "daemonset.")) {
			return fmt.Errorf("invalid repository name")
		}
	}
	if r.Spec.Snapshot != "" {
		if !strings.HasPrefix(r.Spec.Snapshot, r.Spec.Repository+"-") {
			return fmt.Errorf("invalid snapshot name")
		}
	}
	return nil
}
