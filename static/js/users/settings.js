;(() => {
document.getElementById('saveSettingsButton').addEventListener('click', () => {
    const darkModeCheckbox = document.getElementById('darkModeCheckbox');
    const fileBrowserTilesCheckbox = document.getElementById('fileBrowserTilesCheckbox');
    resetError();
    fetch('/users/settings', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            dark_mode_enabled: darkModeCheckbox.checked,
            file_browser_tiles: fileBrowserTilesCheckbox.checked
        })
    })
    .then(response => response.json())
    .then(data => {
        if (!data.Error) {
            location.reload();
        } else {
            showError("Error setting the dark mode: " + data.Error);
        }
    })
    .catch(console.error);
});
})();