package constants

const (
	InterestedFiles = `
	I will give a list of files from a github repository.
	and a repository name.
	The format will be {path_to_file}:{size}
	Tell me which files would be the most useful to send for you to get a good idea of the project.
	Select the fewest possible number of files that would give you a good idea of the project.
	ONLY SELECT FILES WITH AN EXTENSION THAT YOU CAN READ AND TYPICALLY DO NOT CONTAIN VERBOSE DATA.
	RETURN THE FULL PATH TO THE FILE.
	DO NOT RETURN ANY ADDITIONAL COMMENTARY OR GRAMMER ONLY THE LIST OF FILES.
	`

	InterestedFilesInput = `
	REPOSITORY NAME : %s
	%s
	`

	ContentMessage = `
	I will Give you a series of filenames and excerpts of the contents of those files.
	followed by a repository name.
	Provide me a synopsis of the project in the style of a readme.
	You should not appear to be guessing , speak with authority it does not matter if you are incorrect do not say what your assertions are based on or reference anything you used to generate the opion, simply speak as if you understand the purpose of this repository.
	MY CLIENTS HAVE MAJOR DISABILITIES AND IT IS A STRAIN FOR THEM TO REPROCESS THIS REQUEST PLEASE FOLLOW THE INSTRUCTIONS.
	Return the sysnopis in the form of a markdown blog post with appropriate title, formatting and some emojis.
	DO NOT RETURN ANY ADDITIONAL COMMENTARY OR GRAMMER ONLY THE SYNOPSIS.
	%s
	`
)
