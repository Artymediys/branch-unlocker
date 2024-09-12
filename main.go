package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"branch-unlocker/env"

	"github.com/xanzy/go-gitlab"
)

func removeBranchProtection(glc *gitlab.Client, projectID int, branchName string) error {
	_, err := glc.ProtectedBranches.UnprotectRepositoryBranches(projectID, branchName)
	if err != nil {
		return err
	}

	return nil
}

func getAllProjects(glc *gitlab.Client, groupID int) ([]*gitlab.Project, error) {
	includeSubGroups := true
	listProjectsOptions := &gitlab.ListGroupProjectsOptions{
		IncludeSubGroups: &includeSubGroups,
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	var allProjects []*gitlab.Project
	for {
		projects, resp, err := glc.Groups.ListGroupProjects(groupID, listProjectsOptions)
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)
		if resp.NextPage == 0 {
			break
		}

		listProjectsOptions.Page = resp.NextPage
	}

	return allProjects, nil
}

func getProjectPath(project *gitlab.Project) string {
	urlParts := strings.Split(project.WebURL, "/")
	return urlParts[len(urlParts)-1]
}

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	log.Println("Loading configuration...")
	config, err := env.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("Error creating configuration instance: %v\n", err)
	}

	log.Println("Creating GitLab client...")
	glc, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.URL))
	if err != nil {
		log.Fatalf("Error creating gitlab client: %v\n", err)
	}

	log.Println("Getting project list...")
	projects, err := getAllProjects(glc, config.GroupID)
	if err != nil {
		log.Fatalf("Error getting projects: %v\n", err)
	}

	log.Println("Unprotecting projects branches...")
	for _, project := range projects {
		projectSysName := getProjectPath(project)
		log.Println("Working with project: " + projectSysName)
		for _, branch := range config.Branches {
			if err = removeBranchProtection(glc, project.ID, branch); err != nil {
				log.Printf("Error to unprotect branch \"%s\" in project %d: %v", branch, project.ID, err)
			} else {
				log.Printf("Branch \"%s\" in project \"%s\" has been unprotected!\n", branch, projectSysName)
			}
			time.Sleep(2 * time.Second)
		}

		log.Println("Finish with project: " + projectSysName)

		log.Println("Waiting next project checking...")
		time.Sleep(10 * time.Second)
	}
}
