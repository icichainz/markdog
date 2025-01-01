document.addEventListener('DOMContentLoaded', () => {

    const editor = document.getElementById('editor');
    const preview = document.getElementById('preview');
    const newBtn = document.getElementById('newBtn');
    const uploadBtn = document.getElementById('uploadBtn');
    const fileInput = document.getElementById('fileInput');
    const saveBtn = document.getElementById('saveBtn');
    const exportHtmlBtn = document.getElementById('exportHtmlBtn');
    const exportPdfBtn = document.getElementById('exportPdfBtn');

    let debounceTimer;

    // Update preview with debounce
    const updatePreview = async () => {
        const markdown = editor.value;
        try {
            const response = await fetch('http://localhost:3050/api/convert', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ markdown }),
            });

            if (!response.ok) throw new Error('Failed to convert markdown');

            const data = await response.json();
            preview.innerHTML = data.html;
        } catch (error) {
            console.error('Error converting markdown:', error);
        }
    };

    // Debounced preview update
    editor.addEventListener('input', () => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(updatePreview, 300);
    });

    // New document
    newBtn.addEventListener('click', () => {
        if (confirm('Create new document? Current content will be lost.')) {
            editor.value = '';
            preview.innerHTML = '';
        }
    });

    // Handle file upload button click
    uploadBtn.addEventListener('click', () => {
        fileInput.click();
    });

    // Handle file selection
    fileInput.addEventListener('change', async (event) => {

        const file = event.target.files[0];

        if (!file) return;

        try {
            // Create FormData
            const formData = new FormData();
            formData.append('file', file);

            // Upload file
            const response = await fetch('http://localhost:3050/api/upload', {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) throw new Error('Failed to upload file');

            const data = await response.json();
            editor.value = data.content;
            updatePreview();
        } catch (error) {
            console.error('Error uploading file:', error);
            alert('Failed to upload file. Please try again.');
        }

        // Reset file input
        fileInput.value = '';
    });

    // Save markdown document
    saveBtn.addEventListener('click', () => {
        const content = editor.value;
        const blob = new Blob([content], { type: 'text/markdown' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'document.md';
        a.click();
        URL.revokeObjectURL(url);
    });

    // Export HTML
    exportHtmlBtn.addEventListener('click', () => {
        const content = preview.innerHTML;
        const blob = new Blob([content], { type: 'text/html' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'document.html';
        a.click();
        URL.revokeObjectURL(url);
    });

    // Export PDF
    exportPdfBtn.addEventListener('click', async () => {
        const markdown = editor.value;
        try {
            const response = await fetch('http://localhost:3050/api/export-pdf', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ markdown }),
            });

            if (!response.ok) throw new Error('${response.status} ${response.statusText}');

            // Get the PDF blob from the response
            const pdfBlob = await response.blob();
            
            // Create a download link
            const url = URL.createObjectURL(pdfBlob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'document.pdf';
            a.click();
            URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Error generating PDF:', error);
            alert('Failed to generate PDF. Please try again.');
        }
    });
});
