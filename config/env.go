package config

import "os"

type env struct {
	key      string
	variable *string
}

var envs = []env{
	{"PORT", &Port},
	{"MONGO_URI", &MongoUri},
	{"ACCESS_SECRET", &AccessSecret},
	{"REFRESH_SECRET", &RefreshSecret},
	{"SMTP_IDENTITY", &SmtpIdentity},
	{"SMTP_USERNAME", &SmtpUsername},
	{"SMTP_PASSWORD", &SmtpPassword},
	{"SMTP_HOST", &SmtpHost},
	{"SMTP_PORT", &SmtpPort},
}

func loadEnv(env env) bool {
	value, ok := os.LookupEnv(env.key)
	if ok {
		*env.variable = value
	}

	return ok
}
