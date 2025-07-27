;((document) => {
    document.addEventListener('DOMContentLoaded', () => {
        const uploadForm = document.getElementById('upload-form');
        const uploadButton = document.getElementById('upload-button');
        const fileInput = document.getElementById('file-input');
        const spinner = document.querySelector('.spinner-border');
        const createFolderButton = document.getElementById('create-folder-button');
        const folderNameInput = document.getElementById('folder-name');

        const deleteItem = (el) => {
            const itemName = el.currentTarget.getAttribute('data-name');
            if (!itemName) {
                showError('No item name found for delete action');
                return;
            }
            // hide context menu and reset it
            document.getElementById('context-menu').classList.remove('show');
            document.getElementById('context-menu').innerHTML = '';

            // show confirmation dialog
            const promptMsg = `Are you sure you want to delete "${itemName}"? This action cannot be undone.`;
            bootbox.confirm(promptMsg, (confirmed) => {
                if (confirmed) {
                    // fetch request to delete the item (POST /filebrowser/delete). 
                    // Body: path, name
                    const formData = new FormData();
                    // get the path from the current URL
                    const currentPath = new URL(window.location.href).searchParams.get('path') || '';
                    formData.append('path', currentPath);
                    formData.append('name', itemName);
                    fetch('/filebrowser/delete', {
                        method: 'DELETE',
                        body: formData,
                    })
                    .then(response => {
                        if (response.ok) {
                            // get the current path from the URL and reload the page and append ..
                            const currentPath = new URL(window.location.href).searchParams.get('path') || '';
                            if (currentPath) {
                                window.location.href = `/filebrowser?path=${currentPath}`;
                            } else {
                                window.location.href = '/filebrowser';
                            }
                        } else {
                            // if response code 401, show the correct error
                            if (response.status === 401) {
                                showError('Your session has expired. Please log in again.');
                                return;
                            }
                            return response.json().then(data => {
                                showError(data.Error || 'Failed to delete item');
                                console.error('Error deleting item:', data.Error || 'Failed to delete item');
                            });
                        }
                    })
                    .catch(error => {
                        showError('An error occurred while deleting the item');
                        console.error('Error:', error);
                    });
                }
            });
        };
        const renameItem = (el) => {
            const itemName = el.currentTarget.getAttribute('data-name');
            if (!itemName) {
                showError('No item name found for rename action');
                return;
            }
            document.getElementById('context-menu').classList.remove('show');
            document.getElementById('context-menu').innerHTML = '';
            bootbox.prompt({
                title: 'Rename item',
                value: itemName,
                callback: (newName) => {
                    if (newName && newName.trim() !== '' && newName !== itemName) {
                        resetError();
                        // fetch request to rename the item (POST /filebrowser/rename). 
                        // Body: path, old_name, new_name
                        const formData = new FormData();
                        // get the path from the current URL
                        const currentPath = new URL(window.location.href).searchParams.get('path') || '';
                        formData.append('path', currentPath);
                        formData.append('old_name', itemName);
                        formData.append('new_name', newName.trim());
                        fetch('/filebrowser/rename', {
                            method: 'POST',
                            body: formData,
                        })
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
                                    showError(data.Error || 'Failed to rename item');
                                    console.error('Error renaming item:', data.Error || 'Failed to rename item');
                                });
                            }
                        })
                        .catch(error => {
                            showError('An error occurred while renaming the item');
                            console.error('Error:', error);
                        });
                    }
                }
            });
        };
        const downloadItem = (el) => {
            const itemName = el.currentTarget.getAttribute('data-name');
            if (!itemName) {
                return;
            }

            document.getElementById('context-menu').classList.remove('show');
            document.getElementById('context-menu').innerHTML = '';

            // fetch request to download the item (GET /filebrowser/download)
            const currentPath = new URL(window.location.href).searchParams.get('path') || '';
            const downloadUrl = `/filebrowser/download?path=${currentPath}&name=${encodeURIComponent(itemName)}`;
            fetch(downloadUrl)
                .then(response => {
                    if (response.ok) {
                        // create a link element to download the file
                        const link = document.createElement('a');
                        link.href = response.url; // use the response URL for download
                        link.download = itemName; // set the download attribute
                        document.body.appendChild(link);
                        link.click();
                        document.body.removeChild(link);
                    } else {
                        // if response code 401, show the correct error
                        if (response.status === 401) {
                            showError('Your session has expired. Please log in again.');
                            return;
                        }
                        return response.json().then(data => {
                            showError(data.Error || 'Failed to download item');
                            console.error('Error downloading item:', data.Error || 'Failed to download item');
                        });
                    }
                })
                .catch(error => {
                    showError('An error occurred while downloading the item');
                    console.error('Error:', error);
                });
        };
        let menu = [
            {
                Title: 'Download',
                Action: downloadItem,
                Icon: '<i class="bi bi-download"></i>'
            },
            {
                Title: 'Rename',
                Action: renameItem,
                Icon: '<i class="bi bi-pencil"></i>'
            },
            {
                Title: 'Delete',
                Action: deleteItem,
                Icon: '<i class="bi bi-trash"></i>'
            }];

        function renderContextMenu(menu, context, itemName, itemType) {
            context.innerHTML ='';
            menu.forEach((item) => {
                // skip download if itemType is folder
                if (itemType === 'folder' && item.Title === 'Download') {
                    return;
                }
                let menuRow = document.createElement('li');
                let menuitem = document.createElement('a');
                let icon = document.createElement('span');
                
                icon.innerHTML = `${item.Icon} `;

                menuitem.innerHTML = item.Title;
                menuitem.addEventListener('click', item.Action);
                menuitem.classList.add('dropdown-item');
                menuitem.setAttribute('data-name', itemName);
                menuitem.href = '#';

                menuitem.prepend(icon);
                menuRow.appendChild(menuitem);
                context.appendChild(menuRow);
            });
        }

        window.addEventListener('load', () => {
            // on mobile: on long press, show context menu
            let timer;
            let touchduration = 500;
            let touchstartX, touchstartY;

            $('.file-link').on('contextmenu', function(e){
                e.preventDefault();

                const link = e.currentTarget;
                const itemName = link.getAttribute('data-name');
                const itemType = link.getAttribute('data-type');
                renderContextMenu(menu, document.getElementById('context-menu'), itemName, itemType);

                // position menu at the right-click cursor
                document.getElementById('context-menu').style.left = e.clientX + 'px';
                document.getElementById('context-menu').style.top = e.clientY + 'px';
                document.getElementById('context-menu').classList.add('show');
            });

            function touchstart(e) {
                // do not prevent default for data-gall links
                if (!e.currentTarget.classList.contains('venobox')) {
                    if (e.cancelable) e.preventDefault();
                }
                if (!timer) {
                    touchstartX = e.changedTouches[0].clientX;
                    touchstartY = e.changedTouches[0].clientY;
                    timer = setTimeout(showContextMenu(e), touchduration);
                }
            }

            function touchend(e) {
                //stops short touches from firing the event
                if (timer) {
                    clearTimeout(timer);
                    timer = null;
                }
                // if the touchend was far from the touchstart, do not show context menu
                const touchendX = e.changedTouches[0].clientX;
                const touchendY = e.changedTouches[0].clientY;
                const distance = Math.sqrt(Math.pow(touchendX - touchstartX, 2) + Math.pow(touchendY - touchstartY, 2));
                if (distance > 10) {
                    clearTimeout(timer);
                    timer = null;
                }
            }

            const showContextMenu = (e) => {
                return () => {
                    e.preventDefault();
                    const link = e.currentTarget;
                    const itemName = link.getAttribute('data-name');
                    const itemType = link.getAttribute('data-type');
                    renderContextMenu(menu, document.getElementById('context-menu'), itemName, itemType);

                    // position menu at the touch cursor
                    document.getElementById('context-menu').style.left = e.changedTouches[0].clientX + 'px';
                    document.getElementById('context-menu').style.top = e.changedTouches[0].clientY + 'px';
                    document.getElementById('context-menu').classList.add('show');
                }
            };

            $('.file-link').on('touchstart', touchstart);
            $('.file-link').on('touchend', touchend);
            $('.file-link').on('touchmove', function(e) {
                // if the touch moves more than 10px, cancel the context menu
                if (timer) {
                    const touchendX = e.changedTouches[0].clientX;
                    const touchendY = e.changedTouches[0].clientY;
                    const distance = Math.sqrt(Math.pow(touchendX - touchstartX, 2) + Math.pow(touchendY - touchstartY, 2));
                    if (distance > 10) {
                        clearTimeout(timer);
                        timer = null;
                    }
                }
            });

            $('.files-list').on('click', function(e){
                e.preventDefault();
                document.getElementById('context-menu').innerHTML = '';
                document.getElementById('context-menu').classList.remove('show');
            });

            // on touch anywhere outside the context menu, hide it
            document.addEventListener('touchstart', function(e) {
                if (!e.target.closest('.context-menu')) {
                    document.getElementById('context-menu').innerHTML = '';
                    document.getElementById('context-menu').classList.remove('show');
                }
            }, { passive: true });
        });
        
        // folder-item must open on click
        const folderItems = document.querySelectorAll('.folder-item');
        folderItems.forEach(item => {
            item.addEventListener('click', (event) => {
                event.preventDefault();
                const path = item.getAttribute('href');
                window.location.href = path;
            });
        });

        const openViewable = (event, link) => {
            return () => {
                if (!link || !link.getAttribute('data-href')) {
                    return;
                }
                event.preventDefault();
                const href = link.getAttribute('data-href');
                if (href) {
                    const modal = document.querySelector('#my-dialog');
                    const modalContent = modal.querySelector('.modal-body');
                    const modalHeader = modal.querySelector('.modal-header');

                    modalContent.innerHTML = `<iframe src="${href}" style="width: 100%; height: 100%; border: none;"></iframe>`;
                    modalHeader.innerHTML = `<h5 class="modal-title">${link.querySelector('.file-name').textContent}</h5>`;
                    // Show the modal
                    const modalInstance = new bootstrap.Modal(document.querySelector('#my-dialog'));
                    modalInstance.show();
                }
            }
        };

        const detectDoubleTapClosure = () => {
            let lastTap = 0;
            let timeout;
            return (event) => {
                const curTime = new Date().getTime();
                const tapLen = curTime - lastTap;
                if (tapLen < 500 && tapLen > 0) {
                    if (!event || !event.target) {
                        return;
                    }
                    const link = event.target.closest('.viewable-file-link');
                    const dirLink = event.target.closest('.folder-item');
                    if (!link && !dirLink) {
                        return;
                    }

                    if (dirLink) {
                        // if it's a folder item, we want to navigate to the folder
                        event.preventDefault();
                        const path = dirLink.getAttribute('href');
                        window.location.href = path;
                        return;
                    }
                    event.preventDefault();
                    clearTimeout(timeout);
                    openViewable(event, link)();
                } else {
                    timeout = setTimeout(() => {
                        clearTimeout(timeout);
                    }, 500);
                }
                lastTap = curTime;
            };
        };
        // viewable-file-link must open in a modal on double click only. single click should not do anything
        const viewableFileLinks = document.querySelectorAll('.viewable-file-link');
        viewableFileLinks.forEach(link => {
            link.addEventListener('click', (event) => {
                event.preventDefault(); // prevent single click action
            });
            link.addEventListener('dblclick', openViewable(event, link));
            // and for mobile devices, we want to open the modal on single tap
        });
        document.body.addEventListener('touchend', detectDoubleTapClosure(), { passive: false });

        uploadButton.addEventListener('click', (event) => {
            if (fileInput.files.length === 0) {
                event.preventDefault();
                showError('Please select at least one file to upload.');
                return;
            }

            // check files size. If any file is larger than 100MB, show an error
            for (let i = 0; i < fileInput.files.length; i++) {
                if (fileInput.files[i].size > 10 * 1024 * 1024 * 1024) { // 100MB
                    event.preventDefault();
                    showError('File size exceeds the limit of 100MB: ' + fileInput.files[i].name);
                    return;
                }
            }

            // start spinner
            spinner.style.display = 'inline-block';

            // upload files
            const formData = new FormData(uploadForm);
            fetch(uploadForm.action, {
                method: 'POST',
                body: formData,
            })
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
                        showError(data.Error || 'Failed to upload files');
                        console.error('Error uploading files:', data.Error);
                    });
                }
            })
            .catch(error => console.error('Error:', error))
            .finally(() => {
                // stop spinner
                spinner.style.display = 'none';
            });
        });

        const addFolder = (event) => {
            const folderName = folderNameInput.value.trim();
            if (!folderName) {
                event.preventDefault();
                showError('Please enter a folder name.');
                return;
            }

            // start spinner
            spinner.style.display = 'inline-block';
            // reset error block
            resetError();

            // create form data
            const formData = new FormData();
            formData.append('name', folderName);
            formData.append('path', uploadForm.querySelector('input[name="path"]').value);

            // create folder
            fetch('/filebrowser/folder', {
                method: 'POST',
                body: formData,
            })
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
                        showError(data.Error || 'Failed to create folder');
                        console.error('Error creating folder:', data.Error);
                    });
                }
            })
            .catch(error => console.error('Error:', error))
            .finally(() => {
                // stop spinner
                spinner.style.display = 'none';
            });

        };

        createFolderButton.addEventListener('click', addFolder);
        folderNameInput.addEventListener('keypress', (event) => {
            if (event.key === 'Enter') {
                addFolder(event);
            }
        });

        // click and touch handlers for switching between list and tile views
        const btnListView = document.querySelector('.btn-list-view');
        const btnTileView = document.querySelector('.btn-tile-view');
        btnListView.addEventListener('click', () => {
            fetch('/filebrowser/mode', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ mode: 'list' })
            })
            .finally(() => {
                // reload the page to apply the new mode
                window.location.reload();
            });
        });
        btnTileView.addEventListener('click', () => {
            let formData = new FormData();
            formData.append('mode', 'tile');
            fetch('/filebrowser/mode', {
                method: 'POST',
                body: formData
            })
            .finally(() => {
                // reload the page to apply the new mode
                window.location.reload();
            });
        });
        // touch
        btnListView.addEventListener('touchstart', (event) => {
            event.preventDefault();
            btnListView.click();
        });
        btnTileView.addEventListener('touchstart', (event) => {
            event.preventDefault();
            btnTileView.click();
        });

        // initialize venobox for images
        new VenoBox({
            selector: '.venobox',
            spinner: 'rotating-plane',
            titleattr: 'data-title',
            numeration: true,
            infinigall: true,
            share: true,
            shareStyle: 'pill',
            titlePosition: 'bottom',
        });

        const popoverTriggerList = document.querySelectorAll('[data-bs-toggle="popover"]')
        const _popoverList = [...popoverTriggerList].map(popoverTriggerEl => new bootstrap.Popover(popoverTriggerEl, {trigger: 'focus'}))
    });
})(document);