import re

def move_dependencies():
    path = "internal/http/server.go"
    with open(path, "r") as f:
        content = f.read()

    # The block we want to extract starts at "type appRepositories struct {" and ends right before "func NewServer("
    # Let's find it.
    start_str = "type appRepositories struct {"
    end_str = "func NewServer("
    
    start_idx = content.find(start_str)
    end_idx = content.find(end_str)
    
    if start_idx == -1 or end_idx == -1:
        print("Could not find blocks")
        return
        
    extracted_block = content[start_idx:end_idx]
    
    # Save the extracted block to internal/http/dependencies.go
    deps_content = """package http

import (
	"context"
	"database/sql"

	"github.com/docker/docker/client"

	"vessel.dev/vessel/internal/core"
	"vessel.dev/vessel/internal/engine"
	"vessel.dev/vessel/internal/handlers"
	"vessel.dev/vessel/internal/repositories"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/vault"
)

""" + extracted_block

    with open("internal/http/dependencies.go", "w") as f:
        f.write(deps_content)
        
    # Remove the block from server.go
    new_content = content[:start_idx] + content[end_idx:]
    with open(path, "w") as f:
        f.write(new_content)

move_dependencies()
