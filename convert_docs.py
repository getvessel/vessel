import os
import re
import shutil

AEROPLANE_DOCS_DIR = '/home/eminisolomon/Dev/aeroplane/website/src/content/docs/docs'
VESSL_DOCS_DIR = '/home/eminisolomon/Dev/TechX/vessl/docs/src/content/docs'

# Files to ignore (Vessl already has equivalents, or user explicitly excluded)
IGNORE_FILES = {
    'index.mdx',
    'getting-started/install.md', # We have getting-started.md
    'getting-started/first-project.md', # We have tutorial.md
}

def convert_content(content: str) -> str:
    # 1. Basic terminology replacements
    content = content.replace('Aeroplane', 'Vessl')
    content = content.replace('aeroplane', 'vessl')
    content = content.replace('Caddy', 'Traefik')
    content = content.replace('caddy', 'traefik')
    content = content.replace('CADDY', 'TRAEFIK')
    content = content.replace('4310', '8080')
    content = content.replace('.vessl', '.vessl')
    content = content.replace('AEROPLANE_', 'VESSL_')
    
    # 2. Inject components import if we use Starlight components
    # We will replace > Note: or > Warning: with <Aside>
    has_aside = False
    
    def replacer_note(match):
        nonlocal has_aside
        has_aside = True
        inner_text = match.group(1).strip()
        return f'<Aside type="note">\n{inner_text}\n</Aside>'

    def replacer_tip(match):
        nonlocal has_aside
        has_aside = True
        inner_text = match.group(1).strip()
        return f'<Aside type="tip">\n{inner_text}\n</Aside>'

    def replacer_warning(match):
        nonlocal has_aside
        has_aside = True
        inner_text = match.group(1).strip()
        return f'<Aside type="caution">\n{inner_text}\n</Aside>'

    content = re.sub(r'>\s*\*\*Note\*\*:\s*(.*)', replacer_note, content)
    content = re.sub(r'>\s*\*\*Tip\*\*:\s*(.*)', replacer_tip, content)
    content = re.sub(r'>\s*\*\*Warning\*\*:\s*(.*)', replacer_warning, content)

    # 3. Handle <Steps> for numbered lists that are standalone
    # A bit complex to regex, we can just leave numbered lists as is, or attempt a naive wrap.
    # We'll just stick to standard markdown for lists, it's fine. Vessl uses <Steps> manually for complex ones.
    
    if has_aside:
        # insert import after frontmatter
        parts = content.split('---', 2)
        if len(parts) == 3:
            parts[2] = f"\nimport {{ Aside }} from '@astrojs/starlight/components';\n" + parts[2]
            content = "---".join(parts)
            
    return content

def main():
    for root, dirs, files in os.walk(AEROPLANE_DOCS_DIR):
        for file in files:
            if not file.endswith('.md'):
                continue
                
            full_path = os.path.join(root, file)
            rel_path = os.path.relpath(full_path, AEROPLANE_DOCS_DIR)
            
            if rel_path in IGNORE_FILES:
                continue
            
            dest_rel_path = rel_path.replace('aeroplane-bundles.md', 'vessl-bundles.md')
            dest_path = os.path.join(VESSL_DOCS_DIR, dest_rel_path)
            os.makedirs(os.path.dirname(dest_path), exist_ok=True)
            
            with open(full_path, 'r', encoding='utf-8') as f:
                content = f.read()
                
            new_content = convert_content(content)
            
            with open(dest_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            
            print(f"Converted: {rel_path}")

if __name__ == '__main__':
    main()
