# Changelog

All notable changes to this project will be documented in this file.

## [v1.0.0] - 2025-11-28

### Added
- **Dark/Light Mode**: Added a theme toggle button with sun/moon animation. Theme preference is saved in `config.json`.
- **Clear Buttons**: Added "Clear" buttons to quickly reset input and output fields.
- **Background Theming**: Application background now adapts to the selected theme.

### Fixed
- **Config Loading**: Fixed an issue where `config.json` was not found when running with `go run .`.
- **Clipboard**: Fixed a regression where the copy to clipboard functionality was missing.
- **Regex**: Updated default hostname pattern to better match common hostname formats.

### Changed
- **UI**: Improved layout and spacing for better usability.
- **Documentation**: Updated README with new features and screenshots.
