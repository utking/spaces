;(() => {
const submitForm = (tags) => {
    resetError();
    const title = document.getElementById('title').value.trim();
    const url = document.getElementById('url').value.trim();
    const id = document.getElementById('bookmark-id').value.trim();

    if (!id) {
        showError('Bookmark ID is missing. Please refresh the page and try again.');
        return;
    }

    if (!title || !url) {
        showError('Title and URL are required.');
        return;
    }

    if (!/^https?:\/\/.+/i.test(url)) {
        showError('Please enter a valid URL.');
        return;
    }

    if (!tags) {
        showError('Please add at least one tag.');
        return;
    }

    fetch(`/bookmark/${id}/edit`, {
        method: 'PUT',
        headers: {'content-type': 'application/json'},
        body: JSON.stringify({title, url, tags})
    }).then(response => {
        if (response.ok) {
            window.location.href = '/bookmarks';
        } else {
            // if response code 401, show the correct error
            if (response.status === 401) {
                showError('Your session has expired. Please log in again.');
                return;
            }
            return response.json().then((data) => {
                showError(data.Error || 'An error occurred while saving the bookmark.');
            });
        }
    }).catch((error) => {
        showError(error.message);
        console.error('Error:', error);
    });

};

document.addEventListener('DOMContentLoaded', () => {
    const tagSelector = new Tagify(document.getElementById('tags'), {
        enforceWhitelist: false,
        delimiters: ",| ",
        pattern: /^[-!\[\]\(\)/\.=+_a-zA-Z0-9]{1,32}$/,
        required: true,
    });

    document.getElementById('btn-save').addEventListener('click', (e) => {
        if (!tagSelector.value.length) {
            e.preventDefault();
            showError('Please add at least one tag.');
            return;
        }

        submitForm(tagSelector.value.map(tag => tag.value));
    });
});
})();