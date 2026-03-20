@echo off
setlocal EnableExtensions EnableDelayedExpansion

::==============================================================================
:: vic CLI Setup Script for Windows
::==============================================================================
:: Installs vic CLI tool to %USERPROFILE%\AppData\Local\vic\
:: Adds to user PATH automatically
:: Sets up pre-commit hooks
::==============================================================================

set "SCRIPT_VERSION=1.0.0"
set "VIC_INSTALL_DIR=%USERPROFILE%\AppData\Local\vic"
set "VIC_BINARY_NAME=vic.exe"
set "PROJECT_ROOT="
set "GO_CMD=go"
set "UNINSTALL_MODE=0"
set "SKIP_PRECOMMIT=0"
set "FORCE_REBUILD=0"

:: Detect if running from project root (scripts/setup.bat)
set "SCRIPT_DIR=%~dp0"
set "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"
set "DETECTED_PROJECT_ROOT=%SCRIPT_DIR%\.."
set "DETECTED_PROJECT_ROOT=%DETECTED_PROJECT_ROOT:\scripts=%"

:: Parse command line arguments
:parse_args_loop
if "%~1"=="" goto :args_done
if /i "%~1"=="--uninstall" (
    set "UNINSTALL_MODE=1"
    shift
    goto :parse_args_loop
)
if /i "%~1"=="--skip-precommit" (
    set "SKIP_PRECOMMIT=1"
    shift
    goto :parse_args_loop
)
if /i "%~1"=="--force" (
    set "FORCE_REBUILD=1"
    shift
    goto :parse_args_loop
)
if /i "%~1"=="--help" (
    goto :show_help
)
shift
goto :parse_args_loop

:args_done
:: If project root not provided, try to detect from script location
if not defined PROJECT_ROOT (
    set "PROJECT_ROOT=%DETECTED_PROJECT_ROOT%"
)

:: Normalize project root path
pushd "%PROJECT_ROOT%" >nul 2>&1
if errorlevel 1 (
    set "PROJECT_ROOT=%CD%"
)
set "PROJECT_ROOT=%CD%"
popd

:: Dispatch to appropriate handler
if "%UNINSTALL_MODE%"=="1" goto :uninstall
goto :install

::==============================================================================
:: Show Help
::==============================================================================

:show_help
echo vic CLI Setup Script v%SCRIPT_VERSION%
echo.
echo Usage: setup.bat [OPTIONS]
echo.
echo Options:
echo   --uninstall      Remove vic CLI and clean up PATH
echo   --skip-precommit Skip pre-commit hook setup
echo   --force         Force rebuild even if binary exists
echo   --help          Show this help message
echo.
echo Environment:
echo   PROJECT_ROOT     Override project root directory
echo.
echo Examples:
echo   setup.bat                    Install vic CLI
echo   setup.bat --skip-precommit  Install without pre-commit hooks
echo   setup.bat --uninstall       Remove vic CLI
echo.
goto :eof

::==============================================================================
:: Install Procedure
::==============================================================================

:print_step
    echo.
    echo ========================================
    echo %~1
    echo ========================================
    goto :eof

:print_info
    echo [INFO] %~1
    goto :eof

:print_success
    echo [OK] %~1
    goto :eof

:print_error
    echo [ERROR] %~1 >&2
    goto :eof

:print_warning
    echo [WARN] %~1
    goto :eof

:install
call :print_step "Installing vic CLI"

:: Step 1: Check Go Installation
call :print_info "Checking Go installation..."
"%GO_CMD%" version >nul 2>&1
if errorlevel 1 (
    call :print_error "Go is not installed or not in PATH"
    echo.
    echo Please install Go from: https://go.dev/dl/
    echo.
    echo Installation steps:
    echo   1. Download the Windows installer (.msi)
    echo   2. Run the installer and follow the prompts
    echo   3. Restart your terminal/command prompt
    echo   4. Verify: go version
    echo.
    exit /b 1
)

for /f "tokens=*" %%i in ('"%GO_CMD%" version') do set "GO_VERSION=%%i"
call :print_success "Found: !GO_VERSION!"

:: Step 2: Validate Project Structure
call :print_info "Validating project structure at: %PROJECT_ROOT%"

if not exist "%PROJECT_ROOT%\cmd\vic-go\go.mod" (
    call :print_error "go.mod not found in cmd\vic-go\"
    call :print_info "Make sure you're running from the project root"
    exit /b 1
)

call :print_success "Project structure validated"

:: Step 3: Create Installation Directory
call :print_info "Creating installation directory..."
if not exist "%VIC_INSTALL_DIR%" (
    mkdir "%VIC_INSTALL_DIR%" 2>nul
    if errorlevel 1 (
        call :print_error "Failed to create directory: %VIC_INSTALL_DIR%"
        exit /b 1
    )
)
call :print_success "Installation directory: %VIC_INSTALL_DIR%"

:: Step 4: Build Binary
call :print_info "Building vic.exe..."

set "BUILD_OUTPUT=%VIC_INSTALL_DIR%\%VIC_BINARY_NAME%"
set "SOURCE_DIR=%PROJECT_ROOT%\cmd\vic-go"

:: Force rebuild if requested or binary doesn't exist
if "%FORCE_REBUILD%"=="1" goto :do_build
if not exist "%BUILD_OUTPUT%" goto :do_build

call :print_success "Binary already up-to-date: %BUILD_OUTPUT%"
goto :skip_build

:do_build
pushd "%SOURCE_DIR%"
"%GO_CMD%" build -o "%BUILD_OUTPUT%" .
set "BUILD_RESULT=!errorlevel!"
popd

if !BUILD_RESULT! neq 0 (
    call :print_error "Build failed with exit code: !BUILD_RESULT!"
    exit /b 1
)

call :print_success "Built successfully: %BUILD_OUTPUT%"

:skip_build

:: Step 5: Add to User PATH
call :print_info "Configuring PATH..."

:: Use PowerShell for reliable PATH manipulation
powershell -NoProfile -NonInteractive -Command "$UserPath = [Environment]::GetEnvironmentVariable('Path', 'User'); $InstallDir = '%VIC_INSTALL_DIR%'; if ($UserPath -notlike \"*$InstallDir*\") { [Environment]::SetEnvironmentVariable('Path', \"$UserPath;$InstallDir\", 'User'); Write-Host 'PATH updated with: $InstallDir'; exit 0 } else { Write-Host 'Already in PATH'; exit 0 }"

if errorlevel 1 (
    call :print_warning "Failed to update PATH. You may need to add it manually."
    call :print_info "Add this to your PATH: %VIC_INSTALL_DIR%"
) else (
    call :print_success "PATH configuration updated"
)

:: Step 6: Verify Installation
call :print_info "Verifying installation..."
"%BUILD_OUTPUT%" version >nul 2>&1
if errorlevel 1 (
    call :print_warning "Installation completed but verification failed"
    call :print_info "Try restarting your terminal or running: %BUILD_OUTPUT% version"
) else (
    call :print_success "Verification passed"
)

:: Step 7: Setup Pre-commit Hook (if in git repo)
if "%SKIP_PRECOMMIT%"=="1" (
    call :print_info "Skipping pre-commit setup (--skip-precommit)"
) else (
    call :setup_precommit
)

:: Step 8: Show Success Message
call :print_step "Installation Complete!"

echo.
echo vic CLI has been installed successfully!
echo.
echo Installation Location: %VIC_INSTALL_DIR%\%VIC_BINARY_NAME%
echo.
echo IMPORTANT:
echo   1. Restart your terminal/command prompt for PATH changes to take effect
echo   2. Or run the following command to use vic immediately:
echo      "%VIC_INSTALL_DIR%\%VIC_BINARY_NAME%" version
echo.
echo Next Steps:
echo   - Run 'vic init' to initialize a project
echo   - Run 'vic --help' to see all commands
echo   - Run 'vic spec gate 0' to check requirements
echo.
echo To uninstall, run: setup.bat --uninstall
echo.

exit /b 0

::==============================================================================
:: Pre-commit Setup
::==============================================================================

:setup_precommit
call :print_info "Setting up pre-commit hooks..."

:: Check if in git repo
if not exist "%PROJECT_ROOT%\.git" (
    call :print_info "Not a git repository, skipping pre-commit setup"
    goto :eof
)

:: Check for pre-commit CLI
where pre-commit >nul 2>&1
if errorlevel 1 (
    call :print_warning "pre-commit not found. Skipping hook setup."
    call :print_info "To install pre-commit, run: pip install pre-commit"
    call :print_info "Or: winget install pre-commit.pre-commit"
    goto :eof
)

:: Install hooks
pushd "%PROJECT_ROOT%" >nul 2>&1
if not errorlevel 1 (
    pre-commit install --hook-type pre-commit 2>nul
    pre-commit install --hook-type commit-msg 2>nul
    set "PRECOMMIT_RESULT=!errorlevel!"
    popd

    if "!PRECOMMIT_RESULT!"=="0" (
        call :print_success "Pre-commit hooks installed"
    ) else (
        call :print_warning "Failed to install pre-commit hooks"
    )
)
goto :eof

::==============================================================================
:: Uninstall Procedure
::==============================================================================

:uninstall
call :print_step "Uninstalling vic CLI"

:: Step 1: Remove from PATH
call :print_info "Removing from PATH..."

powershell -NoProfile -NonInteractive -Command "$UserPath = [Environment]::GetEnvironmentVariable('Path', 'User'); $InstallDir = '%VIC_INSTALL_DIR%'; $Paths = $UserPath -split ';' | Where-Object { $_ -and $_ -ne $InstallDir }; [Environment]::SetEnvironmentVariable('Path', ($Paths -join ';'), 'User'); Write-Host 'PATH updated (vic removed)'; exit 0"

if errorlevel 1 (
    call :print_warning "Failed to update PATH"
) else (
    call :print_success "Removed from PATH"
)

:: Step 2: Remove binary
if exist "%VIC_INSTALL_DIR%\%VIC_BINARY_NAME%" (
    del /f /q "%VIC_INSTALL_DIR%\%VIC_BINARY_NAME%"
    if errorlevel 1 (
        call :print_warning "Failed to delete binary"
    ) else (
        call :print_success "Binary removed"
    )
)

:: Step 3: Remove installation directory (if empty)
if exist "%VIC_INSTALL_DIR%" (
    rmdir "%VIC_INSTALL_DIR%" 2>nul
    if errorlevel 1 (
        call :print_warning "Installation directory not empty, leaving it"
    ) else (
        call :print_success "Installation directory removed"
    )
)

:: Step 4: Uninstall pre-commit hooks
if exist "%PROJECT_ROOT%\.git" (
    pushd "%PROJECT_ROOT%" >nul 2>&1
    if not errorlevel 1 (
        pre-commit uninstall 2>nul
        popd
        call :print_success "Pre-commit hooks removed"
    )
)

call :print_step "Uninstall Complete!"

echo.
echo vic CLI has been uninstalled.
echo.
echo NOTE: You may need to restart your terminal for PATH changes to take effect.
echo.

exit /b 0
