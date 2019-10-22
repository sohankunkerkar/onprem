/*

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

package v1alpha1

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (j *JoinedCluster) IsCondition(condition JoinedClusterConditionType) *JoinedClusterConditions {
	for _, c := range j.Status.Conditions {
		if c.Type == condition && c.Status == ConditionTrue {
			return &c
		}
	}
	return nil
}

func (j *JoinedCluster) SetCondition(condition JoinedClusterConditionType) error {
	switch condition {
	case ConditionTypeAgentDisconnected:
		j.clearCondition(ConditionTypeAgentStale, "AgentDisconnected")
		j.clearCondition(ConditionTypeAgentConnected, "AgentDisconnected")
		j.setCondition(ConditionTypeAgentDisconnected, "DisconnectTimerExpired")
	case ConditionTypeAgentStale:
		j.clearCondition(ConditionTypeAgentDisconnected, "AgentStale")
		j.clearCondition(ConditionTypeAgentConnected, "AgentStale")
		j.setCondition(ConditionTypeAgentStale, "StaleTimerExpired")
	case ConditionTypeAgentConnected:
		j.clearCondition(ConditionTypeAgentStale, "AgentConnected")
		j.clearCondition(ConditionTypeAgentDisconnected, "AgentConnected")
		j.setCondition(ConditionTypeAgentConnected, "AgentConnected")
	case ConditionTypeReadyToJoin:
		j.setCondition(ConditionTypeReadyToJoin, "AgentReadyToJoin")
	}

	return nil
}

func (j *JoinedCluster) setCondition(condition JoinedClusterConditionType, reason string) {
	// update the condition if found in the slice
	for i, c := range j.Status.Conditions {
		if c.Type == condition {
			if c.Status != ConditionTrue {
				now := metav1.Now()
				c.LastTransitionTime = &now
			}
			c.Status = ConditionTrue
			c.Reason = &reason
			j.Status.Conditions[i] = c
			return
		}
	}
	now := metav1.Now()
	// create the condition since it doesn't exist
	j.Status.Conditions = append(j.Status.Conditions, JoinedClusterConditions{condition, ConditionTrue, &now, &reason, nil})
}

func (j *JoinedCluster) clearCondition(condition JoinedClusterConditionType, reason string) error {
	// update the condition status to false and the timestamp of transition to now
	for i, c := range j.Status.Conditions {
		if c.Type == condition {
			if c.Status != ConditionFalse {
				now := metav1.Now()
				c.LastTransitionTime = &now
			}
			c.Status = ConditionFalse
			c.Reason = &reason
			j.Status.Conditions[i] = c
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Condition %v was not set", condition))
}
