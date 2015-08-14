package action

type DirectorInfo struct {
	UUID string
	URL  string
}

func NewDirectorInfo(directorURL string, cliRunner *CliRunner) (DirectorInfo, error) {
	err := cliRunner.TargetAndLogin(directorURL)
	if err != nil {
		return DirectorInfo{}, err
	}

	uuid, err := cliRunner.RunWithOutput("status", "--uuid")
	if err != nil {
		return DirectorInfo{}, err
	}

	return DirectorInfo{
		UUID: uuid,
		URL:  directorURL,
	}, nil
}
