;(() => {
document.addEventListener("DOMContentLoaded", () => {
    // set up the editor
    const simplemde = new EasyMDE({
        element: document.getElementById('editor-container'),
        placeholder: "Write your note here...",
        renderingConfig: {
            codeSyntaxHighlighting: true,
            hljs: hljs,
        },
        maxHeight: '45vh',
        showIcons: ["horizontal-rule", "strikethrough", "code", "table", "undo", "redo", "side-by-side", "fullscreen"],
        tabSize: 4,
        indentWithTabs: false,
    });
    
    const tagSelector = new Tagify(document.getElementById('tags'), {
        enforceWhitelist: false,
        delimiters: ",| ",
        pattern: /^[-!\[\]\(\)/\.=+_a-zA-Z0-9]{1,32}$/,
        required: true,
    });

    // set up the note create form handler
    document.getElementById('add-item').addEventListener('click', (event) => {
        addItem(tagSelector.value.map(tag => tag.value), simplemde);
    });

    const anchorContent = document.getElementById('anchor-content');
    const needScrolling = $(window).width() < 992;
    if (needScrolling && anchorContent) {
        anchorContent.scrollIntoView({ behavior: 'smooth' });
    }
});

const addItem = (tags, editorEl) => {
    // use Fetch API to send a POST request to the /note/create endpoint
    const noteTitleEl = document.getElementById('note-title');
    const content = editorEl.value().trim();
    const title = noteTitleEl ? noteTitleEl.value.trim() : '';
    
    // reset error block
    resetError();

    if (!title) {
        showError('Note title cannot be empty.');
        return;
    }

    if (!tags) {
        showError('Please add at least one tag.');
        return;
    }

    // prepare and send the request
    fetch('/note/create', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ title, tags, content }),
    }).then(response => {
        if (response.ok) {
            // handle success
            response.json().then((data) => {
                // redirect to the note view page /notes?tag=&note_id=
                window.location = `/notes?tag=${data.Tag}&note_id=${data.ID}`;
            });
        } else {
            // if response code 401, show the correct error
            if (response.status === 401) {
                showError('Your session has expired. Please log in again.');
                return;
            }
            response.json().then((data) => {
                showError(data.Error || 'An error occurred while creating the note.');
                console.error('Error:', data.Error);
            });
        }
    });
}
})();