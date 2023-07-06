# Mythic GoLang Scripting Interface

The `Mythic_Go_Scripting` package creates a way to interact and control a Mythic instance programmatically and is the GoLang implementation of the Python `mythic` package. Mythic is a Command and Control (C2) framework for Red Teaming. The code is on GitHub (https://github.com/its-a-feature/Mythic) and the Mythic project's documentation is on GitBooks (https://docs.mythic-c2.net).
## Installation

You can install the mythic scripting interface from github using `go get`:

```
go get github.com/antman1p/Mythic_Go_Scripting
```
Then import the package into your git project file from where you intend to interact with Mythic.  For example:
```
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	mythic "github.com/antman1p/Mythic_Go_Scripting"
)
```
