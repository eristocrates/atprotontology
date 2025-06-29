import os
import sys
import requests
from urllib.parse import urljoin, urlparse, urldefrag
from bs4 import BeautifulSoup

if len(sys.argv) != 3:
    print("Usage: python shallow_crawler.py <start_url> <output_dir>")
    sys.exit(1)

start_url = sys.argv[1]
output_dir = sys.argv[2]
visited = set()

os.makedirs(output_dir, exist_ok=True)

base_domain = urlparse(start_url).netloc

def sanitize_filename(url):
    parsed = urlparse(url)
    path = parsed.path.rstrip('/')
    if not path or path == '/':
        filename = 'index'
    else:
        filename = path.lstrip('/').replace('/', '_')
    return f"{filename}.html"

def download_and_save(url, filename):
    try:
        print(f"Downloading: {url}")
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        filepath = os.path.join(output_dir, filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(response.text)
        visited.add(url)
    except Exception as e:
        print(f"Failed to download {url}: {e}")

# Step 1: Download main page
main_filename = sanitize_filename(start_url)
download_and_save(start_url, main_filename)

# Step 2: Follow same-domain links (1 level deep)
main_path = os.path.join(output_dir, main_filename)
with open(main_path, encoding='utf-8') as f:
    soup = BeautifulSoup(f, 'html.parser')
    for link in soup.find_all('a', href=True):
        href = link['href']
        if href.startswith('#'):
            continue
        full_url = urljoin(start_url, href)
        clean_url, _ = urldefrag(full_url)

        # Only visit same-domain links
        if clean_url not in visited and urlparse(clean_url).netloc == base_domain:
            download_and_save(clean_url, sanitize_filename(clean_url))
