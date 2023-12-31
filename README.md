# Coding Bot
A GPT bot that tries to be a useful coding assistant by automating your common tasks.


> ⚠️ **NOTE** ⚠️
>
> This project is still very much new and rough. It has A LOT of edge cases... But there's a glimmer of something fun and useful here.

> ⚠️ **CAUTION** ⚠️
>
> There's an "API" app, but that's not currently attached to the new way of doing things. Main focus is on the project structure & CLI tool. See roadmap for more information.

## Project details
The goal of this project is to have a coding bot that can answer your coding questions and understand your current project.

## Features
### Plugins
Create unique and powerful plugins that leverage OpenAI Function Calling. Developers can create plugins that can perform a wide range of actions with the results from a chatbot, such as creating files, appending to existing files, initiating new projects, and more.

Available plugins:
- [x] Coding style guide
- [x] Creating or updating a file

## Installation
### Requirements 
- Go 1.21 (uses new features)
- OpenAI API token 

```shell
$ OPENAI_TOKEN=[REPLACE_WITH_TOKEN] >> backend/.env # only necessary if you are developing locally
$ make build/cli # build the cli tool
```

#### Running, Prompts and commands
```shell
$ export OPENAI_TOKEN=[REPLACE_WITH_TOKEN] 
$ bin/cli -p link-to-your-prompt.md -i link-to-optional-input-file.ts
```

## Challenges

### Excluding necessary code
Sometimes the bot will return incomplete file code. For example, if I request the bot to convert some TS into a Go struct and to then create a file out of that output, it will usually return the correct type, but it may choose not to include the package name for the file; which is required in a Go file.

#### Solution / workaround
I'm working on two approaches. 
- [x] Add explicit instructions in the prompts to always include necessary information
- [ ] Add file validation checks


### Output changes
Differences in output from one run to another. Ie. Sometimes a function will be called `GetName()` and another time it will be called `GetXXXName()`. This can be problematic if you reference the first function and it gets updated to another name.

#### Solution / workaround
This problem is a little harder to get around until there's a way to map a project to the bot to gain context. 

Until then, 
- A bot should never completely overwrite a file, unless explicitly told to.
- A bot should never overwrite the contents of the file, instead it should append it.

### Context Length
#### Problem
Since LLMs do not possess inherent memory and rely on tokens as their source of memory, it is common to provide an agent with sufficient information in the system prompt or messages. This enables the agent to build context and pattern recognition necessary to respond appropriately to a given question. However, it is important to note that most LLMs, including OpenAI models used in this repository, have a token limit of 4,000. It may be surprising how quickly these tokens are consumed when dealing with files.

#### Solution / workaround
Like the previous problem this is hard to get around since it's a current technical limitation with most LLMs. But there are some ways to reduce the amount of tokens through some prompt engineering tricks and leveraging embeddings.

Techniques and strategies being investigated.
- [ ] Improve conversations via a "[select and query](https://martinfowler.com/articles/building-boba.html#carry-context)" strategy
- [ ] Leverage embeddings with a vector database

### Understanding project structure
#### Problem
Get the assistant to understand project structure, specifically, what the package name should be for a file. The bot always wants to give it "package main" when other files in that directory are "package XXX".

#### Solution / workaround
Give the bot a mock conversation that is added to the system prompt.

example:
system prompt: below is an example conversation from the users project structure

```markdown
[//]: # (project structure)
| Name | Filepath | Is Directory | Size | Last Modified | Metadata |
| ---- | -------- | ------------ | ---- | ------------- | ------- |
| . | . | true | 224 | 2023-06-26 14:15:05.789880218 -0700 PDT | null |
| cmd | cmd | true | 128 | 2023-06-22 13:22:13.2712506 -0700 PDT | null |
| api | cmd/api | true | 288 | 2023-06-26 16:16:58.213583944 -0700 PDT | null |
| errors.go | cmd/api/errors.go | false | 1905 | 2023-06-21 20:13:35.06768917 -0700 PDT | {"package_name":"package main"} |
| read.go | internal/file/read.go | false | 216 | 2023-06-27 11:31:18.045689437 -0700 PDT | {"package_name":"package file"} |
| write.go | internal/file/write.go | false | 297 | 2023-06-25 20:23:02.117814613 -0700 PDT | {"package_name":"package file"} |

user: what's the file path for the errors file in the api directory?
assistant: cmd/api/errors.go
user: What is the package name for the api directory?
assistant: package main
user: What is the package name for the file directory?
assistant: package file
user: create foo.go in the file package
assistant: I created the file internal/file/foo.go here is its contents

package file

type foo struct{}

user: create bar.go in the api directory
assistant: I created the file cmd/api/bar.go here is its contents
 
package main

type bar struct{}
```

Note how I am priming the bot to understand the projects file structure and file paths, and how I prime it to use the metadata for the package name. This back and forth conversation example has worked really well so far. 

## Roadmap
- [x] Internal plugin implementation
- [x] OpenAI function calling
- [x] Enhanced plugin capabilities
  - [x] Inject system prompts: this allows the developer to prime the bot with examples and additional context on what they or the plugin will do. 
  - [x] Add pre-flight messages / conversation: this appears to improve the bots context of what is being asked of them by "faking" a conversation.
- [x] CLI chat ux
  - [x] remove default invocation
  - [x] cleaning printing
  - [x] better message response
  - [x] chat loop
- [ ] Additional plugin Functionality
  - [ ] Handle multiple files 
  - [ ] Create specific kind of files
    - tests
    - services
    - repository
    - views or project
- [ ] Context / Token utilization improvements
  - [ ] Implement prompt optimization strategies 
    -  Select and query
  - [ ] Implement embeddings / vector database
- [ ] API implementation
- [ ] Web application