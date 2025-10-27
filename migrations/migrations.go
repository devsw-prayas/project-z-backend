package migrations

func SetupMigration() {
	UsersMigration()
	ProblemsMigration()
	SubmissionsMigration()
}
