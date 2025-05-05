document.addEventListener("DOMContentLoaded", function () {
    const threadName = getCurrentThreadName(); // Remplace par le nom réel du thread si nécessaire

    const tagNameInput = document.getElementById('tagName');
    const tagColorInput = document.getElementById('tagColor');
    const addTagBtn = document.getElementById('addTagBtn');
    const tagList = document.getElementById('tagList');
    const editTagList = document.getElementById('editTagList');
    const errorDiv = document.getElementById('error');

    const tags = {};
    const tagIds = {}; // map tagName -> tagId

    function renderTags() {
        tagList.innerHTML = '';
        editTagList.innerHTML = '';
        Object.keys(tags).forEach(name => {
            const color = tags[name];
            const tagId = tagIds[name];

            const tagEl = document.createElement('div');
            tagEl.className = 'tag';
            tagEl.style.backgroundColor = color;
            tagEl.innerHTML = `${name} <button data-name="${name}">x</button>`;
            tagList.appendChild(tagEl);

            const editEl = document.createElement('div');
            editEl.className = 'tag-editor';
            editEl.innerHTML = `
          <input type="text" value="${name}" data-old-name="${name}" />
          <input type="color" value="${color}" />
          <button data-action="save">Save</button>
          <button data-action="delete">Delete</button>
        `;
            editTagList.appendChild(editEl);

            const saveBtn = editEl.querySelector('[data-action="save"]');
            const deleteBtn = editEl.querySelector('[data-action="delete"]');
            const nameInput = editEl.querySelector('input[type="text"]');
            const colorInput = editEl.querySelector('input[type="color"]');

            saveBtn.addEventListener('click', () => {
                const newName = nameInput.value.trim();
                const newColor = colorInput.value;
                const oldName = nameInput.dataset.oldName;

                if (!newName) {
                    alert('Tag name cannot be empty.');
                    return;
                }

                if (newName !== oldName && tags.hasOwnProperty(newName)) {
                    alert('Tag name must be unique.');
                    return;
                }

                const tagId = tagIds[oldName];

                editThreadTag(threadName, tagId, newName, newColor)
                    .then(response => {
                        if (!response.ok) throw new Error('Failed to edit tag.');
                        delete tags[oldName];
                        delete tagIds[oldName];
                        tags[newName] = newColor;
                        tagIds[newName] = tagId;
                        renderTags();
                    })
                    .catch(err => {
                        alert(err.message);
                    });
            });

            deleteBtn.addEventListener('click', () => {
                const tagId = tagIds[name];
                deleteThreadTag(threadName, tagId)
                    .then(response => {
                        if (!response.ok) throw new Error('Failed to delete tag.');
                        delete tags[name];
                        delete tagIds[name];
                        renderTags();
                    })
                    .catch(err => {
                        alert(err.message);
                    });
            });
        });

        document.querySelectorAll('.tag button').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const name = e.target.getAttribute('data-name');
                const tagId = tagIds[name];
                deleteThreadTag(threadName, tagId)
                    .then(response => {
                        if (!response.ok) throw new Error('Failed to delete tag.');
                        delete tags[name];
                        delete tagIds[name];
                        renderTags();
                    })
                    .catch(err => {
                        alert(err.message);
                    });
            });
        });
    }

    addTagBtn.addEventListener('click', () => {
        console.log("Add tag button clicked");
        const name = tagNameInput.value.trim();
        const color = tagColorInput.value;

        if (!name) {
            errorDiv.textContent = 'Tag name cannot be empty.';
            return;
        }

        if (tags.hasOwnProperty(name)) {
            errorDiv.textContent = 'Tag name must be unique.';
            return;
        }

        errorDiv.textContent = '';

        createThreadTag(threadName, name, color)
            .then(response => {
                if (!response.ok) throw new Error('Failed to create tag.');
                return response.json();
            })
            .then(data => {
                tags[name] = color;
                tagIds[name] = data.tag_id;
                tagNameInput.value = '';
                tagColorInput.value = '#ff0000';
                renderTags();
            })
            .catch(err => {
                errorDiv.textContent = err.message;
                console.error('Failed to create tag:', err);
            });
    });

    // Chargement initial
    getThreadTags(threadName)
        .then(response => {
            console.log(response)
            if (!response.ok) throw new Error('Failed to load tags.');
            return response.json();
        })
        .then(data => {
            data.forEach(tag => {
                tags[tag.tag_name] = tag.tag_color;
                tagIds[tag.tag_name] = tag.tag_id;
            });
            renderTags();
        })
        .catch(err => {
            console.error('Failed to load tags:', err);
            errorDiv.textContent = err.message;
        });
});