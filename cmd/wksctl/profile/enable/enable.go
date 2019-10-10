package enable

import (
	"fmt"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weaveworks/wksctl/cmd/wksctl/profile/constants"
	"github.com/weaveworks/wksctl/pkg/git"
)

var Cmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable profile",
	Long: `To enable the profile from a specific git URL, run

wksctl profile enable --git-url=<profile_repository> [--revision=master] [--profile-dir=profiles] [--push=true]

If you'd like to specify the revision other than the master branch, use --revision flag.
To disable auto-push, pass --push=false.
`,
	Args: profileEnableArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		return profileEnableRun(profileEnableParams)
	},
	SilenceUsage: true,
}

type profileEnableFlags struct {
	gitUrl     []string
	revision   []string
	push       bool
	profileDir string
}

var profileEnableParams profileEnableFlags

func init() {
	Cmd.Flags().StringVar(&profileEnableParams.profileDir, "profile-dir", "profiles", "specify a directory for storing profiles")
	Cmd.Flags().StringSliceVar(&profileEnableParams.gitUrl, "git-url", []string{""}, "enable profile from the gitUrl")
	Cmd.Flags().StringSliceVar(&profileEnableParams.revision, "revision", []string{"master"}, "use this revision of the profile")
	Cmd.Flags().BoolVar(&profileEnableParams.push, "push", true, "auto push after enable the profile")
}

func profileEnableArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("profile enable does not require any argument")
	}
	return nil
}

func enableSingleProfile(profileDir, repoUrl, revision string) error {
	if repoUrl == constants.AppDevAlias {
		repoUrl = constants.AppDevRepoURL
	}

	if err := git.IsGitURL(repoUrl); err != nil {
		return err
	}

	hostName, repoName, err := git.HostAndRepoPath(repoUrl)
	if err != nil {
		return err
	}
	clonePath := path.Join(profileDir, hostName, repoName)

	log.Info("Adding the profile to the local repository...")
	err = git.SubtreeAdd(clonePath, repoUrl, revision)
	if err != nil {
		return err
	}
	log.Info("Added the profile to the local repository.")

	return nil
}

func profileEnableRun(params profileEnableFlags) error {
	if len(params.gitUrl) != len(params.revision) {
		return errors.New("length of Git URLs and Revisions don't match")
	}

	errs := []error{}
	for i := 0; i < len(params.gitUrl); i++ {
		err := enableSingleProfile(params.profileDir, params.gitUrl[i], params.revision[i])
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("profile enabling errors: %v", errs)
	}

	// The default behaviour is auto-commit and push
	if params.push {
		if err := git.Push(); err != nil {
			return err
		}
	}

	return nil
}
