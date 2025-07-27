;(() => {
document.addEventListener("DOMContentLoaded", () => {
    $('.js-search-secrets').select2({
        placeholder: 'Ctrl + / - search',
        allowClear: true,
        maximumInputLength: 20,
        minimumInputLength: 2,
        ajax: {
            url: '/search/secrets',
            delay: 170, // delay in ms before sending the request
            dataType: 'json',
            processResults: (data) => {
                // Transforms the top-level key of the response object from 'items' to 'results'
                return {
                    results: data.items.map((item) => {
                        return {
                            id: `secret_id=${item.id}&tag=${item.tag}`,
                            text: item.text,
                        };
                    })
                };
            }
        },
    });

    // open a new tab with the selected tag (URL)
    $('.js-search-secrets').on('select2:select', (e) => {
        const secretQueryParams = e.params.data.id;
        if (secretQueryParams) {
            document.location.href = `/secrets?${secretQueryParams}`;
        }
    });

    // focus the search box on '/' key press
    document.addEventListener('keydown', (event) => {
        // check if the key is 'Ctrl+/'
        if (event.key === '/' && (event.ctrlKey || event.metaKey)) {
            event.preventDefault(); // Prevent default browser search
            const searchBox = document.querySelector('.js-search-secrets');
            if (searchBox) {
                searchBox.focus();
                searchBox.select(); // Select the text in the input
            }
        }
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

    // scroll-to handlers
    const anchorItems = document.getElementById('anchor-items');
    const anchorContent = document.getElementById('anchor-content');
    const needScrolling = $(window).width() < 992;
    if (needScrolling) {
        if (anchorContent) {
            anchorContent.scrollIntoView({ behavior: 'smooth' });
        } else if (anchorItems) {
            anchorItems.scrollIntoView({ behavior: 'smooth' });
        }
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