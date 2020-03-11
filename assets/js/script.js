function parseDate() {
    let ps = document.getElementsByClassName("date")
    let start = "Started: "
    let end = "End: "
    for (let i = 0; i < ps.length; i++) {
        let d = ps[i].innerHTML
        if (d.includes(start)) {
            newDate = new Date(d.slice(9, 36))
            if (isValidDate(newDate)) {
                ps[i].innerHTML = start + newDate.toUTCString()
            }
        }
        if (d.includes(end)) {
            newDate = new Date(d.slice(5, 34))
            if (isValidDate(newDate)) {
                ps[i].innerHTML = end + newDate.toUTCString()
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
            var form = document.getElementById("form-register-player")
            for (let i = 0; i < plrs.items.length; i++) {
                ul.innerHTML += "<li class='hover-green' id='" + plrs.items[i].id + "'>" + plrs.items[i].label + "</li>"
            }
            ul.addEventListener("click", e => {
                if (e.target.tagName == "LI") {
                    switch (e.target.classList.contains("selected")) {
                        case true:
                            deleteHiddenInput(form)
                            toggleSelected(e.target.classList)
                            break
                        case false:
                            if (document.getElementsByClassName("selected").length == 0) {
                                createHiddenInput(plrs, e.target)
                                toggleSelected(e.target.classList)
                            } else {
                                toggleSelected(document.getElementsByClassName("selected")[0].classList)
                                deleteHiddenInput(form)
                                createHiddenInput(plrs, e.target)
                                toggleSelected(e.target.classList)
                            }
                    }
                }
            })
        }
    }
    xhttp.open("GET", "/api/players/", true)
    xhttp.send()
}

function deleteHiddenInput(form) {
    var inp = document.getElementById("hidden-input-selected")
    form.removeChild(inp)
}

function createHiddenInput(plrs, elem) {
    for (let i = 0; i < plrs.items.length; i++) {
        if (elem.id == plrs.items[i].id) {
            elem.parentNode.parentNode.insertAdjacentHTML("afterbegin", "<input id='hidden-input-selected' type='hidden' name='pid' value='" + plrs.items[i].id + "'>")
        }
    }
}

function toggleSelected(classList) {
    classList.toggle("selected")
    classList.toggle("hover-green")
}

function distributeSeats() {
    var radius = 125;
    var seats = document.getElementsByClassName("seat")
    var container = document.getElementById("draftContainer")
    var containerStyles = window.getComputedStyle(container)
    var width = containerStyles.width.replace("px", "")
    var height = containerStyles.height.replace("px", "")
    var angle = 1.5 * Math.PI
    var step = (2 * Math.PI) / seats.length;
    for (let i = 0; i < seats.length; i++) {
        let seatWidth = window.getComputedStyle(seats[i]).width.replace("px", "")
        let seatHeight = window.getComputedStyle(seats[i]).height.replace("px", "")
        var x = Math.round(width / 2 + radius * Math.cos(angle) - seatWidth / 2);
        var y = Math.round(height / 2 + radius * Math.sin(angle) - seatHeight / 2);
        seats[i].style.left = x + "px"
        seats[i].style.top = y + "px"
        angle += step;
    }
}