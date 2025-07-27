;(() => {
document.addEventListener("DOMContentLoaded", () => {
    $('.delete-item').click((el) => {
        const id = $(el.currentTarget).data('id');
        const name = $(el.currentTarget).data('name');

        bootbox.confirm(`Are you sure you want to suspend user [${name}]?`, (confirmed) => {
            if (confirmed) {
                fetch(`/user/${id}`, {method: 'DELETE'})
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

    // submit form on status change in the filter
    document.getElementById('status-select').addEventListener('change', () => {
        const form = document.querySelector('form');
        if (form) {
            form.submit();
        }
    });
});
})();