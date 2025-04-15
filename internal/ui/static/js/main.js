// Main JavaScript for Universal API

document.addEventListener('DOMContentLoaded', function() {
    // Show loading indicator when form is submitted
    const form = document.getElementById('scrapeForm');
    if (form) {
        form.addEventListener('submit', function() {
            document.getElementById('loading').style.display = 'block';
        });
    }
});
