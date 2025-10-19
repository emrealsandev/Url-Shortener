// DOM Elements
const form = document.getElementById('shortenForm');
const urlInput = document.getElementById('urlInput');
const customAliasToggle = document.getElementById('customAliasToggle');
const customAliasInput = document.getElementById('customAliasInput');
const customAlias = document.getElementById('customAlias');
const submitBtn = document.getElementById('submitBtn');
const result = document.getElementById('result');
const error = document.getElementById('error');
const loading = document.getElementById('loading');
const shortUrl = document.getElementById('shortUrl');
const originalUrl = document.getElementById('originalUrl');
const copyBtn = document.getElementById('copyBtn');
const newBtn = document.getElementById('newBtn');
const errorMessage = document.getElementById('errorMessage');

// Toggle custom alias input
customAliasToggle.addEventListener('change', (e) => {
    if (e.target.checked) {
        customAliasInput.style.display = 'block';
    } else {
        customAliasInput.style.display = 'none';
        customAlias.value = '';
    }
});

// Form submission
form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const url = urlInput.value.trim();
    if (!url) {
        showError('Lütfen geçerli bir URL girin');
        return;
    }

    // Validate URL format
    try {
        new URL(url);
    } catch {
        showError('Lütfen geçerli bir URL formatı girin (örn: https://example.com)');
        return;
    }

    // Hide previous results/errors
    hideAll();
    showLoading();

    // Prepare request body
    const requestBody = {
        url: url
    };

    if (customAliasToggle.checked && customAlias.value.trim()) {
        requestBody.custom_alias = customAlias.value.trim();
    }

    try {
        const response = await fetch('/v1/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        const data = await response.json();
        showResult(data.short_url, url);
    } catch (err) {
        let errorMsg = 'Bir hata oluştu. Lütfen tekrar deneyin.';

        if (err.message.includes('invalid_url')) {
            errorMsg = 'Geçersiz URL formatı';
        } else if (err.message.includes('conflict')) {
            errorMsg = 'Bu özel link adı zaten kullanılıyor. Lütfen başka bir isim deneyin.';
        } else if (err.message.includes('bad_request')) {
            errorMsg = 'Geçersiz istek. Lütfen bilgileri kontrol edin.';
        } else if (err.message.includes('rate limit')) {
            errorMsg = 'Çok fazla istek gönderdiniz. Lütfen biraz bekleyin.';
        }

        showError(errorMsg);
    }
});

// Copy button
copyBtn.addEventListener('click', async () => {
    const url = shortUrl.href;

    try {
        await navigator.clipboard.writeText(url);

        // Visual feedback
        copyBtn.classList.add('copied');
        const originalHTML = copyBtn.innerHTML;
        copyBtn.innerHTML = `
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
        `;

        setTimeout(() => {
            copyBtn.classList.remove('copied');
            copyBtn.innerHTML = originalHTML;
        }, 2000);
    } catch (err) {
        showError('Kopyalama başarısız oldu');
    }
});

// New shortening button
newBtn.addEventListener('click', () => {
    hideAll();
    form.style.display = 'block';
    urlInput.value = '';
    customAlias.value = '';
    customAliasToggle.checked = false;
    customAliasInput.style.display = 'none';
    urlInput.focus();
});

// Helper functions
function showLoading() {
    loading.style.display = 'flex';
    submitBtn.disabled = true;
}

function showResult(shortLink, originalLink) {
    loading.style.display = 'none';
    form.style.display = 'none';
    result.style.display = 'block';
    submitBtn.disabled = false;

    shortUrl.href = shortLink;
    shortUrl.textContent = shortLink;
    originalUrl.textContent = originalLink;
}

function showError(message) {
    loading.style.display = 'none';
    error.style.display = 'flex';
    errorMessage.textContent = message;
    submitBtn.disabled = false;

    setTimeout(() => {
        error.style.display = 'none';
    }, 5000);
}

function hideAll() {
    result.style.display = 'none';
    error.style.display = 'none';
    loading.style.display = 'none';
    form.style.display = 'block';
}

// Auto-focus on load
window.addEventListener('load', () => {
    urlInput.focus();
});