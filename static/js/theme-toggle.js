// Funzione per il cambio dinamico del tema
function toggleTheme() {
    const currentTheme = document.body.classList.contains('dark') ? 'light' : 'dark';
    document.body.classList.toggle('dark', currentTheme === 'dark');
    localStorage.setItem('theme', currentTheme); // Salva la preferenza dell'utente
}

// Carica la preferenza del tema salvato
window.addEventListener('load', () => {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.body.classList.add(savedTheme);
});

// Aggiungi il listener per il pulsante di cambio tema
document.getElementById('theme-toggle-button').addEventListener('click', toggleTheme);
