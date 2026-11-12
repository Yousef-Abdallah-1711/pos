param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('Up', 'Seed', 'Test', 'Down')]
    [string]$Action
)

$ErrorActionPreference = 'Stop'
$project = 'pos-p2a-ws05-e2e-20261103'
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$composeFile = Join-Path $repoRoot 'docker-compose.phase2a-e2e.yml'
$seedFile = Join-Path $PSScriptRoot 'seed.ps1'
$composeFileFullPath = (Resolve-Path $composeFile).Path

function Invoke-Compose([string[]]$Arguments) {
    & docker compose -p $project -f $composeFile @Arguments
    if ($LASTEXITCODE -ne 0) { throw "Docker Compose action failed: $($Arguments -join ' ')" }
}

function Assert-DisposableProject {
    $containerIDs = @(docker ps -aq --filter "label=com.docker.compose.project=$project")
    if ($LASTEXITCODE -ne 0) { throw 'Could not inspect Docker project containers.' }
    if ($containerIDs.Count -eq 0) { throw "No containers found for expected disposable project $project." }

    foreach ($containerID in $containerIDs) {
        $labels = docker inspect --format '{{json .Config.Labels}}' $containerID | ConvertFrom-Json
        if ($LASTEXITCODE -ne 0 -or $labels.'com.docker.compose.project' -ne $project -or $labels.'com.docker.compose.project.config_files' -ne $composeFileFullPath) {
            throw "Refusing teardown: container $containerID is not owned by this exact Compose project and file."
        }
    }

    foreach ($service in @('postgres', 'api')) {
        $container = docker ps -q --filter "label=com.docker.compose.project=$project" --filter "label=com.docker.compose.service=$service"
        if ($LASTEXITCODE -ne 0 -or -not $container) { throw "Refusing teardown: expected isolated $service service is not running." }
        $labels = docker inspect --format '{{json .Config.Labels}}' $container | ConvertFrom-Json
        if ($LASTEXITCODE -ne 0 -or $labels.'com.pos.disposable' -ne 'true') { throw "Refusing teardown: the $service service lacks the disposable label." }
    }
}

Push-Location $repoRoot
try {
    switch ($Action) {
        'Up' {
            Invoke-Compose @('up', '-d', '--build')
            Invoke-Compose @('ps')
        }
        'Seed' {
            & $seedFile
            if ($LASTEXITCODE -ne 0) { throw 'Disposable fixture seed failed.' }
        }
        'Test' {
            $apiContainer = "$project-api-1"
            $dashboardContainer = "$project-dashboard-1"
            foreach ($container in @($apiContainer, $dashboardContainer)) {
                $labels = docker inspect --format '{{json .Config.Labels}}' $container | ConvertFrom-Json
                if ($LASTEXITCODE -ne 0 -or $labels.'com.docker.compose.project' -ne $project -or $labels.'com.docker.compose.project.config_files' -ne $composeFileFullPath) {
                    throw "Refusing browser run: $container is not owned by this exact disposable Compose project and file."
                }
            }
            & docker run --rm --network "container:$apiContainer" --volumes-from $dashboardContainer -w /app/apps/dashboard `
                -e WS05_DASHBOARD_URL=http://localhost:3002 `
                -e WS05_LANDING_URL=http://localhost:13101 `
                -e WS05_LANDING_SERVICE_URL=http://localhost:3001 `
                -e WS05_API_URL=http://localhost:18080/api/v1 `
                -e WS05_STOREFRONT_URL=http://phase2a-restaurant-a.localhost:3003 `
                -e WS05_SUPER_ADMIN_URL=http://localhost:3004 `
                -e WS05_MAP_LOCALHOST=127.0.0.1 `
                mcr.microsoft.com/playwright:v1.64.1-noble sh -lc './node_modules/.bin/playwright test --config playwright.phase2a.config.ts'
            if ($LASTEXITCODE -ne 0) { throw 'WS05 Playwright suite failed.' }
        }
        'Down' {
            Assert-DisposableProject
            Invoke-Compose @('down', '--volumes', '--remove-orphans')
        }
    }
}
finally {
    Pop-Location
}
