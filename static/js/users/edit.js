;(() => {
document.addEventListener("DOMContentLoaded", () => {
    const idInput = document.getElementById('id');
    const form = document.getElementById('edit-form');
    const submitBtn = document.getElementById('submit-update');
    submitBtn.addEventListener("click", () => {
        const formData = new FormData(form);
        const id = idInput.value;
        if (!id) {
            showError('User ID is missing. Please refresh the page and try again.');
            return false;
        }

        // reset error block
        resetError();

        const url = `/user/${id}/edit`;
        fetch(url, {
            method: 'PUT',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams(formData)
        }).then((response) => {
            if (response.ok) {
                // handle success
                response.json().then((data) => {
                    // reload the page with the updated secret
                    location.href = `/user/${id}`;
                });
            } else {
                // if response code 401, show the correct error
                if (response.status === 401) {
                    showError('Your session has expired. Please log in again.');
                    return;
                }
                response.json().then((data) => {
                    showError(data.Error || 'An error occurred while updating the user.');
                });
            }
        }).catch((error) => {
            showError(error.message);
            console.error('Error:', error);
        });
    });
});
})();