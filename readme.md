# Docker Store Publish

This tools has been build to help people to manage their Docker store products.

```bash
DSP - Docker store publisher

Usage:
  dsp [flags]
  dsp [command]

Available Commands:
  help        Help about any command
  submit      Submit available image.
  version     Display version

Flags:
      --dry-run               Dry run mode. (default true)
  -h, --help                  help for dsp
      --password string       Docker password.
      --project string        Docker hub project.
      --publisher-id string   Docker publisher ID.
      --repository string     Docker hub repository.
      --username string       Docker username.

Use "dsp [command] --help" for more information about a command.
```

All flags are configurable with environment variables.

```bash
# This environment variable can be set if you want to load DOCKER_* env varibles.
DSP_ENV_FILE_PATH=/keybase/private/mmatur/dsp/.env

DOCKER_USER=username
DOCKER_PASSWORD=password
DOCKER_PUBLISHER_ID=123456789-123456-123456
DOCKER_REPOSITORY=mmatur
DOCKER_PROJECT=dsp
```

## Usage

- [dsp](docs/dsp.md)
- [dsp submit](docs/dsp_submit.md)
- [dsp version](docs/dsp_version.md)
