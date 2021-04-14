/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cancel

import (
	"fmt"

	"github.com/spf13/cobra"
)

const description = `Cancel a staged command that has not been executed.

This can be done as the issuer of a command to revoke it if there was
an error while issuing it, or for a properly permissioned user to
explicitly deny a command to be run.

This does not remove the command from the audit history, but it does
mark it as deleted and prevent it from being scheduled for execution.

Commands that are already running or have already completed cannot be
canceled.
`

func NewCmdCancel() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a staged command that has not been executed",
		Long:  description,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("history called")
		},
	}
}
