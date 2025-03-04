# Citizenship Tracker CLI Homebrew Tap

This is the Homebrew Tap for the [Citizenship Tracker CLI](https://github.com/mendesbarreto/citizenship-tracker-cli).

## Installation

```bash
brew tap mendesbarreto/citizenship
brew install citizen
```

## Usage

```bash
# Check the application version
citizen --version

# Run the application
citizen status
```
```

## Step 6: Create a GitHub secret for tap repository access

1. Go to your GitHub repository settings
2. Navigate to "Secrets and Variables" > "Actions"
3. Create a new repository secret:
   - Name: `TAP_REPO_TOKEN`
   - Value: A GitHub personal access token with repo scope

## Step 7: Commit and push your changes

```bash
git add .
git commit -m "Add version support and GitHub Actions workflow"
git push
```

## Step 8: Create and push a tag to trigger a release

```bash
git tag -a v0.0.1 -m "Initial release"
git push origin v0.0.1
