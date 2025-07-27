;(() => {
document.addEventListener("DOMContentLoaded", () => {
    $('.js-search-notes').select2({
        placeholder: 'Ctrl + / - search',
        allowClear: true,
        maximumInputLength: 20,
        minimumInputLength: 2,
        ajax: {
            url: '/search/notes',
            delay: 170, // delay in ms before sending the request
            dataType: 'json',
            processResults: (data) => {
                // Transforms the top-level key of the response object from 'items' to 'results'
                return {
                    results: data.items.map((item) => {
                        const tag = item.tags ? item.tags[0] : '';
                        return {
                            id: `note_id=${item.id}&tag=${tag}`,
                            text: item.title,
                        };
                    })
                };
            }
        },
    });

    // open a new tab with the selected tag (URL)
    $('.js-search-notes').on('select2:select', (e) => {
        const noteQueryParams = e.params.data.id;
        if (noteQueryParams) {
            document.location.href = `/notes?${noteQueryParams}`;
        }
    });

    // focus the search box on '/' key press
    // ignore if pressed in input or textarea
    document.addEventListener('keydown', (event) => {
        // check if the key is 'Ctrl+/'
        if (event.key === '/' && (event.ctrlKey || event.metaKey)) {
            event.preventDefault(); // Prevent default browser search
            const searchBox = document.querySelector('.js-search-notes');
            if (searchBox) {
                searchBox.focus();
                searchBox.select(); // Select the text in the input
            }
        }
    });

    // set up the export notes button
    document.querySelector('.open-export-page').addEventListener('click', (e) => {
        e.preventDefault();
        const url = e.currentTarget.getAttribute('href');
        bootbox.confirm('Proceed with exporting notes?', (confirmed) => {
            if (confirmed) {
                const a = document.createElement('a');
                a.href = url;
                a.click();
            }
        });
    });

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