package canu

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/belitre/canu/pkg/canu/profile"
	"github.com/belitre/canu/pkg/config"
	"github.com/nexidian/gocliselect"
	"gopkg.in/ini.v1"
)

const (
	ssoStartURL = "sso_start_url"
)

type canu struct {
	cfg *config.Config
}

type Canu interface {
	Run() error
}

func New(cfg *config.Config) Canu {
	return &canu{
		cfg: cfg,
	}
}

func (c *canu) Run() error {
	if len(c.cfg.ConfigPath) == 0 {
		return fmt.Errorf("config-path for AWS config file is empty")
	}

	listProfiles, err := getLocalAwsProfiles(c.cfg.ConfigPath, c.cfg.Include, c.cfg.Exclude)

	if err != nil {
		return err
	}

	if len(listProfiles) == 0 {
		return fmt.Errorf("no profiles found on file %s", c.cfg.ConfigPath)
	}

	slices.Sort(listProfiles)

	menu := gocliselect.NewMenu("Choose a profile")

	for _, p := range listProfiles {
		menu.AddItem(p, p)
	}

	choice := menu.Display()

	if len(choice) == 0 {
		fmt.Println("\n\nno profile selected, bye!")
		return nil
	}

	p := profile.New(choice, c.cfg.ConfigPath)

	if err := p.CheckAWSProfileStatus(); err != nil {
		return err
	}

	// ok, we have permissions, so let's save the profile name so we can set it later with the script and the alias
	profileSaveFileName := path.Join(c.cfg.UserHomePath, config.CanuSaveProfileFileName)

	if err := os.WriteFile(profileSaveFileName, []byte(choice), 0600); err != nil {
		return fmt.Errorf("error while saving canu config with profile name %s in path %s: %v", choice, profileSaveFileName, err)
	}

	return err
}

func getLocalAwsProfiles(configPath string, includes []string, excludes []string) ([]string, error) {
	listProfiles := []string{}

	f, err := ini.Load(configPath)

	if err != nil {
		return listProfiles, fmt.Errorf("error while reading AWS config file from path %s: %v", configPath, err)
	}

	for _, v := range f.Sections() {
		if len(v.Keys()) != 0 {
			parts := strings.Split(v.Name(), " ")

			if len(parts) == 2 && parts[0] == "profile" { // skip default
				// for now only support profiles with sso_start_url!
				if !v.HasKey(ssoStartURL) {
					continue
				}

				if len(excludes) > 0 {

					isExcluded := false

					for _, exclude := range excludes {
						if strings.Contains(strings.ToUpper(parts[1]), strings.ToUpper(exclude)) {
							isExcluded = true
							break
						}
					}

					if isExcluded {
						continue
					}
				}

				if len(includes) > 0 {
					for _, include := range includes {
						if strings.Contains(strings.ToUpper(parts[1]), strings.ToUpper(include)) {
							listProfiles = append(listProfiles, parts[1])
						}
					}

					continue
				}

				listProfiles = append(listProfiles, parts[1])
			}

		}
	}

	return listProfiles, nil
}
