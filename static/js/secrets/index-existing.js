;(() => {
const updateSecret = (secret_id, tags) => {
    // use Fetch API to send a PUT request to the /secrets endpoint
    const secretNameEl = document.querySelector('#update-secret-form input[name="name"]');
    const name = document.querySelector('#update-secret-form input[name="name"]').value.trim();
    const username_value = document.querySelector('#update-secret-form input[name="username_value"]').value.trim();
    const secret_value = document.querySelector('#update-secret-form input[name="secret_value"]').value.trim();
    const url = document.querySelector('#update-secret-form input[name="url"]').value.trim();
    const description = document.querySelector('#update-secret-form textarea[name="description"]').value.trim();
    
    // reset error block
    resetError();

    if (!name) {
        showError('Secret name cannot be empty.');
        return;
    }

    if (!tags || tags.length === 0) {
        showError('Tags cannot be empty.');
        return;
    }

    if (!secret_id) {
        showError('Secret ID is missing. Please refresh the page and try again.');
        return;
    }

    // prepare and send the request
    fetch(`/secrets`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            secret_id,
            name,
            tags,
            username_value,
            secret_value,
            url,
            description,
        }),
    }).then((response) => {
        if (response.ok) {
            // handle success
            response.json().then((data) => {
                // reload the page with the updated secret
                document.location.reload();
            });
        } else {
            // if response code 401, show the correct error
            if (response.status === 401) {
                showError('Your session has expired. Please log in again.');
                return;
            }
            response.json().then((data) => {
                showError(data.Error || 'An error occurred while updating the secret.');
            });
        }
    }).catch((error) => {
        showError(error.message);
        console.error('Error:', error);
    });
    
}

const deleteItem = (secret_id) => {
    // use Fetch API to send a DELETE request to the /secrets endpoint
    fetch(`/secret/${secret_id}`, {method: 'DELETE'}).then(response => {
        if (response.ok) {
            const url = new URL(window.location.href);
            const tag = url.searchParams.get('tag') || '';
            window.location.href = `/secrets?tag=${tag}`;
        } else {
            // if response code 401, show the correct error
            if (response.status === 401) {
                showError('Your session has expired. Please log in again.');
                return;
            }
            response.json().then((data) => {
                showError(data.Error || 'An error occurred while deleting the secret.');
            });
        }
    }).catch(error => {
        showError(error.message);
        console.error('Error:', error);
    });
}

document.addEventListener("DOMContentLoaded", () => {
    const tagSelector = new Tagify(document.getElementById('tags'), {
        enforceWhitelist: false,
        delimiters: ",| ",
        pattern: /^[-!\[\]\(\)/\.=+_a-zA-Z0-9]{1,32}$/,
        required: true,
    });

    // set up the delete buttons
    document.querySelectorAll('.btn-del-list-secret').forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();
            const itemID = event.currentTarget.getAttribute('data-id');
            const name = event.currentTarget.getAttribute('data-name');
            bootbox.confirm(`Are you sure you want to delete this secret [${name}]?`, (confirmed) => {
                if (confirmed) {
                    deleteItem(itemID);
                }
            });
        });
    });

    // set up the update secret form handler
    if (document.querySelector('#update-secret-form #btn-update')) {
        document.querySelector('#update-secret-form #btn-update').
        addEventListener('click', (event) => {
            event.preventDefault();
            const secret_id = document.getElementById('secret-id').value;
            updateSecret(secret_id, tagSelector.value.map(tag => tag.value));
        });
    }
});
})();