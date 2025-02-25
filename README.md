# Overview

The official API of sunnah.com for retrieving information about hadith collections.

# Getting started

Please follow the instructions below.

First create a local `.env.local` configuration file and update values as needed.
A sample file is provided at `.env.local.sample`.

Run manually:
```bash
git clone REPO
cd REPO
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
export FLASK_ENV=development FLASK_APP=main.py
flask run --host=0.0.0.0
```

Or alternatively use `docker-compose` which will give a full environment with a MySQL instance loaded with a sample dataset:

```bash
docker-compose up
```

* Use `--build` option to re-build.
* Use the `-d` option to run in detached mode.

You can then visit [localhost:5000](http://localhost:5000) to verify that it's running on your machine. Or, alternatively:

```bash
$ curl http://localhost:5000
```

## Deployment

Configuration files are located at `env.local` and `uwsgi.ini`.

A production ready uWSGI daemon (uwsgi socket exposed on port 5001) can be started with:

```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

## Routes

Visit https://sunnah.stoplight.io/docs/api/ for full API documentation.

## Linting and Formatting

`flake8` and `black` are used for code linting and formatting respectively. Before submitting pull requests, make sure black and flake8 is run against the code. Follow the instructions below for using `black` and `flake8`:

```sh
# goto repository root directory
# make sure the virtual environment is activated
black .
flake8 .
# fix any linting issues
# Then you are ready to submit your PR
```

To add more rules for linting and formatting, make changes to `.flake8` and `pyproject.toml` accordingly.

# Guidelines for Sending a Pull Request

1. Only change one thing at a time.
2. Don't mix a lot of formatting changes with logic change in the same pull request.
3. Keep code refactor and logic change in separate pull requests.
4. Squash your commits. When you address feedback, squash it as well. No one benefits from "addressed feedback" commit in the history.
5. Break down bigger changes into smaller separate pull requests.
6. If changing UI, attach a screenshot of how the changes look.
7. Reference the issue being fixed by adding the issue tag in the commit message.
8. Do not send a big change before first proposing it and getting a buy-in from the maintainer.
