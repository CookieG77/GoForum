/**
 * @param {Object} options - options de configuration pour la fonction.
 * @param {HTMLElement} options.container - Le conteneur général auquel la barre de défilement est appliquée. À entrer en argument.
 * @param {string} [options.contentSelector='.scroll-content'] - sélecteur de l'élément dans lequel le contenu est mis.
 * @param {string} [options.scrollbarSelector='.custom-scrollbar'] - sélecteur de la piste de la barre de défilement.
 * @param {string} [options.arrowUpSelector='.scrollbar-arrow.up'] - sélecteur de la flèche du haut de la barre de défilement.
 * @param {string} [options.arrowDownSelector='.scrollbar-arrow.down'] - sélecteur de la flèche du bas de la barre de défilement.
 * @param {string} [options.thumbSelector='.scrollbar-thumb'] - sélecteur du pouce de la barre dde défilement.
 * @param {number} [options.scrollStep=30] - distance de scroll lors du clic d'une des flèches
 */
function initCustomScrollbar(options) {
    const {
        container,       // Requis
        contentSelector = '.scroll-content',
        scrollbarSelector = '.custom-scrollbar',
        arrowUpSelector = '.scrollbar-arrow.up',
        arrowDownSelector = '.scrollbar-arrow.down',
        thumbSelector = '.scrollbar-thumb',
        scrollStep = 30
    } = options;

    // Récupère les sous-éléments relatifs au conteneur.
    const scrollContent = container.querySelector(contentSelector);
    const customScrollbar = container.querySelector(scrollbarSelector);
    const arrowUp = container.querySelector(arrowUpSelector);
    const arrowDown = container.querySelector(arrowDownSelector);
    const scrollbarThumb = container.querySelector(thumbSelector);

    // Fonction pour mettre à jour la position et la taille du pouce.
    function updateThumb() {
        const containerHeight = scrollContent.clientHeight;
        const contentHeight = scrollContent.scrollHeight;
        const scrollTop = scrollContent.scrollTop;

        // Utilise les hauteurs réelles des flèches plutôt que des valeurs fixes.
        const arrowUpHeight = arrowUp.offsetHeight;
        const arrowDownHeight = arrowDown.offsetHeight;
        const availableHeight =  Math.round(containerHeight - arrowUpHeight - arrowDownHeight) + 8;

        // Cacher la barre de défilement si le contenu rentre dans le conteneur.
        if (contentHeight <= containerHeight) {
            customScrollbar.style.display = 'none';
            return;
        } else {
            customScrollbar.style.display = 'block';
        }

        // Calcule la taille du pouce proportionnellement à l'espace disponible.
        const thumbHeight = Math.max(30, (containerHeight / contentHeight) * availableHeight);

        // La zone de déplacement du pouce est la hauteur disponible moins la hauteur du pouce.
        const maxTop = availableHeight - thumbHeight;
        // Position du pouce en se basant sur la position de scroll et en ajoutant l'offset de la flèche du haut.
        const thumbTop = arrowUpHeight + (scrollTop / (contentHeight - containerHeight)) * maxTop;

        scrollbarThumb.style.height = thumbHeight + 'px';
        scrollbarThumb.style.top = thumbTop + 'px';
    }

    // Fonction qui fait défiler vers le haut ou le bas selon un delta.
    function scrollContentBy(delta) {
        scrollContent.scrollTop += delta;
    }

    // Écouteurs d'événements pour les flèches.
    arrowUp.addEventListener('click', () => {
        scrollContentBy(-scrollStep);
    });

    arrowDown.addEventListener('click', () => {
        scrollContentBy(scrollStep);
    });

    // Mise à jour du pouce lors du défilement.
    scrollContent.addEventListener('scroll', updateThumb);

    // Gestion du drag du pouce.
    let isDragging = false;
    let startY, startTop;

    scrollbarThumb.addEventListener('mousedown', (e) => {
        isDragging = true;
        startY = e.clientY;
        startTop = parseInt(window.getComputedStyle(scrollbarThumb).top, 10);
        document.body.style.userSelect = 'none';
    });

    document.addEventListener('mousemove', (e) => {
        if (!isDragging) return;
        const deltaY = e.clientY - startY;
        let newTop = startTop + deltaY;

        const containerHeight = scrollContent.clientHeight;
        const arrowUpHeight = arrowUp.offsetHeight;
        const arrowDownHeight = arrowDown.offsetHeight;
        const availableHeight =  Math.round(containerHeight - arrowUpHeight - arrowDownHeight) + 8;
        const thumbHeight = scrollbarThumb.offsetHeight;
        const maxTop = availableHeight - thumbHeight;

        // Contrain le pouce a resté entre les flèches.
        newTop = Math.max(arrowUpHeight, Math.min(newTop, arrowUpHeight + maxTop));
        scrollbarThumb.style.top = newTop + 'px';

        // Met à jour la position du contenu selon la position du pouce.
        // Attention à éviter une division par zéro.
        const scrollFraction = maxTop ? (newTop - arrowUpHeight) / maxTop : 0;
        scrollContent.scrollTop = scrollFraction * (scrollContent.scrollHeight - containerHeight);
    });

    document.addEventListener('mouseup', () => {
        isDragging = false;
        document.body.style.userSelect = ''; // Réinitialise la sélection de texte
    });

    // Met à jour le pouce lors du redimensionnement de la fenêtre.
    window.addEventListener('resize', updateThumb);

    // Observe les mutations dans le contenu pour ajuster le pouce.
    const observer = new MutationObserver(() => {
        updateThumb();
    });
    observer.observe(scrollContent, { childList: true, subtree: true, characterData: true });

    // Mise à jour initiale.
    updateThumb();
}

// Initialise les barres de défilement.
document.addEventListener('DOMContentLoaded', () => {
    const mainContainer = document.querySelector('main');
    initCustomScrollbar({ container: mainContainer });
});
