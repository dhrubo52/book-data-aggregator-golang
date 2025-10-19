const getData = () => {
    let request = new XMLHttpRequest();

    let titleInputElement = document.getElementById("title-input");
    
    let loadingText = document.getElementById("loading");
    loadingText.hidden = false;

    if(!titleInputElement.value.trim()) {
        alert("Please enter a title name");
        return
    }

    let params = new URLSearchParams({
        title: titleInputElement.value.trim()
    });

    request.open("GET", `/get-data?${params.toString()}`, true);
    
    let resultElement = document.getElementById("result-text");

    let keyMap = {
        totalBooks: "Total number books in Open Library that has similar title",
        earliestPublicationYear: "The earliest publication year among the books",
        latestPublicationYear: "The latest publication year among the books",
        authors: "List of all authors of the books",
        languages: "List of languages of the books"
    }

    request.onload = () => {
        if(request.status>=200 && request.status<300) {
            let resultHtml = "";
            let responseJSON = JSON.parse(request.responseText);
            for(const [key, value] of Object.entries(responseJSON)) {
                if(key==="totalBooks" && value===0) {
                    resultHtml = "No books found";
                    break
                } 
                resultHtml = resultHtml + `<p class="mt-1">${keyMap[key]}: ${(key!=="authors" && key!=="languages") ? value : ""}<p>`;

                if(key==="authors") {
                    resultHtml += "<ul class='list mt-1'>"
                    for(let i=0; i<value.length; ++i) {
                        resultHtml += `<li>${value[i]}</li>`
                    }
                    resultHtml += "</ul>"
                }

                if(key==="languages") {
                    resultHtml += "<ul class='list mt-1'>"
                    for(let i=0; i<value.length; ++i) {
                        resultHtml += `<li>${value[i]}</li>`
                    }
                    resultHtml += "</ul>"
                }
            }

            loadingText.hidden = true;
            resultElement.innerHTML = resultHtml;
        }
        else {
            resultElement.innerHTML = `<p>${request.responseText}</p>`;
        }
    }

    request.send()
}