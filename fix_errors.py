import os
import re

directory = 'internal/repositories'

for filename in os.listdir(directory):
    if not filename.endswith('.go'):
        continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        lines = f.readlines()

    out_lines = []
    changed = False
    
    current_func = None
    current_arg = None
    
    i = 0
    while i < len(lines):
        line = lines[i]
        
        # Track function name and first param after ctx
        func_match = re.search(r'func \([^)]+\)\s+(Get[A-Za-z0-9_]+)\([^,]+,\s*([a-zA-Z0-9_]+)\s', line)
        if func_match:
            current_func = func_match.group(1)
            current_arg = func_match.group(2)
            
            # e.g. GetUserByEmail -> Entity: User
            entity = current_func.replace('Get', '')
            if 'By' in entity:
                entity = entity.split('By')[0]
            if entity == '':
                entity = 'Entity'
        
        # Check for sql.ErrNoRows check
        if 'errors.Is(err, sql.ErrNoRows)' in line or 'err == sql.ErrNoRows' in line:
            # Only process if we are inside a Get function and returning nil, nil or similar
            if current_func and (i + 1) < len(lines):
                next_line = lines[i+1]
                if 'return nil, nil' in next_line or 'return nil, err' in next_line:
                    line = line.replace('if err == sql.ErrNoRows', 'if errors.Is(err, sql.ErrNoRows)')
                    out_lines.append(line)
                    
                    if current_arg:
                        # use current_arg if it's not a generic name, or if it is string we can use it
                        replacement = f'\t\treturn nil, utils.NewNotFoundError("{entity}", {current_arg})\n'
                    else:
                        replacement = f'\t\treturn nil, utils.NewNotFoundError("{entity}", "")\n'
                        
                    # Fix indentation
                    indent = next_line[:len(next_line) - len(next_line.lstrip())]
                    out_lines.append(indent + replacement.strip() + '\n')
                    i += 1
                    changed = True
                    i += 1
                    continue

        out_lines.append(line)
        i += 1

    if changed:
        # Add utils import
        content = "".join(out_lines)
        if '"vessel.dev/vessel/internal/utils"' not in content:
            # find last import
            if 'import (' in content:
                content = content.replace('import (\n', 'import (\n\t"vessel.dev/vessel/internal/utils"\n')
        
        with open(filepath, 'w') as f:
            f.write(content)
        print(f"Modified {filename}")
