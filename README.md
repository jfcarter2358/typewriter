# Typewriter

## About

Typewriter is a framework for hosting websites in a simple manner with Markdown being the main source of content.

## Getting Started

Create a git repo with the following structure (the only ones that _need_ to be there are `.theme` with `template.html` inside of it:

```
.theme
|-- template.html
page_1.md
page_2.md
path_1
|-- page_1_1.md
|-- path_1_1
    |-- page_1_1_1.md 
```

Push this repo up somewhere and then run Typewriter with the following config:

| ENV Var    | Example value                                                                                                                                               | Description                                              |
|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------|
| `PORT`     | `8080`                                                                                                                                                      | The port for Typewriter to listen on                              |
| `REDIRECT` | `page_1`                                                                                                                                                    | The landing page path for your site                          |
| `GIT`      | `{"url": "<your url>","username": "<your username>","password": "<your password>","branch": "<your branch>","email": "<your email>","name": "<your name>"}` | Git configuration for the repo containing your site contents |

Typewriter will then pull in the repo and serve your site with paths reflecting the file structure/file names.

# Configuring reverse proxies

> TODO

# Markdown file metadata

> TODO
