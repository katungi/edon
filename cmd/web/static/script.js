document.addEventListener('DOMContentLoaded', () => {
    const input = document.getElementById('input');
    const output = document.getElementById('output');
    let history = [];
    let historyIndex = -1;

    function appendOutput(text, isError = false) {
        const line = document.createElement('div');
        line.className = `output-line ${isError ? 'error' : ''}`;
        line.textContent = text;
        output.appendChild(line);
        output.scrollTop = output.scrollHeight;
    }

    async function executeCommand(command) {
        try {
            const response = await fetch('/eval', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ code: command }),
            });
            
            const result = await response.json();
            if (result.error) {
                appendOutput(result.error, true);
            } else {
                appendOutput(result.output);
            }
        } catch (error) {
            appendOutput('Error connecting to server: ' + error.message, true);
        }
    }

    input.addEventListener('keydown', async (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            const command = input.value.trim();
            
            if (command) {
                appendOutput('>> ' + command);
                history.push(command);
                historyIndex = history.length;
                await executeCommand(command);
                input.value = '';
            }
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            if (historyIndex > 0) {
                historyIndex--;
                input.value = history[historyIndex];
            }
        } else if (e.key === 'ArrowDown') {
            e.preventDefault();
            if (historyIndex < history.length - 1) {
                historyIndex++;
                input.value = history[historyIndex];
            } else {
                historyIndex = history.length;
                input.value = '';
            }
        }
    });
});
