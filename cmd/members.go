/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type RoleResponse struct {
	Role    Role     `json:"role"`
	Members []Member `json:"members"`
}

type Role struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Member struct {
	URI       string `json:"uri"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	ID        int    `json:"id"`
}

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "List project members",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		role, _ := cmd.Flags().GetString("role")

		// Get persistent flags from root command
		baseURL, _ := rootCmd.PersistentFlags().GetString("base")
		username, _ := rootCmd.PersistentFlags().GetString("user")
		password, _ := rootCmd.PersistentFlags().GetString("pwd")

		// Build URL
		if role == "" {
			listMembers(baseURL, project, username, password)
		} else {
			listRoleMembers(baseURL, project, username, password, role)
		}
	},
}

func listMembers(baseURL, project, username, password string) {
	url := fmt.Sprintf("%s/rest/project/%s/roles/members",
		baseURL,
		url.PathEscape(project))

	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Add basic auth
	req.SetBasicAuth(username, password)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}

		var roleResponses []RoleResponse
		err = json.Unmarshal(body, &roleResponses)
		if err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			return
		}

		uniqueMembers := make(map[string]Member)
		for _, roleResp := range roleResponses {
			for _, member := range roleResp.Members {
				if _, exists := uniqueMembers[member.Name]; !exists {
					uniqueMembers[member.Name] = member
				}
			}
		}

		fmt.Print("Name\t\tEmail\n")
		for _, member := range uniqueMembers {
			fmt.Printf("%s\t\t%s\n", member.Name, member.Email)
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed. Status code: %d, error reading response: %v\n", resp.StatusCode, err)
			return
		}
		fmt.Printf("Failed. Status code: %d, message: %s\n", resp.StatusCode, string(body))
	}
}

func listRoleMembers(baseURL, project, username, password, role string) {
	url := fmt.Sprintf("%s/rest/project/%s/role/%s/members",
		baseURL,
		url.PathEscape(project),
		url.PathEscape(role))

	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Add basic auth
	req.SetBasicAuth(username, password)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}

		var members []Member
		err = json.Unmarshal(body, &members)
		if err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			return
		}

		fmt.Print("Name\t\tEmail\n")
		for _, member := range members {
			fmt.Printf("%s\t\t%s\n", member.Name, member.Email)
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed. Status code: %d, error reading response: %v\n", resp.StatusCode, err)
			return
		}
		fmt.Printf("Failed. Status code: %d, message: %s\n", resp.StatusCode, string(body))
	}
}

func init() {
	projectCmd.AddCommand(membersCmd)

	membersCmd.Flags().String("project", "", "Specify the project name (required)")
	membersCmd.Flags().String("role", "", "role name (optional)")

	grantCmd.MarkFlagRequired("project")
}
