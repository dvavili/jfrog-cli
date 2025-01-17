package scan

import (
	"os"
	"os/exec"
	"strings"

	"github.com/jfrog/jfrog-cli/utils/progressbar"

	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/utils"
	commandsutils "github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/utils"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/common/spec"
	corecommondocs "github.com/jfrog/jfrog-cli-core/v2/docs/common"
	coreconfig "github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit"
	_go "github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/go"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/java"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/npm"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/nuget"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/python"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/scan"
	"github.com/jfrog/jfrog-cli/docs/common"
	auditdocs "github.com/jfrog/jfrog-cli/docs/scan/audit"
	auditgodocs "github.com/jfrog/jfrog-cli/docs/scan/auditgo"
	auditgradledocs "github.com/jfrog/jfrog-cli/docs/scan/auditgradle"
	"github.com/jfrog/jfrog-cli/docs/scan/auditmvn"
	auditnpmdocs "github.com/jfrog/jfrog-cli/docs/scan/auditnpm"
	auditpipdocs "github.com/jfrog/jfrog-cli/docs/scan/auditpip"
	auditpipenvdocs "github.com/jfrog/jfrog-cli/docs/scan/auditpipenv"
	buildscandocs "github.com/jfrog/jfrog-cli/docs/scan/buildscan"
	scandocs "github.com/jfrog/jfrog-cli/docs/scan/scan"
	"github.com/jfrog/jfrog-cli/utils/cliutils"
	"github.com/urfave/cli"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const auditScanCategory = "Audit & Scan"

func GetCommands() []cli.Command {
	return cliutils.GetSortedCommands(cli.CommandsByName{
		{
			Name:         "audit",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.Audit),
			Aliases:      []string{"aud"},
			Usage:        auditdocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit", auditdocs.GetDescription(), auditdocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditCmd,
		},
		{
			Name:         "audit-mvn",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.AuditMvn),
			Aliases:      []string{"am"},
			Usage:        auditmvn.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit-mvn", auditmvn.GetDescription(), auditmvn.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditMvnCmd,
		},
		{
			Name:         "audit-gradle",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.AuditGradle),
			Aliases:      []string{"ag"},
			Usage:        auditgradledocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit-gradle", auditgradledocs.GetDescription(), auditgradledocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditGradleCmd,
		},
		{
			Name:         "audit-npm",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.AuditNpm),
			Aliases:      []string{"an"},
			Usage:        auditnpmdocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit-npm", auditnpmdocs.GetDescription(), auditnpmdocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditNpmCmd,
		},
		{
			Name:         "audit-go",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.AuditGo),
			Aliases:      []string{"ago"},
			Usage:        auditgodocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit-go", auditgodocs.GetDescription(), auditgodocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditGoCmd,
		},
		{
			Name:         "audit-pip",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.AuditPip),
			Aliases:      []string{"ap"},
			Usage:        auditpipdocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit-pip", auditpipdocs.GetDescription(), auditpipdocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditPipCmd,
		},
		{
			Name:         "audit-pipenv",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.AuditPip),
			Aliases:      []string{"ape"},
			Usage:        auditpipenvdocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("audit-pipenv", auditpipenvdocs.GetDescription(), auditpipenvdocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       AuditPipenvCmd,
		},
		{
			Name:         "scan",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.XrScan),
			Aliases:      []string{"s"},
			Usage:        scandocs.GetDescription(),
			HelpName:     corecommondocs.CreateUsage("scan", scandocs.GetDescription(), scandocs.Usage),
			UsageText:    scandocs.GetArguments(),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       ScanCmd,
		},
		{
			Name:         "build-scan",
			Category:     auditScanCategory,
			Flags:        cliutils.GetCommandFlags(cliutils.BuildScan),
			Aliases:      []string{"bs"},
			Usage:        buildscandocs.GetDescription(),
			UsageText:    buildscandocs.GetArguments(),
			HelpName:     corecommondocs.CreateUsage("build-scan", buildscandocs.GetDescription(), buildscandocs.Usage),
			ArgsUsage:    common.CreateEnvVars(),
			BashComplete: corecommondocs.CreateBashCompletionFunc(),
			Action:       BuildScan,
		},
	})
}

func AuditCmd(c *cli.Context) error {
	wd, err := os.Getwd()
	if errorutils.CheckError(err) != nil {
		return err
	}
	detectedTechnologies, err := coreutils.DetectTechnologies(wd, false, false)
	if err != nil {
		return err
	}
	detectedTechnologiesString := coreutils.DetectedTechnologiesToString(detectedTechnologies)
	if detectedTechnologiesString == "" {
		log.Info("Could not determine the package manager / build tool used by this project.")
		return nil
	}
	log.Info("Detected: " + detectedTechnologiesString)
	var failBuildErr error
	for tech := range detectedTechnologies {
		switch tech {
		case coreutils.Maven:
			err = AuditMvnCmd(c)
		case coreutils.Gradle:
			err = AuditGradleCmd(c)
		case coreutils.Npm:
			err = AuditNpmCmd(c)
		case coreutils.Go:
			err = AuditGoCmd(c)
		case coreutils.Pip:
			err = AuditPipCmd(c)
		case coreutils.Pipenv:
			err = AuditPipenvCmd(c)
		case coreutils.Dotnet:
			break
		case coreutils.Nuget:
			err = AuditNugetCmd(c)
		default:
			log.Info(string(tech), " is currently not supported")
		}

		// If error is failBuild error, remember it and continue to next tech
		if e, ok := err.(*exec.ExitError); ok {
			if e.ExitCode() == coreutils.ExitCodeVulnerableBuild.Code {
				failBuildErr = err
				break
			}
		}

		if err != nil {
			return err
		}
	}
	return failBuildErr
}

func AuditMvnCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	xrAuditMvnCmd := java.NewAuditMavenCommand(*genericAuditCmd).SetInsecureTls(c.Bool(cliutils.InsecureTls))
	return commands.Exec(xrAuditMvnCmd)
}

func AuditGradleCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	xrAuditGradleCmd := java.NewAuditGradleCommand(*genericAuditCmd).SetExcludeTestDeps(c.Bool(cliutils.ExcludeTestDeps)).SetUseWrapper(c.Bool(cliutils.UseWrapper))
	return commands.Exec(xrAuditGradleCmd)
}

func AuditNpmCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	var setNpmArgs []string
	switch c.String("dep-type") {
	case "devOnly":
		setNpmArgs = []string{"--dev"}
	case "prodOnly":
		setNpmArgs = []string{"--prod"}
	}
	auditNpmCmd := npm.NewAuditNpmCommand(*genericAuditCmd).SetNpmArgs(setNpmArgs)
	return commands.Exec(auditNpmCmd)
}

func AuditGoCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	auditGoCmd := _go.NewAuditGoCommand(*genericAuditCmd)
	return commands.Exec(auditGoCmd)
}

func AuditPipCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	auditPipCmd := python.NewAuditPipCommand(*genericAuditCmd)
	return commands.Exec(auditPipCmd)
}

func AuditPipenvCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	auditPipenvCmd := python.NewAuditPipenvCommand(*genericAuditCmd)
	return commands.Exec(auditPipenvCmd)
}

func AuditNugetCmd(c *cli.Context) error {
	genericAuditCmd, err := createGenericAuditCmd(c)
	if err != nil {
		return err
	}
	auditNugetCmd := nuget.NewAuditNugetCommand(*genericAuditCmd)
	return commands.Exec(auditNugetCmd)
}

func createGenericAuditCmd(c *cli.Context) (*audit.AuditCommand, error) {
	auditCmd := audit.NewAuditCommand()
	err := validateXrayContext(c)
	if err != nil {
		return nil, err
	}
	serverDetails, err := createServerDetailsWithConfigOffer(c)
	if err != nil {
		return nil, err
	}
	format, err := commandsutils.GetXrayOutputFormat(c.String("format"))
	if err != nil {
		return nil, err
	}

	auditCmd.SetServerDetails(serverDetails).
		SetOutputFormat(format).
		SetTargetRepoPath(addTrailingSlashToRepoPathIfNeeded(c)).
		SetProject(c.String("project")).
		SetIncludeVulnerabilities(shouldIncludeVulnerabilities(c)).
		SetIncludeLicenses(c.Bool("licenses")).
		SetFail(c.BoolT("fail")).
		SetPrintExtendedTable(c.Bool(cliutils.ExtendedTable))

	if c.String("watches") != "" {
		auditCmd.SetWatches(strings.Split(c.String("watches"), ","))
	}
	return auditCmd, err
}

func ScanCmd(c *cli.Context) error {
	err := validateXrayContext(c)
	if err != nil {
		return err
	}
	serverDetails, err := createServerDetailsWithConfigOffer(c)
	if err != nil {
		return err
	}
	var specFile *spec.SpecFiles
	if c.IsSet("spec") {
		specFile, err = cliutils.GetFileSystemSpec(c)
	} else {
		specFile, err = createDefaultScanSpec(c, addTrailingSlashToRepoPathIfNeeded(c))
	}
	if err != nil {
		return err
	}
	err = spec.ValidateSpec(specFile.Files, false, false)
	if err != nil {
		return err
	}
	threads, err := cliutils.GetThreadsCount(c)
	if err != nil {
		return err
	}
	format, err := commandsutils.GetXrayOutputFormat(c.String("format"))
	if err != nil {
		return err
	}
	cliutils.FixWinPathsForFileSystemSourcedCmds(specFile, c)
	scanCmd := scan.NewScanCommand().SetServerDetails(serverDetails).SetThreads(threads).SetSpec(specFile).SetOutputFormat(format).
		SetProject(c.String("project")).SetIncludeVulnerabilities(shouldIncludeVulnerabilities(c)).
		SetIncludeLicenses(c.Bool("licenses")).SetFail(c.BoolT("fail")).SetPrintExtendedTable(c.Bool(cliutils.ExtendedTable))
	if c.String("watches") != "" {
		scanCmd.SetWatches(strings.Split(c.String("watches"), ","))
	}
	return commands.Exec(scanCmd)
}

// Scan published builds with Xray
func BuildScan(c *cli.Context) error {
	if c.NArg() > 2 {
		return cliutils.WrongNumberOfArgumentsHandler(c)
	}
	buildConfiguration := cliutils.CreateBuildConfiguration(c)
	if err := buildConfiguration.ValidateBuildParams(); err != nil {
		return err
	}

	serverDetails, err := createServerDetailsWithConfigOffer(c)
	if err != nil {
		return err
	}
	err = validateXrayContext(c)
	if err != nil {
		return err
	}
	format, err := commandsutils.GetXrayOutputFormat(c.String("format"))
	if err != nil {
		return err
	}
	buildScanCmd := scan.NewBuildScanCommand().SetServerDetails(serverDetails).SetFailBuild(c.BoolT("fail")).SetBuildConfiguration(buildConfiguration).
		SetIncludeVulnerabilities(c.Bool("vuln")).SetOutputFormat(format).SetPrintExtendedTable(c.Bool(cliutils.ExtendedTable))
	return commands.Exec(buildScanCmd)
}

func DockerScan(c *cli.Context, image string) error {
	if show, err := cliutils.ShowGenericCmdHelpIfNeeded(c, c.Args(), "dockerscanhelp"); show || err != nil {
		return err
	}
	if image == "" {
		return cli.ShowCommandHelp(c, "dockerscanhelp")
	}
	serverDetails, err := createServerDetailsWithConfigOffer(c)
	if err != nil {
		return err
	}
	containerScanCommand := scan.NewDockerScanCommand()
	fail, licenses, formatArg, project, watches, serverDetails, err := utils.ExtractDockerScanOptionsFromArgs(c.Args())
	if err != nil {
		return err
	}
	format, err := commandsutils.GetXrayOutputFormat(formatArg)
	if err != nil {
		return err
	}
	containerScanCommand.SetServerDetails(serverDetails).SetOutputFormat(format).SetProject(project).
		SetIncludeVulnerabilities(shouldIncludeVulnerabilities(c)).SetIncludeLicenses(licenses).
		SetFail(fail).SetPrintExtendedTable(c.Bool(cliutils.ExtendedTable))
	if watches != "" {
		containerScanCommand.SetWatches(strings.Split(watches, ","))
	}
	containerScanCommand.SetImageTag(c.Args().Get(1))
	return progressbar.ExecWithProgress(containerScanCommand, true)
}

func addTrailingSlashToRepoPathIfNeeded(c *cli.Context) string {
	repoPath := c.String("repo-path")
	if repoPath != "" && !strings.Contains(repoPath, "/") {
		// In case a only repo name was provided (no path) we are adding a trailing slash.
		repoPath += "/"
	}
	return repoPath
}

func createDefaultScanSpec(c *cli.Context, defaultTarget string) (*spec.SpecFiles, error) {
	return spec.NewBuilder().
		Pattern(c.Args().Get(0)).
		Target(defaultTarget).
		Recursive(c.BoolT("recursive")).
		Exclusions(cliutils.GetStringsArrFlagValue(c, "exclusions")).
		Regexp(c.Bool("regexp")).
		Ant(c.Bool("ant")).
		IncludeDirs(c.Bool("include-dirs")).
		BuildSpec(), nil
}

func createServerDetailsWithConfigOffer(c *cli.Context) (*coreconfig.ServerDetails, error) {
	return cliutils.CreateServerDetailsWithConfigOffer(c, true, "xr")
}

func shouldIncludeVulnerabilities(c *cli.Context) bool {
	// If no context was provided by the user, no Violations will be triggered by Xray, so include general vulnerabilities in the command output
	return c.String("watches") == "" && !isProjectProvided(c) && c.String("repo-path") == ""
}

func validateXrayContext(c *cli.Context) error {
	contextFlag := 0
	if c.String("watches") != "" {
		contextFlag++
	}
	if isProjectProvided(c) {
		contextFlag++
	}
	if c.String("repo-path") != "" {
		contextFlag++
	}
	if contextFlag > 1 {
		return errorutils.CheckErrorf("only one of the following flags can be supplied: --watches, --project or --repo-path")
	}
	return nil
}

func isProjectProvided(c *cli.Context) bool {
	return c.String("project") != "" || os.Getenv(coreutils.Project) != ""
}
