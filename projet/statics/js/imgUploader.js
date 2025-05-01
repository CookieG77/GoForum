
const maxSize = 20 * 1024 * 1024; // 20 Mo
const allowedTypes = ['image/png', 'image/jpeg', 'image/gif'];

/**
 * UploadImages takes a file input that can contains multiple images and uploads them to the server.
 * @param imageHolder {HTMLInputElement} - The file input element containing the images to upload.
 * @param imgType {string} - The type of the images.
 */
async function UploadImages(imageHolder, imgType) {
    const results = [];
    const errors = [];
    const files = imageHolder.files;
    for (const file of files) {
        if (!allowedTypes.includes(file.type)) {
            alert("Format non autorisé (PNG, JPEG, GIF)");
            continue;
        }
        if (file.size > maxSize) {
            alert("Taille de fichier trop grande (max 20 Mo)");
            continue;
        }

        await UploadImg(file, imgType)
            .then((data) => {
                results.push([data.url, data.id]);
                errors.push(null);
                console.log("Image uploaded successfully", data);
            }).catch((error) => {
                results.push(null);
                errors.push(error);
                console.error(error);
            });
    }

    return {
        results: results,
        errors: errors
    };
}

/**
 * UploadImg uploads a single image to the server.
 * @param file {File} - The file to upload.
 * @param imgType {string} - The type of the image.
 * @returns {Promise<any>}
 */
function UploadImg(file, imgType) {
    if (!allowedTypes.includes(file.type)) {
        alert("Format non autorisé (PNG, JPEG, GIF)");
        return null;
    }
    if (file.size > maxSize) {
        alert("Taille de fichier trop grande (max 20 Mo)");
        return null;
    }
    const formData = new FormData();
    formData.append(`image` , file);
    return fetch(`/api/upload/${imgType}`, {
        method: "POST",
        body: formData
    }).then((response) => {
        if (!response.ok) {
            throw new Error("Erreur lors de l'upload de l'image : " + response.statusText);
        }
        return response.json();
    }).then((data) => {
        return data;
    }).catch((error) => {
        console.error(error);
    });
}