document.getElementById("check-button").addEventListener("click", function() {
    const text = document.getElementById("text-input").value;
    if (text.trim() === "") {
        alert("Please enter some text to check.");
        return;
    }

    // Call LanguageTool API to check grammar
    checkGrammar(text);
});

function checkGrammar(text) {
    const url = 'https://api.languagetool.org/v2/check';
    
    // Prepare request data
    const data = {
        text: text,
        language: 'en-US',
    };

    // Fetch request to LanguageTool API
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: new URLSearchParams(data)
    })
    .then(response => response.json())
    .then(data => {
        const correctedText = applyCorrections(text, data.matches);
        displayCorrectedText(correctedText);
    })
    .catch(error => {
        console.error("Error checking grammar:", error);
        alert("Something went wrong. Please try again later.");
    });
}

function applyCorrections(text, matches) {
    let correctedText = text;

    // Loop through all grammar mistakes and apply corrections
    matches.forEach(match => {
        const replacement = match.replacements[0]?.value;  // Choose the first suggested replacement
        if (replacement) {
            const from = match.context.text.substring(match.context.offset, match.context.offset + match.context.length);
            correctedText = correctedText.replace(from, replacement);
        }
    });

    return correctedText;
}

function displayCorrectedText(correctedText) {
    // Show the corrected text in the output area
    document.getElementById("corrected-text").value = correctedText;
    document.getElementById("output-container").style.display = "block";
}
