(function () {
    const alerts = document.querySelectorAll('.alert');
    if (!alerts.length) {
        return;
    }

    setTimeout(() => {
        alerts.forEach((alert) => {
            alert.style.transition = 'opacity 0.4s ease';
            alert.style.opacity = '0.9';
        });
    }, 3000);
})();
