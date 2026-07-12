import re

with open('dependencies.go', 'r') as f:
    content = f.read()

repo_constructors = re.findall(r'(\w+):\s*repositories\.New(\w+)\(', content)
svc_constructors = re.findall(r'(\w+):\s*services\.New(\w+)\(', content)

with open('wire.go', 'w') as f:
    f.write('//go:build wireinject\n')
    f.write('// +build wireinject\n\n')
    f.write('package http\n\n')
    f.write('import (\n')
    f.write('\t"github.com/google/wire"\n')
    f.write(')\n\n')
    
    f.write('var repositorySet = wire.NewSet(\n')
    for field, name in repo_constructors:
        f.write(f'\trepositories.New{name},\n')
    f.write(')\n')
