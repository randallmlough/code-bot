# Coding Bot
A GPT bot that tries to be a useful coding assistant.


> ⚠️ **NOTE** ⚠️
>
> This project is still very much new and rough. It has A LOT of edge cases... But there's a glimmer of something fun and useful here.

> ⚠️ **CAUTION** ⚠️
>
> There's an "API" app, but that's not currently attached to the new way of doing things. Main focus is on the project structure & CLI tool. See roadmap for more information.

## Project details
The goal of this project is to have a coding bot that can answer your coding questions and understand your current project.

## Features
Leverages OpenAI Function calling by its powerful plugin implementation. Developers can create plugins that do almost anything they want to do with the results from a chat bot.

Available plugins:
- [x] Coding style guide
- [x] Creating or updating a file

### Requirements 
- Go 1.21 (uses new features)
- OpenAI API token 

### Installation
```shell
$ OPENAI_TOKEN=[REPLACE_WITH_TOKEN] >> backend/.env # only necessary if you are developing locally
$ make build/cli # build the cli tool
$ OPENAI_TOKEN=[REPLACE_WITH_TOKEN] bin/cli # to run the cli tool
```

#### Prompts and commands

Edit files in `testdata/` to tell the cli tool what to do. 
> ⚠️ **NOTE** ⚠️
>
> This behavior will be deprecated very soon in favor of a much more organic and fluid process. This has just been the easiest way for me to build and test the bot.

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

### Understanding project structure
#### Problem
Get the assistant to understand project structure, specifically, what the package name should be for a file. Always wants to give it "package main" when other files in that directory are "package XXX".

#### Solution / workaround
Give the bot a mock conversational example 
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
```
user: what's the file path for the errors file in the api directory?
assistant: cmd/api/errors.go
user: What is the package name for the api directory?
assistant: package main
user: What is the package name for the file directory?
assistant: package file
user: create foo.go in the file package
assistant: I created the file internal/file/foo.go here is its contents
```go 
package file

type foo struct{}
```
user: create bar.go in the api directory
assistant: I created the file cmd/api/bar.go here is its contents
```go 
package main

type bar struct{}
```

## Roadmap
- [x] OpenAI function calling
- [x] Internal plugin implementation
- [x] Enhanced plugin capabilities
  - [x] Inject system prompts: this allows the developer to prime the bot with examples and additional context on what they or the plugin will do. 
  - [x] Add pre-flight messages / conversation: this appears to improve the bots context of what is being asked of them by "faking" a conversation.
- [ ] CLI chat loop
- [ ] API implementation
- [ ] Web application