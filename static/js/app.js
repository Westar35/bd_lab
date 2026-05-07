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

document.querySelectorAll('.add-related').forEach((link) => {
    link.addEventListener('click', (event) => {
        const form = link.closest('form');
        if (!form) {
            return;
        }

        event.preventDefault();
        const url = new URL(link.href, window.location.origin);
        const data = new FormData(form);
        for (const [key, value] of data.entries()) {
            if (!key || key.startsWith('_')) {
                continue;
            }
            url.searchParams.set(`draft_${key}`, value);
        }
        window.location.href = url.toString();
    });
});

document.querySelectorAll('[data-cascade-target]').forEach((source) => {
    source.addEventListener('change', async () => {
        const target = document.getElementById(source.dataset.cascadeTarget);
        if (!target) {
            return;
        }

        const url = new URL(source.dataset.cascadeUrl, window.location.origin);
        if (source.value) {
            url.searchParams.set(source.dataset.cascadeParam, source.value);
        }
        if (source.dataset.cascadeExtra) {
            source.dataset.cascadeExtra.split(',').forEach((pair) => {
                const [param, elementID] = pair.split(':');
                const element = document.getElementById(elementID);
                if (param && element && element.value) {
                    url.searchParams.set(param, element.value);
                }
            });
        }

        target.disabled = true;
        const keepFirst = target.options.length ? target.options[0].textContent : 'Все';
        target.innerHTML = `<option value="">${keepFirst}</option>`;

        try {
            const response = await fetch(url.toString(), { headers: { 'Accept': 'application/json' } });
            if (!response.ok) {
                throw new Error('options');
            }
            const options = await response.json();
            options.forEach((option) => {
                const item = document.createElement('option');
                item.value = option.id || option.ID;
                item.textContent = option.label || option.Label;
                target.appendChild(item);
            });
        } catch (error) {
            const item = document.createElement('option');
            item.value = '';
            item.textContent = 'Не удалось загрузить варианты';
            target.appendChild(item);
        } finally {
            target.disabled = false;
        }
    });
});
