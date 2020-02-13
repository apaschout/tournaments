document.addEventListener("DOMContentLoaded", e => {
    parseDate()
    loadPlayers()
})

function parseDate() {
    let lis = document.getElementsByTagName("li")
    let start = "Start: "
    let end = "End: "
    for (let i = 0; i < lis.length; i++) {
        let d = lis[i].innerHTML
        if (d.includes(start)) {
            newDate = new Date(d.slice(7, 36))
            if (isValidDate(newDate)) {
                lis[i].innerHTML = start + newDate.toUTCString()
            }
        }
        if (d.includes(end)) {
            newDate = new Date(d.slice(5, 34))
            if (isValidDate(newDate)) {
                lis[i].innerHTML = end + newDate.toUTCString()
            }
        }
    }
}

function isValidDate(d) {
    return d instanceof Date && !isNaN(d);
}

function loadPlayers() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = () => {
        if (xhttp.readyState == 4 && xhttp.status == 200) {
            var plrs = JSON.parse(xhttp.responseText)
            var ul = document.getElementById("ul-players")
            // var form = document.getElementById("form-register-player")
            ul.addEventListener("click", e => {
                if (e.target.tagName == "LI") {
                    if (!e.target.classList.contains("selected")) {
                        for (let i = 0; i < plrs.items.length; i++) {
                            if (e.target.id == plrs.items[i].id) {
                                ul.innerHTML += "<input id='hidden-input-selected' type='hidden' value='" + plrs.items[i].id + "'>"
                            }
                        }
                        e.target.classList.add("selected")
                    } else {
                        var inp = document.getElementById("hidden-input-selected")
                        e.target.parentNode.removeChild(inp)
                        e.target.classList.remove("selected")
                    }
                }
            })
            for (let i = 0; i < plrs.items.length; i++) {
                ul.innerHTML += "<li class='w3-hoverable' id='" + plrs.items[i].id + "'>" + plrs.items[i].label + "</li>"
            }
        }
    }
    xhttp.open("GET", "/api/players/", true)
    xhttp.send()
}