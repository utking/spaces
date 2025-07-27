;(() => {
document.addEventListener("DOMContentLoaded", () => {
    const tagSelector = new Tagify(document.getElementById('tags'), {
        enforceWhitelist: false,
        delimiters: ",| ",
        pattern: /^[-!\[\]\(\)/\.=+_a-zA-Z0-9]{1,32}$/,
        required: true,
    });

    // set up the note create form handler
    document.getElementById('add-item').addEventListener('click', (event) => {
        addItem(tagSelector.value.map(tag => tag.value));
    });

    // set up the show password button. show password while pressed
    const btnShowPassword = document.getElementById('btn-show-password');
    const passwordInput = document.getElementById('password');
    if (btnShowPassword && passwordInput) {
        btnShowPassword.addEventListener('mousedown', () => {
            passwordInput.type = 'text';
        });
        btnShowPassword.addEventListener('mouseup', () => {
            passwordInput.type = 'password';
        });
        btnShowPassword.addEventListener('mouseleave', () => {
            passwordInput.type = 'password';
        });

        // also, show password on tap on mobile devices
        btnShowPassword.addEventListener('touchstart', () => {
            passwordInput.type = 'text';
        });
        btnShowPassword.addEventListener('touchend', () => {
            passwordInput.type = 'password';
        });
        btnShowPassword.addEventListener('touchcancel', () => {
            passwordInput.type = 'password';
        });
        btnShowPassword.addEventListener('touchmove', () => {
            passwordInput.type = 'password';
        });
    }

    // set up copy secret value to clipboard on click
    const btnCopyPassword = document.getElementById('btn-copy-password');
    if (btnCopyPassword) {
        btnCopyPassword.addEventListener('click', () => {
            const passwordInput = document.getElementById('password');
            navigator.clipboard.writeText(passwordInput.value).then(() => {
                // set text to "Copied!" for 3 second. then set it back to "Copy"
                btnCopyPassword.innerText = 'Copied!';
                setTimeout(() => {
                    btnCopyPassword.innerText = 'Copy';
                }, 3000);
            }).catch(err => {
                btnCopyPassword.innerText = 'Error copying password';
            });
        });
    }

    // set up copy username to clipboard on click
    const btnCopyUsername = document.getElementById('btn-copy-username');
    if (btnCopyUsername) {
        btnCopyUsername.addEventListener('click', () => {
            const usernameInput = document.getElementById('username');
            navigator.clipboard.writeText(usernameInput.value).then(() => {
                // set text to "Copied!" for 3 second. then set it back to "Copy"
                btnCopyUsername.innerText = 'Copied!';
                setTimeout(() => {
                    btnCopyUsername.innerText = 'Copy';
                }, 3000);
            }).catch(err => {
                btnCopyUsername.innerText = 'Error copying username';
            });
        });
    }

    const anchorContent = document.getElementById('anchor-content');
    const needScrolling = $(window).width() < 992;
    if (needScrolling) {
        if (anchorContent) {
            anchorContent.scrollIntoView({ behavior: 'smooth' });
        }
    }
});

const addItem = (tags) => {
    // use Fetch API to send a POST request to the /secret/create endpoint
    const secretNameEl = document.querySelector('#secret-create-form input[name="name"]');
    const secretValueEl = document.querySelector('#secret-create-form input[name="secret_value"]');
    const username_value = document.querySelector('#secret-create-form input[name="username_value"]').value.trim();
    const name = secretNameEl ? secretNameEl.value.trim() : '';
    const url = document.querySelector('#secret-create-form input[name="url"]').value.trim();
    const description = document.querySelector('#secret-create-form textarea[name="description"]').value.trim();
    const secret_value = secretValueEl ? secretValueEl.value.trim() : '';
    
    // reset error block
    resetError();

    if (!name) {
        showError('Secret name cannot be empty.');
        return;
    }

    if (!tags) {
        showError('Please add at least one tag.');
        return;
    }

    // prepare and send the request
    fetch('/secret/create', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            name,
            url,
            username_value,
            description,
            tags,
            secret_value,
        }),
    }).then(response => {
        if (response.ok) {
            // handle success
            response.json().then((data) => {
                // redirect to the note view page /secrets?tag=&note_id=
                window.location = `/secrets?tag=${data.Tag}&secret_id=${data.ID}`;
            });
        } else {
            // if response code 401, show the correct error
            if (response.status === 401) {
                showError('Your session has expired. Please log in again.');
                return;
            }
            response.json().then((data) => {
                showError(data.Error || 'An error occurred while creating the secret.');
            });
        }
    }).catch(error => {
        showError(error.message || 'An error occurred while creating the secret.');
        console.error('Error:', error);
    });
}
})();