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

function createFields(n) {
    var container = document.getElementById("draftContainer");
    for (let i = 0; i < n; i++) {
        container.insertAdjacentHTML("beforeend", '<div class="field flex-container">' + (i + 1) + '</div>')
    }
}

function distributeFields() {
    var radius = 200;
    var fields = document.getElementsByClassName("field")
    var container = document.getElementById("draftContainer")
    var containerStyles = window.getComputedStyle(container)
    var width = containerStyles.width.replace("px", "")
    var height = containerStyles.height.replace("px", "")
    var angle = 1.5 * Math.PI
    var step = (2 * Math.PI) / fields.length;
    console.log(width, height)
    for (let i = 0; i < fields.length; i++) {
        let fieldWidth = window.getComputedStyle(fields[i]).width.replace("px", "")
        let fieldHeight = window.getComputedStyle(fields[i]).height.replace("px", "")
        var x = Math.round(width / 2 + radius * Math.cos(angle) - fieldWidth / 2);
        var y = Math.round(height / 2 + radius * Math.sin(angle) - fieldHeight / 2);
        fields[i].style.left = x + "px"
        fields[i].style.top = y + "px"
        angle += step;
    }
}