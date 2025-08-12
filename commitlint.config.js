// Commitlint configuration for GoReleaser Template Project
// Enforces conventional commit message format
// https://conventionalcommits.org/

module.exports = {
  // Extend the conventional configuration
  extends: ['@commitlint/config-conventional'],
  
  // Custom rules
  rules: {
    // Type enum - allowed commit types
    'type-enum': [
      2, // Error level
      'always',
      [
        // Main types
        'feat',     // New feature
        'fix',      // Bug fix
        'docs',     // Documentation changes
        'style',    // Code style changes (formatting, etc.)
        'refactor', // Code refactoring
        'perf',     // Performance improvements
        'test',     // Test additions or modifications
        'chore',    // Maintenance tasks
        
        // Additional types for this project
        'build',    // Build system or external dependency changes
        'ci',       // CI/CD configuration changes
        'revert',   // Revert a previous commit
        'security', // Security fixes
        'deps',     // Dependency updates
        'config',   // Configuration changes
        'docker',   // Docker-related changes
        'release',  // Release commits
        'hotfix',   // Critical hotfixes
        'init',     // Initial commit
        'wip',      // Work in progress (for development branches)
      ]
    ],
    
    // Type case - must be lowercase
    'type-case': [2, 'always', 'lowercase'],
    
    // Type empty - type is required
    'type-empty': [2, 'never'],
    
    // Scope case - must be lowercase
    'scope-case': [2, 'always', 'lowercase'],
    
    // Scope enum - allowed scopes for this project
    'scope-enum': [
      1, // Warning level (not error)
      'always',
      [
        // Core components
        'cli',
        'server',
        'api',
        'web',
        'ui',
        
        // Services
        'config',
        'validation',
        'license',
        'container',
        'template',
        
        // Infrastructure
        'docker',
        'ci',
        'build',
        'release',
        'deploy',
        
        // Development
        'dev',
        'test',
        'bench',
        'lint',
        'fmt',
        
        // Documentation
        'docs',
        'readme',
        'guide',
        'examples',
        
        // Dependencies
        'deps',
        'mod',
        'vendor',
        
        // Security
        'security',
        'auth',
        'crypto',
        
        // Performance
        'perf',
        'cache',
        'memory',
        
        // Project structure
        'cmd',
        'internal',
        'pkg',
        'scripts',
        'assets',
        'templates',
      ]
    ],
    
    // Subject case - sentence case (first letter lowercase)
    'subject-case': [2, 'always', 'sentence-case'],
    
    // Subject empty - subject is required
    'subject-empty': [2, 'never'],
    
    // Subject full stop - no period at the end
    'subject-full-stop': [2, 'never', '.'],
    
    // Subject max length
    'subject-max-length': [2, 'always', 72],
    
    // Subject min length
    'subject-min-length': [2, 'always', 3],
    
    // Header max length (type + scope + subject)
    'header-max-length': [2, 'always', 100],
    
    // Body leading blank - require blank line before body
    'body-leading-blank': [2, 'always'],
    
    // Body max line length
    'body-max-line-length': [1, 'always', 100],
    
    // Footer leading blank - require blank line before footer
    'footer-leading-blank': [2, 'always'],
    
    // Footer max line length
    'footer-max-line-length': [1, 'always', 100],
  },
  
  // Custom prompt configuration
  prompt: {
    questions: {
      type: {
        description: 'Select the type of change that you\'re committing:',
        enum: {
          feat: {
            description: 'A new feature',
            title: 'Features'
          },
          fix: {
            description: 'A bug fix',
            title: 'Bug Fixes'
          },
          docs: {
            description: 'Documentation only changes',
            title: 'Documentation'
          },
          style: {
            description: 'Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)',
            title: 'Styles'
          },
          refactor: {
            description: 'A code change that neither fixes a bug nor adds a feature',
            title: 'Code Refactoring'
          },
          perf: {
            description: 'A code change that improves performance',
            title: 'Performance Improvements'
          },
          test: {
            description: 'Adding missing tests or correcting existing tests',
            title: 'Tests'
          },
          build: {
            description: 'Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)',
            title: 'Builds'
          },
          ci: {
            description: 'Changes to our CI configuration files and scripts (example scopes: Circle, BambooHR, Jenkins)',
            title: 'Continuous Integrations'
          },
          chore: {
            description: 'Other changes that don\'t modify src or test files',
            title: 'Chores'
          },
          revert: {
            description: 'Reverts a previous commit',
            title: 'Reverts'
          },
          security: {
            description: 'Security fixes or improvements',
            title: 'Security'
          },
          deps: {
            description: 'Dependency updates',
            title: 'Dependencies'
          },
          config: {
            description: 'Configuration changes',
            title: 'Configuration'
          },
          docker: {
            description: 'Docker-related changes',
            title: 'Docker'
          },
          release: {
            description: 'Release commits',
            title: 'Releases'
          },
          hotfix: {
            description: 'Critical hotfixes',
            title: 'Hotfixes'
          }
        }
      },
      scope: {
        description: 'What is the scope of this change (e.g. cli, server, api, docs)?'
      },
      subject: {
        description: 'Write a short, imperative tense description of the change'
      },
      body: {
        description: 'Provide a longer description of the change'
      },
      isBreaking: {
        description: 'Are there any breaking changes?',
        default: false
      },
      breakingBody: {
        description: 'A BREAKING CHANGE commit requires a body. Please enter a longer description of the commit itself'
      },
      breaking: {
        description: 'Describe the breaking changes'
      },
      isIssueAffected: {
        description: 'Does this change affect any open issues?',
        default: false
      },
      issuesBody: {
        description: 'If issues are closed, the commit requires a body. Please enter a longer description of the commit itself'
      },
      issues: {
        description: 'Add issue references (e.g. "fix #123", "re #123")'
      }
    }
  },
  
  // Custom parser preset
  parserPreset: {
    parserOpts: {
      // Header pattern for parsing commits
      headerPattern: /^(\w*)(?:\(([^)]*)\))?: (.*)$/,
      headerCorrespondence: ['type', 'scope', 'subject'],
      
      // Reference actions (for linking to issues)
      referenceActions: [
        'close',
        'closes', 
        'closed',
        'fix',
        'fixes',
        'fixed',
        'resolve',
        'resolves',
        'resolved'
      ],
      
      // Issue prefixes
      issuePrefixes: ['#'],
      
      // Note keywords
      noteKeywords: ['BREAKING CHANGE', 'BREAKING CHANGES'],
      
      // Field pattern
      fieldPattern: /^-(.*?)-$/,
      
      // Revert pattern
      revertPattern: /^(?:Revert|revert:)\s"?([\s\S]+?)"?\s*This reverts commit (\w*)\./i,
      revertCorrespondence: ['header', 'hash']
    }
  },
  
  // Ignore patterns - commits that should be ignored
  ignores: [
    // Merge commits
    (commit) => commit.includes('Merge'),
    
    // Revert commits (handled separately)
    (commit) => commit.includes('Revert'),
    
    // Automated commits
    (commit) => commit.includes('[automated]'),
    (commit) => commit.includes('[skip ci]'),
    (commit) => commit.includes('[ci skip]'),
    
    // Dependabot commits
    (commit) => commit.includes('dependabot'),
    
    // Release commits from tools
    (commit) => commit.includes('chore(release)'),
    (commit) => commit.includes('Release '),
  ],
  
  // Default ignores (can be overridden)
  defaultIgnores: true,
  
  // Help URL
  helpUrl: 'https://github.com/conventional-changelog/commitlint/#what-is-commitlint',
  
  // Custom formatters for different output formats
  formatter: '@commitlint/format',
  
  // Custom plugins (if any)
  plugins: [
    // Add any custom plugins here if needed
    // Example: '@commitlint/plugin-function-rules'
  ]
};

/*
Examples of valid commit messages for this configuration:

feat(cli): add new validation command
fix(server): resolve memory leak in handler
docs(readme): update installation instructions
style(web): format CSS files according to standards
refactor(config): simplify configuration loading logic
perf(api): optimize database queries
test(validation): add unit tests for environment validation
chore(deps): update Go dependencies
build(docker): optimize Docker image size
ci(github): add automated testing workflow
security(auth): fix authentication bypass vulnerability
deps(mod): update goreleaser to v1.0.0

Examples of commit messages with body and footer:

feat(license): add support for custom license templates

This change allows users to provide their own license templates
in addition to the built-in ones. Templates are loaded from the
templates/licenses directory and validated before use.

Closes #123
Breaking change: Custom template format has changed

fix(validation): handle empty configuration files gracefully

Previously, the application would crash when encountering empty
configuration files. Now it provides a helpful error message and
suggests using the 'config init' command.

Fixes #456
*/