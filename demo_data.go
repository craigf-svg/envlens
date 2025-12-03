package main

func getDemoEnvVars() []string {
	return []string{
		"PATH=/usr/local/bin:/usr/bin:/bin",
		"HOME=/home/developer",
		"USER=developer",
		"TERM=xterm-256color",
		"SHELL=/bin/bash",
		"EDITOR=vim",
		"LANG=en_US.UTF-8",
		"GOPATH=/home/developer/go",
		"GOROOT=/usr/local/go",
		"NODE_ENV=development",
		"NPM_CONFIG_PREFIX=/home/developer/.npm-global",
		"DOCKER_HOST=unix:///var/run/docker.sock",
		"AWS_REGION=us-east-1",
		"AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
		"AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"DATABASE_URL=postgres://user:pass@localhost:5432/mydb",
		"REDIS_URL=redis://localhost:6379",
		"API_KEY=sk-demo-1234567890abcdef",
		"JWT_SECRET=super-secret-jwt-token-for-demo",
		"LOG_LEVEL=debug",
	}
}

func getDemoLocalEnvVars() []string {
	return []string{
		"DATABASE_URL=postgres://admin:localpass@localhost:5432/devdb",
		"API_KEY=sk-local-dev-key-12345",
		"REDIS_URL=redis://localhost:6379/0",
		"JWT_SECRET=local-dev-jwt-secret",
		"DEBUG=true",
		"PORT=3000",
		"SMTP_HOST=localhost",
		"SMTP_PORT=1025",
	}
}
