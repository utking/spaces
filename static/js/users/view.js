;(() => {
document.addEventListener("DOMContentLoaded", () => {
    const deleteBtn = document.getElementById('delete-item');
    const idInput = document.getElementById('id');

    if (!deleteBtn) { return; }

    deleteBtn.addEventListener('click', () => {
        bootbox.confirm('Suspend this user?', (confirmed) => {
            if (confirmed) {
                resetError();

                fetch(`/user/${idInput.value}`, {method: 'DELETE'})
                // if status is 204, it means the user was deleted successfully
                .then(response => {
                    if (response.ok) {
                        window.location.reload();
                    } else {
                        // if response code 401, show the correct error
                        if (response.status === 401) {
                            showError('Your session has expired. Please log in again.');
                            return;
                        }
                        return response.json().then(data => {
                            showError(data.Error || 'Failed to delete user');
                            console.error('Error deleting user:', data.Error);
                        });
                    }
                })
                .catch((error) => {
                    showError(error.message || 'An error occurred while deleting the user.');
                    console.error('Error:', error)
                });
            }
        });
    });
});
})();