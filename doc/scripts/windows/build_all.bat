ECHO OFF
ECHO Building all versions of nuiforms...
CALL build_for_linux_amd64.bat
CALL build_for_linux_arm64.bat
CALL build_for_macos_arm64.bat
CALL build_for_windows_amd64.bat
CALL build_for_windows_arm64.bat
ECHO Finished building all versions of nuiforms.
PAUSE
