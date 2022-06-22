// ref https://github.com/vuejs/core/blob/main/scripts/release.js
/* eslint-disable no-undef */
const args = require('minimist')(process.argv.slice(2))
const fs = require('fs')
const path = require('path')
const chalk = require('chalk')
const semver = require('semver')
const currentVersion = require('../package.json').version
const { prompt } = require('enquirer')
const execa = require('execa')

const preId = args.preid || (semver.prerelease(currentVersion) && semver.prerelease(currentVersion)[0])
const isDryRun = args.dry
const skipBuild = args.skipBuild
const packages = fs.readdirSync(path.resolve(__dirname, '../packages')).filter((p) => !p.endsWith('.ts') && !p.startsWith('.'))

const skippedPackages = []

const versionIncrements = ['patch', 'minor', 'major', ...(preId ? ['prepatch', 'preminor', 'premajor', 'prerelease'] : [])]

const inc = (i) => semver.inc(currentVersion, i, preId)
const run = (bin, args, opts = {}) => execa(bin, args, { stdio: 'inherit', ...opts })
const dryRun = (bin, args, opts = {}) => console.log(chalk.blue(`[dryrun] ${bin} ${args.join(' ')}`), opts)
const runIfNotDry = isDryRun ? dryRun : run
const getPkgRoot = (pkg) => path.resolve(__dirname, '../packages/' + pkg)
const step = (msg) => console.log(chalk.cyan(msg))

async function main() {
    let targetVersion = args._[0]

    if (!targetVersion) {
        // no explicit version, offer suggestions
        const { release } = await prompt({
            type: 'select',
            name: 'release',
            message: 'Select release type',
            choices: versionIncrements.map((i) => `${i} (${inc(i)})`).concat(['custom']),
        })

        if (release === 'custom') {
            targetVersion = (
                await prompt({
                    type: 'input',
                    name: 'version',
                    message: 'Input custom version',
                    initial: currentVersion,
                })
            ).version
        } else {
            targetVersion = release.match(/\((.*)\)/)[1]
        }
    }

    if (!semver.valid(targetVersion)) {
        throw new Error(`invalid target version: ${targetVersion}`)
    }

    const { yes } = await prompt({
        type: 'confirm',
        name: 'yes',
        message: `Releasing v${targetVersion}. Confirm?`,
    })

    if (!yes) {
        return
    }

    // update all package versions and inter-dependencies
    step('\nUpdating cross dependencies...')
    updateVersions(targetVersion)

    // build all packages with types
    step('\nBuilding all packages...')
    if (!skipBuild) {
        await run('pnpm', ['run', 'build'])
    } else {
        console.log(`(skipped)`)
    }

    // update pnpm-lock.yaml
    step('\nUpdating lockfile...')
    await run(`pnpm`, ['install', '--prefer-offline'])

    // publish packages
    step('\nPublishing packages...')
    for (const pkg of packages) {
        await publishPackage(pkg, targetVersion, runIfNotDry)
    }
}

function updateVersions(version) {
    // update website package.json
    updatePackage(path.resolve(__dirname, '../website'), version)
    // update all packages
    packages.forEach((p) => updatePackage(getPkgRoot(p), version))
}

function updatePackage(pkgRoot, version) {
    const pkgPath = path.resolve(pkgRoot, 'package.json')
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'))
    pkg.version = version
    updateDeps(pkg, 'dependencies', version)
    updateDeps(pkg, 'peerDependencies', version)
    fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2) + '\n')
}

function updateDeps(pkg, depType, version) {
    const deps = pkg[depType]
    if (!deps) return
    Object.keys(deps).forEach((dep) => {
        if (dep.startsWith('@natosoft')) {
            console.log(chalk.yellow(`${pkg.name} -> ${depType} -> ${dep}@${version}`))
            deps[dep] = version
        }
    })
}

async function publishPackage(pkgName, version, runIfNotDry) {
    if (skippedPackages.includes(pkgName)) {
        return
    }
    const pkgRoot = getPkgRoot(pkgName)
    const pkgPath = path.resolve(pkgRoot, 'package.json')
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'))
    if (pkg.private) {
        return
    }

    let releaseTag = null
    if (args.tag) {
        releaseTag = args.tag
    } else if (version.includes('alpha')) {
        releaseTag = 'alpha'
    } else if (version.includes('beta')) {
        releaseTag = 'beta'
    } else if (version.includes('rc')) {
        releaseTag = 'rc'
    }

    step(`Publishing ${pkgName}...`)
    try {
        await runIfNotDry('pnpm', ['publish', ...(releaseTag ? ['--tag', releaseTag] : []), '--access', 'public'], {
            cwd: pkgRoot,
            stdio: 'pipe',
        })
        console.log(chalk.green(`Successfully published ${pkgName}@${version}`))
    } catch (e) {
        if (e.stderr.match(/previously published/)) {
            console.log(chalk.red(`Skipping already published: ${pkgName}`))
        } else {
            throw e
        }
    }
}

main().catch((err) => {
    // 回滚
    updateVersions(currentVersion)
    console.error(err)
})
