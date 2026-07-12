import os
import glob
import re

replacements = {
    "Oauth": "Auth",
    "Profile": "Users",
    "ProjectSettings": "Projects",
    "S3-destinations": "Backups",
    "Services": "AppServices",
    "ServiceVariables": "AppServices",
    "EmailSettings": "Settings",
    "Team-invites": "Teams",
    "Ws": "Terminal",
    "Vercel": "Projects",
    "Domains": "Projects"
}

for root, _, files in os.walk("internal/handlers"):
    for file in files:
        if file.endswith(".go"):
            path = os.path.join(root, file)
            with open(path, "r") as f:
                content = f.read()
            
            changed = False
            for old, new in replacements.items():
                pattern = r"// @Tags\s+" + old + r"\b"
                if re.search(pattern, content):
                    content = re.sub(pattern, "// @Tags " + new, content)
                    changed = True
                    
            if changed:
                with open(path, "w") as f:
                    f.write(content)
                print(f"Updated tags in {path}")
