package redistore

import (
	"fmt"
)

const namespace = "okpock:"

func key(k string, args ...interface{}) string {
	return fmt.Sprintf(namespace+k, args...)
}

const (
	kzUsers = "users"

	khUser                  = "user:%s"
	khUserUsername          = khUser + ":index:username"
	khUserEmail             = khUser + ":index:email"
	khUserConfirmationToken = khUser + ":index:confirmation_token"
	khUserRecoveryToken     = khUser + ":index:recovery_token"
	khUserEmailChangeToken  = khUser + ":index:change_token"
)
