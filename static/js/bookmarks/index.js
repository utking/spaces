;(() => {
const addItem = (tags) => {
    const title = document.querySelector('#title').value;
    const url = document.querySelector('#url').value;

    // reset error block
    resetError();

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

    const formData = new FormData();
    formData.append('tags', tags ? tags : '');
    const data = {
        title,
        url,
        tags,
    };

    fetch('/bookmark/create', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data)
    })
    .then((response) => {
        response.json().then((data) => {
            if (response.ok) {
                window.location.reload();
            } else {
                // if response code 401, show the correct error
                if (response.status === 401) {
                    showError('Your session has expired. Please log in again.');
                    return;
                }
                showError(data.Error || 'An error occurred while adding the bookmark.');
                console.error('Error adding bookmark:', data.Error);
            }
        });
    })
    .catch((error) => {
        showError(error.message);
        console.error('Error:', error);
    });
};

// document ready
document.addEventListener('DOMContentLoaded', () => {
    $('.js-search-selector').select2({
        placeholder: 'Ctrl + / - search',
        allowClear: true,
        maximumInputLength: 20,
        minimumInputLength: 2,
        ajax: {
            url: '/search/bookmarks',
            delay: 170, // delay in ms before sending the request
            dataType: 'json',
            processResults: (data) => {
                // Transforms the top-level key of the response object from 'items' to 'results'
                return {
                    results: data.items.map((item) => {
                        return {
                            id: item.id,
                            text: item.text,
                        };
                    })
                };
            }
        },
    });

    // open a new tab with the selected tag (URL)
    $('.js-search-selector').on('select2:select', (e) => {
        const tag = e.params.data.id;
        if (tag) {
            window.open(tag, '_blank');
        }
    });

    // focus the search box on '/' key press
    document.addEventListener('keydown', (event) => {
        // check if the keys are 'Ctrl+/'
        if (event.key === '/' && (event.ctrlKey || event.metaKey)) {
            event.preventDefault(); // Prevent default browser search
            const searchBox = document.querySelector('.js-search-selector');
            if (searchBox) {
                searchBox.focus();
                searchBox.select(); // Select the text in the input
            }
        }
    });
    
    const tagSelector = new Tagify(document.getElementById('tags'), {
        enforceWhitelist: false,
        delimiters: ",| ",
        pattern: /^[-!\[\]\(\)/\.=+_a-zA-Z0-9]{1,32}$/,
        required: true,
    });

    document.getElementById('add-item').addEventListener('click', (event) => {
        addItem(tagSelector.value.map(tag => tag.value));
    });

    document.querySelector('.open-export-page').addEventListener('click', (e) => {
        e.preventDefault();
        const url = e.currentTarget.getAttribute('href');
        bootbox.confirm('Proceed with exporting bookmarks?', (confirmed) => {
            if (confirmed) {
                const a = document.createElement('a');
                a.href = url;
                a.click();
            }
        });
    });
    
    document.querySelectorAll('.btn-delete').forEach((button) => {
        button.addEventListener('click', (event) => {
            const id = event.currentTarget.getAttribute('data-id');
            bootbox.confirm('Sure you want to delete this bookmark?', (confirmed) => { 
                if (confirmed) {
                    fetch('/bookmark/' + id, {method: 'DELETE'})
                    // if status is 204, it means the bookmark was deleted successfully
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
                                showError(data.Error || 'Failed to delete bookmark');
                                console.error('Error deleting bookmark:', data.Error);
                            });
                        }
                    })
                    .catch((error) => {
                        showError(error.message || 'An error occurred while deleting the bookmark.');
                        console.error('Error:', error)
                    });
                }
            });
        });
    });

    // Show add bookmark form on button click. Toggle visibility for the button, form and the table
    document.getElementById('show-add-form').addEventListener('click', () => {
        const form = document.getElementById('add-bookmark-form');
        const titleElement = document.getElementById('title');
        const bookmarksList = document.getElementById('bookmarks-list');
        const showAddFormButton = document.getElementById('show-add-form');
        if (form.style.display === 'none' || form.style.display === '') {
            form.style.display = 'block';
            if (titleElement) {
                titleElement.focus(); // Focus on the title input when the form is shown
            }
            if (bookmarksList) {bookmarksList.style.display = 'none';}
            showAddFormButton.style.display = 'none';
        } else {
            form.style.display = 'none';
            if (bookmarksList) {bookmarksList.style.display = 'block';}
            showAddFormButton.style.display = 'inline-block';
        }
    });

    // Cancel button to hide the form and show the bookmarks list
    document.getElementById('cancel-add').addEventListener('click', () => {
        document.getElementById('add-bookmark-form').style.display = 'none';
        document.getElementById('bookmarks-list').style.display = 'block';
        document.getElementById('show-add-form').style.display = 'inline-block';
        // Clear the error block
        resetError();
    });

    // scroll-to handlers
    const anchorItems = document.getElementById('anchor-items');
    const needScrolling = $(window).width() < 992;
    if (anchorItems && needScrolling) {
        anchorItems.scrollIntoView({ behavior: 'smooth' });
    }
    // scroll to the active tag in the tag list
    const tagList = document.getElementById('tag-list');
    const activeTag = document.querySelector('.list-group-item-primary');
    if (tagList && activeTag) {
        tagList.scrollTo({
            top: activeTag.offsetTop - tagList.offsetTop - 10,
            behavior: 'smooth',
        });
    }
});
})();