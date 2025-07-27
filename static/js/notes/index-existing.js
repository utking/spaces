;(() => {
const updateNote = (note_id, tags, editorEl) => {
    const title = document.querySelector('#note-update-form input[name="title"]').value.trim();
    const content = editorEl.value().trim();
    
    // reset error block
    resetError();

    if (!tags) {
        showError('at least one tag is required');
        return;
    }

    if (!title) {
        showError('note title cannot be empty');
        return;
    }

    if (!note_id) {
        showError('note id is not specified');
        return;
    }

    // prepare and send the request
    fetch(`/notes`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            title,
            tags,
            content,
            note_id,
        }),
    }).then(response => {
        if (response.ok) {
            // handle success
            response.json().then((data) => {
                // reload the page
                window.location.reload();
            });
        } else {
            // if response code 401, show the correct error
            if (response.status === 401) {
                showError('Your session has expired. Please log in again.');
                return;
            }
            response.json().then((data) => {
                // if response code 401, show the correct error
                if (response.status === 401) {
                    showError('Your session has expired. Please log in again.');
                    return;
                }
                showError(data.Error || 'An error occurred while updating the note.');
            });
        }
    }).catch((error) => {
        showError(error.message || 'An error occurred while updating the note.');
        console.error('Error:', error);
    });
}

const deleteNote = (note_id) => {
    bootbox.confirm('Are you sure you want to delete this note?', (confirmed) => {
    if (confirmed) {
        fetch(`/note/${note_id}`, {method: 'DELETE'}).then(response => {
            if (response.ok) {
                // redirect to the notes page without the deleted note
                // tag is from the query parameters
                const url = new URL(window.location.href);
                const tag = url.searchParams.get('tag') || '';
                window.location.href = `/notes?tag=${tag}`;
            } else {
                // if response code 401, show the correct error
                if (response.status === 401) {
                    showError('Your session has expired. Please log in again.');
                    return;
                }
                response.json().then((data) => {
                    showError(data.Error || 'An error occurred while deleting the note.');
                    console.error('Error:', data.Error);
                });
            }
        });
    }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    const noteId = document.getElementById('note-id').value;

    // set up the editor
    const simplemde = new EasyMDE({
        element: document.getElementById(`editor-container-${noteId}`),
        renderingConfig: {
            codeSyntaxHighlighting: true,
            hljs: hljs,
        },
        maxHeight: '53vh',
        showIcons: ["horizontal-rule", "strikethrough", "code", "table", "undo", "redo", "side-by-side", "fullscreen"],
        tabSize: 4,
        indentWithTabs: false,
    });
    simplemde.togglePreview();

    const tagSelector = new Tagify(document.getElementById('tags'), {
        enforceWhitelist: false,
        delimiters: ",| ",
        pattern: /^[-!\[\]\(\)/\.=+_a-zA-Z0-9]{1,32}$/,
        required: true,
    })

    // set up the note update form handler
    if (document.querySelector('#note-update-form > #btn-update')) {
        document.querySelector('#note-update-form > #btn-update').addEventListener('click', (event) => {
            event.preventDefault();
            updateNote(noteId, tagSelector.value.map(tag => tag.value), simplemde);
        });
    }

    // set up the delete button for notes in the list
    const deleteButtons = document.querySelectorAll('.btn-del-list-note');
    deleteButtons.forEach(button => {
        button.addEventListener('click', (event) => {
            event.preventDefault();
            const noteId = event.currentTarget.getAttribute('data-id');
            if (noteId) {
                deleteNote(noteId);
            } else {
                showError('note id is not specified');
            }
        });
    });
});
})();