# Spaces

A simple Personal Space multi-user Web application

## Installation

1. Build it from the source (`make build`) or download a pre-built binary
2. Run migrations (`bin/spaces migrate`)
3. Create `.env` (you can copy it from `.env.example` and adjust as needed)
    * The application works only with MySQL or SQLite to store app data in the database (set `SQL_DRIVER` accordingly)
    * Also, it requires write access to a disk storage (for logs and FileBrowser functionality)
    * SMTP, though the configuration is required, is not essential for the application functionality
    * You can allow self-registrations if you need this
4. Start the application (`bin/spaces serve`)
5. Log in with the default user/password (root/password) and create a new user with a strong password
    * It is recommended to set a strong password for the root user and disable it

## Functionality

### Ready/InProgress

* Note taking
    * [x] notes have tags for better categorization
    * [x] notes are Markdown-based
    * [x] notes' visibility is limited to the user-owner
    * [x] notes import/export as JSON
    * [x] seach notes by content and/or title
* Password storage / Vault
    * [x] passwords have tags for better categorization
    * [x] passwords encryption is per user and having one user's key won't expose other users' secrets
    * [x] passwords' visibility is limited to the user-owner
    * [x] passwords import/export as JSON
    * [x] seach passwords by their username/url/description/name
* Bookmarking
    * [x] bookmarks have tags for better categorization
    * [x] bookmarks import/export as JSON
    * [x] search bookmarks by title/url
* File storage / File browser
    * [x] Tile/list views
    * [x] Type-aware icons for some files
    * [x] Upload files
    * [x] Download files
    * [x] Delete files/folders
    * [x] Rename files/folders
    * [x] Preview some file types

### Markdown editor options

* Currently in use: https://github.com/Ionaru/easy-markdown-editor
