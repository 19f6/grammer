document.getElementById("correct-btn").addEventListener("click", async function() {
    const textInput = document.getElementById("text-input").value;
    
    if (textInput.trim() === "") {
        alert("Please enter some text.");
        return;
    }

    const response = await fetch('http://127.0.0.1:5000/correct', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ text: textInput })
    });

    const result = await response.json();

    if (response.ok) {
        document.getElementById("corrected-text").textContent = result.corrected_text;
    } else {
        document.getElementById("corrected-text").textContent = `Error: ${result.error}`;
    }
});
