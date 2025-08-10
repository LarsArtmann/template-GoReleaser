// GoReleaser Configuration Tool - Client-side JavaScript

// Configuration templates
const configTemplates = {
    basic: `# Basic GoReleaser configuration
project_name: my-project

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
    format_overrides:
    - goos: windows
      format: zip

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ incpatch .Version }}-next"

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'`,

    docker: `# GoReleaser configuration with Docker
project_name: my-project

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}

dockers:
  - image_templates:
    - "ghcr.io/username/{{ .ProjectName }}:{{ .Version }}"
    - "ghcr.io/username/{{ .ProjectName }}:latest"
    dockerfile: Dockerfile
    build_flag_templates:
    - "--pull"
    - "--label=org.opencontainers.image.created={{.Date}}"
    - "--label=org.opencontainers.image.title={{.ProjectName}}"
    - "--label=org.opencontainers.image.revision={{.FullCommit}}"
    - "--label=org.opencontainers.image.version={{.Version}}"

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ incpatch .Version }}-next"

changelog:
  sort: asc`,

    homebrew: `# GoReleaser configuration with Homebrew
project_name: my-project

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}

brews:
  - name: my-project
    tap:
      owner: username
      name: homebrew-tap
    homepage: "https://github.com/username/my-project"
    description: "My awesome project description"
    license: "MIT"

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ incpatch .Version }}-next"

changelog:
  sort: asc`
};

// Load configuration template
function loadTemplate(templateName) {
    const editor = document.getElementById('config-editor');
    if (editor && configTemplates[templateName]) {
        editor.value = configTemplates[templateName];
        showToast('Template loaded successfully', 'success');
    }
}

// Toast notification system
function showToast(message, type = 'info', duration = 3000) {
    // Remove existing toasts
    const existingToasts = document.querySelectorAll('.toast');
    existingToasts.forEach(toast => toast.remove());

    // Create new toast
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    // Show toast
    setTimeout(() => toast.classList.add('show'), 100);
    
    // Hide toast after duration
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 300);
    }, duration);
}

// Format bytes for display
function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// Format validation results
function formatValidationResult(result) {
    if (!result) return '';
    
    let html = '';
    
    if (result.valid) {
        html += '<div class="validation-success">✅ Configuration is valid</div>';
    } else {
        html += '<div class="validation-error">❌ Configuration has errors</div>';
    }
    
    if (result.issues && result.issues.length > 0) {
        html += '<div class="mt-2"><strong>Issues:</strong><ul class="list-disc list-inside mt-1">';
        result.issues.forEach(issue => {
            html += `<li class="text-red-600">${issue}</li>`;
        });
        html += '</ul></div>';
    }
    
    if (result.warnings && result.warnings.length > 0) {
        html += '<div class="mt-2"><strong>Warnings:</strong><ul class="list-disc list-inside mt-1">';
        result.warnings.forEach(warning => {
            html += `<li class="text-yellow-600">${warning}</li>`;
        });
        html += '</ul></div>';
    }
    
    return html;
}

// HTMX event handlers
document.addEventListener('DOMContentLoaded', function() {
    // Add loading indicators
    document.body.addEventListener('htmx:beforeRequest', function(evt) {
        const target = evt.target;
        if (target.tagName === 'BUTTON') {
            target.disabled = true;
            const originalText = target.textContent;
            target.setAttribute('data-original-text', originalText);
            target.innerHTML = '<span class="spinner"></span>' + originalText;
        }
    });
    
    document.body.addEventListener('htmx:afterRequest', function(evt) {
        const target = evt.target;
        if (target.tagName === 'BUTTON') {
            target.disabled = false;
            const originalText = target.getAttribute('data-original-text');
            if (originalText) {
                target.textContent = originalText;
                target.removeAttribute('data-original-text');
            }
        }
    });
    
    // Handle successful responses
    document.body.addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.successful) {
            // Show elements that were hidden
            const targetElement = document.getElementById(evt.detail.target.id);
            if (targetElement && targetElement.classList.contains('hidden')) {
                targetElement.classList.remove('hidden');
            }
        }
    });
    
    // Handle validation responses
    document.body.addEventListener('htmx:afterSettle', function(evt) {
        const target = evt.target;
        if (target.id === 'validation-results' || target.id === 'quick-validation') {
            try {
                const response = JSON.parse(evt.detail.xhr.responseText);
                if (response.valid !== undefined) {
                    target.innerHTML = formatValidationResult(response);
                }
            } catch (e) {
                // Response is already HTML, no need to format
            }
        }
        
        if (target.id === 'metrics-display') {
            try {
                const response = JSON.parse(evt.detail.xhr.responseText);
                if (response.memory) {
                    let html = '<h3 class="text-xl font-semibold mb-4">System Metrics</h3>';
                    html += '<div class="grid grid-cols-2 md:grid-cols-4 gap-4">';
                    html += `<div class="metric-card">
                        <div class="metric-value">${formatBytes(response.memory.alloc)}</div>
                        <div class="metric-label">Memory Allocated</div>
                    </div>`;
                    html += `<div class="metric-card">
                        <div class="metric-value">${formatBytes(response.memory.sys)}</div>
                        <div class="metric-label">System Memory</div>
                    </div>`;
                    html += `<div class="metric-card">
                        <div class="metric-value">${response.goroutines}</div>
                        <div class="metric-label">Goroutines</div>
                    </div>`;
                    html += `<div class="metric-card">
                        <div class="metric-value">${response.memory.num_gc}</div>
                        <div class="metric-label">GC Cycles</div>
                    </div>`;
                    html += '</div>';
                    target.innerHTML = html;
                }
            } catch (e) {
                // Response is already HTML
            }
        }
    });
    
    // Auto-save functionality (save config every 30 seconds if modified)
    let configModified = false;
    let autoSaveInterval;
    
    const configEditor = document.getElementById('config-editor');
    if (configEditor) {
        configEditor.addEventListener('input', function() {
            configModified = true;
            if (!autoSaveInterval) {
                autoSaveInterval = setInterval(function() {
                    if (configModified) {
                        // Trigger auto-save
                        htmx.trigger(configEditor.closest('form'), 'submit');
                        configModified = false;
                        showToast('Configuration auto-saved', 'info', 2000);
                    }
                }, 30000);
            }
        });
    }
});

// Error handling
document.body.addEventListener('htmx:responseError', function(evt) {
    showToast('Request failed: ' + evt.detail.xhr.statusText, 'error');
});

document.body.addEventListener('htmx:timeout', function(evt) {
    showToast('Request timed out', 'error');
});