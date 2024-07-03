document.getElementById("upload-form").addEventListener("submit", function(event){
    event.preventDefault();

    const formData = new FormData();
    formData.append("file", document.getElementById("file").files[0]);
    formData.append("query", document.getElementById("query").value);

    fetch("/predict", {
        method: "POST",
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        const responseDiv = document.getElementById("response");
        responseDiv.innerHTML = `
            <p><strong>Answer:</strong> ${data.answer}</p>
            <p><strong>Coordinates:</strong> ${JSON.stringify(data.coordinates)}</p>
            <p><strong>Cells:</strong> ${data.cells.join(", ")}</p>
            <p><strong>Aggregator:</strong> ${data.aggregator}</p>
        `;
    })
    .catch(error => {
        console.error("Error:", error);
    });
});
