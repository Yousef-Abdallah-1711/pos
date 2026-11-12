$ErrorActionPreference = 'Stop'

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$containerName = 'pos-postgres'
$expectedProject = 'pos-dev'
$expectedDatabase = 'pos_saas'

$containerInfo = docker inspect --format '{{ index .Config.Labels "com.docker.compose.project" }}|{{ index .Config.Labels "com.docker.compose.project.working_dir" }}|{{ index .Config.Labels "com.docker.compose.service" }}' $containerName
if ($LASTEXITCODE -ne 0) {
    throw "Cannot inspect $containerName. Start the local POS Compose stack first."
}

$labels = $containerInfo -split '\|', 3
$containerRoot = if ($labels.Count -ge 2) { [IO.Path]::GetFullPath($labels[1]) } else { '' }
if ($labels.Count -ne 3 -or $labels[0] -ne $expectedProject -or $labels[2] -ne 'postgres' -or $containerRoot -ne $repoRoot) {
    throw "Refusing to seed: $containerName is not the $expectedProject PostgreSQL container for $repoRoot."
}

$databaseName = docker exec $containerName psql -U postgres -d $expectedDatabase -Atc 'SELECT current_database()'
if ($LASTEXITCODE -ne 0 -or $databaseName.Trim() -ne $expectedDatabase) {
    throw "Refusing to seed: expected local database '$expectedDatabase'."
}

$seedPath = Join-Path $PSScriptRoot 'seed-local-test-accounts.sql'
$containerSeedPath = '/tmp/pos-local-test-accounts.sql'
$passwordEnvironmentNames = @(
    'POS_LOCAL_SUPERADMIN_PASSWORD',
    'POS_LOCAL_OWNER_PASSWORD',
    'POS_LOCAL_CASHIER_PASSWORD'
)
$accountEmails = @(
    'superadmin@pos-saas.com',
    'owner@pos-demo.local',
    'cashier@pos-demo.local'
)
$previousEnvironment = @{}
$seedPasswords = @{}
$seedCopied = $false

try {
    for ($index = 0; $index -lt $passwordEnvironmentNames.Count; $index++) {
        $securePassword = Read-Host -AsSecureString -Prompt "Choose a local password for $($accountEmails[$index])"
        $passwordPointer = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($securePassword)
        try {
            $password = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($passwordPointer)
        } finally {
            [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($passwordPointer)
            $securePassword.Dispose()
        }
        if ($password.Length -lt 16) {
            throw 'Each local test password must be at least 16 characters.'
        }
        $seedPasswords[$passwordEnvironmentNames[$index]] = $password
        $password = $null
    }
    if (@($seedPasswords.Values | Select-Object -Unique).Count -ne $passwordEnvironmentNames.Count) {
        throw 'Choose a different password for each local account.'
    }

    foreach ($environmentName in $passwordEnvironmentNames) {
        $previousEnvironment[$environmentName] = [Environment]::GetEnvironmentVariable($environmentName, 'Process')
        [Environment]::SetEnvironmentVariable($environmentName, $seedPasswords[$environmentName], 'Process')
    }

    docker cp $seedPath "${containerName}:$containerSeedPath"
    if ($LASTEXITCODE -ne 0) {
        throw 'Could not copy the local account seed into the verified development database container.'
    }
    $seedCopied = $true

    docker exec -e POS_LOCAL_SUPERADMIN_PASSWORD -e POS_LOCAL_OWNER_PASSWORD -e POS_LOCAL_CASHIER_PASSWORD $containerName psql -q -v ON_ERROR_STOP=1 -U postgres -d $expectedDatabase -f $containerSeedPath
    if ($LASTEXITCODE -ne 0) {
        throw 'Local account seeding failed; inspect the database error above.'
    }
} finally {
    if ($seedCopied) {
        docker exec $containerName rm -f $containerSeedPath | Out-Null
    }
    foreach ($environmentName in $passwordEnvironmentNames) {
        if ($previousEnvironment.ContainsKey($environmentName)) {
            [Environment]::SetEnvironmentVariable($environmentName, $previousEnvironment[$environmentName], 'Process')
        }
    }
    $password = $null
    $seedPasswords.Clear()
}

Write-Output 'Local development accounts seeded successfully.'
