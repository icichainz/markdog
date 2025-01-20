document.addEventListener('DOMContentLoaded', () => {

    const editor = document.getElementById('editor');
    const preview = document.getElementById('preview');
    const newBtn = document.getElementById('newBtn');
    const uploadBtn = document.getElementById('uploadBtn');
    const fileInput = document.getElementById('fileInput');
    const saveBtn = document.getElementById('saveBtn');
    const exportHtmlBtn = document.getElementById('exportHtmlBtn');
    const exportPdfBtn = document.getElementById('exportPdfBtn');
    const exportDocxBtn = document.getElementById('exportDocxBtn');
    const scrollSyncBtnToggle = document.getElementById('scrollSyncBtnToggle');
    const scrollSyncToggleEl = document.getElementById('scrollSyncBtn');

    let scrollSyncToggleStatus = false;
    const currentDocumentKey = "current_doc";
    

    let debounceTimer;
   
    function toggleScrollSync() {
        scrollSyncToggleStatus = !scrollSyncToggleStatus;
        updateScrollListeners();
    }

    scrollSyncBtnToggle.addEventListener('click', () => {
        toggleScrollSync();
        scrollSyncToggleEl.innerText = scrollSyncToggleStatus ? 'ON' : 'OFF';
    });

    function onLoad(){
        if (localStorage.getItem(currentDocumentKey)) {
            editor.value = localStorage.getItem(currentDocumentKey) ;
        }
        
        updatePreview();

    }
    function addToStorage(value) {
        localStorage.setItem(currentDocumentKey,value)

    }

    function clearStorage(){
        localStorage.removeItem(currentDocumentKey);
    }
    // Update preview with debounce
    const updatePreview = async () => {
        const markdown = editor.value;
         addToStorage(markdown) ;

         console.log('the data : '+localStorage.getItem(currentDocumentKey));
        
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

        clearStorage();
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
    exportPdfBtn.addEventListener('click', () => {
        const element = preview;
        const options = {
            margin: 0.2,
            filename: 'document.pdf',
            image: { type: 'jpec', quality: 0.98 },
            html2canvas: { scale: 2 },
            jsPDF: {
                unit: 'in',
                format: 'letter',
                orientation: 'portrait',
            }
        };

        html2pdf().set(options).from(element).save();
    });

    // Export DOCX
    exportDocxBtn.addEventListener('click', async () => {
        const { Document, Packer, Paragraph } = docx;

        const content = preview.innerHTML; // Get the HTML content from preview

        // Convert HTML to plain text or handle it for Word compatibility
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = content;
        const textContent = tempDiv.textContent || tempDiv.innerText || '';

        // Create a Word document
        const doc = new Document({
            sections: [
                {
                    properties: {},
                    children: [
                        new Paragraph({
                            text: textContent,
                            spacing: { after: 200 },
                        }),
                    ],
                },
            ],
        });

        // Generate the .docx file
        const blob = await Packer.toBlob(doc);
        const url = URL.createObjectURL(blob);

        // Trigger download
        const a = document.createElement('a');
        a.href = url;
        a.download = 'document.docx';
        a.click();

        URL.revokeObjectURL(url);
    });

    let isScrolling = false;
    let scrollDebounceTimer;
    
    function debounce(func, wait) {
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(scrollDebounceTimer);
                func(...args);
            };
            clearTimeout(scrollDebounceTimer);
            scrollDebounceTimer = setTimeout(later, wait);
        };
    }
    
    function syncScroll(source, target) {
        if (!scrollSyncToggleStatus || isScrolling) return;
        
        isScrolling = true;
        const sourceScrollPercent = source.scrollTop / (source.scrollHeight - source.clientHeight);
        const targetScrollPosition = sourceScrollPercent * (target.scrollHeight - target.clientHeight);
        target.scrollTop = targetScrollPosition;
        
        setTimeout(() => {
            isScrolling = false;
        }, 50);
    }

    const debouncedSync = debounce((source, target) => syncScroll(source, target), 16);

    // Add/remove scroll listeners based on toggle status
    function updateScrollListeners() {
        if (scrollSyncToggleStatus) {
            editor.addEventListener('scroll', () => debouncedSync(editor, preview), { passive: true });
            preview.addEventListener('scroll', () => debouncedSync(preview, editor), { passive: true });
        } else {
            editor.removeEventListener('scroll', () => debouncedSync(editor, preview));
            preview.removeEventListener('scroll', () => debouncedSync(preview, editor));
        }
    }

    // Initialize scroll sync state
    updateScrollListeners();

    onLoad();

});
