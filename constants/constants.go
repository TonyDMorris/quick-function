package constants

const (
	InterestedFiles = `
	I will give a list of files from a github repository.
	and a repository name.
	The format will be {path_to_file}:{size}
	Tell me which files would be the most useful to send for you to get a good idea of the project.
	DO NOT RETURN ANY ADDITIONAL COMMENTARY OR GRAMMER ONLY THE LIST OF FILES.
	`

	InterestedFilesInput = `
	REPOSITORY NAME : %s
	%s
	`
)
