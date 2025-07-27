;(() => {
    $('.dropdown-menu a.dropdown-toggle').on('click', (e) => {
        if (!$(this).next().hasClass('show')) {
            $(this).parents('.dropdown-menu').first().find('.show').removeClass('show');
        }
        const $subMenu = $(this).next('.dropdown-menu');
        $subMenu.toggleClass('show');

        $(this).parents('li.nav-item.dropdown.show').on('hidden.bs.dropdown', (e) => {
            $('.dropdown-submenu .show').removeClass('show');
        });

        return false;
    });
    document.addEventListener('DOMContentLoaded', () => {
        const errorBlock = document.getElementById('error-block');
        if (errorBlock) {
            const btnClose = errorBlock.querySelector('.btn-error-hide');
            if (btnClose) {
                btnClose.addEventListener('click', () => {
                    errorBlock.style.display = 'none';
                });
            }
        }
    });

    window.showError = (message) => {
        const errorBlock = document.getElementById('error-block');
        const errorBlockContent = document.querySelector('#error-block .error-content');
        if (errorBlock && errorBlockContent) {
            let exclIcon = document.createElement('i');
            exclIcon.className = 'bi bi-exclamation-triangle';
            errorBlockContent.innerHTML = ''; // Clear previous content
            errorBlockContent.appendChild(exclIcon);
            errorBlockContent.appendChild(document.createTextNode(' ' + message));
            // Show the error block
            errorBlock.style.display = 'block';
            // scroll to the top of the document
            window.scrollTo({
                top: 0,
                behavior: 'smooth'
            });
        }
    };

    window.resetError = () => {
        const errorBlock = document.getElementById('error-block');
        const errorBlockContent = document.querySelector('#error-block .error-content');
        if (errorBlock) {
            errorBlockContent.textContent = '';
            errorBlock.style.display = 'none';
        }
    };
})();