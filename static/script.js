require.config({ paths: { 'vs': 'https://unpkg.com/monaco-editor@0.33.0/min/vs' }});
require(['vs/editor/editor.main'], function() {
  const editor = monaco.editor.create(document.getElementById('editor'), {
    value: `# Write your code here\nprint("Hello, World!")`,
    language: 'python',
    theme: 'vs-dark',
  });

  const languageSelector = document.getElementById('language');
  languageSelector.addEventListener('change', () => {
    const language = languageSelector.value;
    monaco.editor.setModelLanguage(editor.getModel(), language);
  });

  const runButton = document.getElementById('run-btn');
  runButton.addEventListener('click', async () => {
    const code = editor.getValue();
    const language = languageSelector.value;
    const outputElement = document.getElementById('output');

    outputElement.textContent = 'Running...';

    try {
      const response = await fetch('/run', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code, language }),
      });
      const result = await response.json();
      outputElement.textContent = result.output || result.error;
    } catch (error) {
      outputElement.textContent = 'Error: ' + error.message;
    }
  });
});
